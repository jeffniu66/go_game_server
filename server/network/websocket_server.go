package network

import (
	"encoding/binary"
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/handler"
	"go_game_server/server/logger"
	"log"
	"net/http"
	"os"
	"reflect"
	"runtime/debug"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

var clientMutex sync.Mutex

var ServerStatus bool // 默认关闭状态

type wsClient struct {
	readClose  bool
	login      int
	userId     int32
	socket     websocket.Conn
	userPid    *global.PidObj
	touristPid *global.PidObj
}

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Message  string `json:"message"`
}

func StartWsServer() {
	cls := make(chan os.Signal, 1)
	go Signals(cls) // 用于接收信号量
	// Configure websocket route
	http.HandleFunc("/ws", handleServer)

	// Start the server on localhost port 8000 and log any errors
	//go game.SchedulePushClient() // 定时帧同步数据
	addr := global.MyConfig.Read("server", "address")
	logger.Log.Infof("websocket server started on %s\n", addr)
	go func() {
		httpType := global.MyConfig.Read("http", "http_type")
		switch httpType {
		case "http":
			err := http.ListenAndServe(addr, nil)
			if err != nil {
				log.Fatal("ListenAndServe: ", err)
			}
		case "https":
			err := http.ListenAndServeTLS(addr, "./config/ninja.crt", "./config/ninja.key", nil)
			if err != nil {
				log.Fatal("ListenAndServeTLS: ", err)
			}
		}
	}()

	// 结束主线程
	for {
		flag := false
		sig := <-cls
		switch sig {
		case syscall.SIGINT:
			flag = true
		default:
		}
		if flag {
			break
		}
	}
}

func handleServer(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Log.Errorf("client web socket conn disconnect, url:%v err:%v", r.URL.Path, err)
		http.NotFound(w, r)
		return
	}
	logger.Log.Info("client remote addr: ", ws.RemoteAddr().String())
	client := &wsClient{}
	defer func() {
		// 移除连接
		// delete(game.ConnMap, ws.RemoteAddr().String())
		_ = ws.Close()
		client.close()
		deleteClient(ws)
		logger.Log.Info("client disconnect!!!")
	}()

	logger.Log.Info("new websocket client : ", len(handler.Clients))
	writeChan := make(chan interface{}, 1024)
	go client.writeServer(ws, writeChan) // 是否创建玩家成功后才开这个协程

	//game.ConnMap[ws.RemoteAddr().String()] = ws

	// Register our new client
	handler.Clients[ws] = true
	client.readServer(ws, writeChan)
}

func (w *wsClient) readServer(ws *websocket.Conn, writeChan chan interface{}) {
	for {
		err := ws.SetReadDeadline(time.Now().Add(time.Duration(time.Second * timeout)))
		if err != nil {
			logger.Log.Warnf("ws userID: %v timeout:%v SetReadDeadline failed", w.userId, timeout)
			return
		}
		_, msgBin, err := ws.ReadMessage()
		if err != nil {
			logger.Log.Errorf("read server error~~~~: %v", err)
			return
		}
		msgID := binary.LittleEndian.Uint16(msgBin[0:2])
		cmd := proto3.ProtoCmd(msgID)
		if cmd < 0 {
			return
		}
		pbData := game.Handler.GetPbData(cmd)
		msgPb := reflect.New(reflect.TypeOf(pbData).Elem()).Interface()
		err = proto.Unmarshal(msgBin[2:], msgPb.(proto.Message))
		if err != nil {
			logger.Log.Info("proto unmarshal err: ", err)
			return
		}
		// logger.Log.Info("time>>>>>>>>msg>>>>>>>>>>>>>cmd:", cmd, " pb:", msgPb, " client:", ws.RemoteAddr().String())

		// 无账号
		if cmd == proto3.ProtoCmd_CMD_RandNameReq || cmd == proto3.ProtoCmd_CMD_RegisterReq || cmd == proto3.ProtoCmd_CMD_CreateUserReq {
			if w.touristPid == nil {
				w.touristPid = game.CreateTourist(writeChan)
			}
			msgData := &game.Message{Cmd: cmd, PbData: msgPb}
			w.touristPid.Cast("SOCKET_EVENT", msgData)
			continue
		}

		if w.login == login && cmd == proto3.ProtoCmd_CMD_LoginReq { // 未断开socket就重复发登录协议
			return
		}

		if w.login == logout && cmd == proto3.ProtoCmd_CMD_LoginReq { // 第一次登陆
			msg := msgPb.(*proto3.LoginReq)
			ip := getClientIp(ws)
			ret, userPid, userId := game.CreatePid(writeChan, msg, ip)
			logger.Log.Infof("CreatePid ret = %v, userPid = %v, userId = %v", ret, userPid, userId)
			if ret != 0 {
				continue
			}
			w.userPid = userPid
			w.userId = userId
			w.login = login
		}

		if w.login == login {
			// 只在登录情况下 异步向玩家协程发送消息
			msgData := &game.Message{Cmd: cmd, PbData: msgPb}
			if cmd != proto3.ProtoCmd_CMD_FrameSyncReq && cmd != proto3.ProtoCmd_CMD_HeartBeatReq {
				logger.Log.Infof("------- userId = %v, cmd = %v, pbData = %v\n", w.userId, cmd, msgPb)
			}
			w.userPid.Cast("SOCKET_EVENT", msgData)
		}
	}
}

func (w *wsClient) writeServer(ws *websocket.Conn, writeChan chan interface{}) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("websocket_server send packet err000---: ", err, len(writeChan), string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
			debug.PrintStack()
		}
	}()
	for {
		msgData, ok := <-writeChan
		if !ok {
			//c.readClose = true
			//_ = c.socket.Close() // 通知read协程退出
			logger.Log.Infof("userID:%v player quit server", w.userId)
			_ = ws.Close()
			deleteClient(ws)
			return
		}
		if msgData != nil {
			switch msg := msgData.(type) {
			case *game.Message:
				if msg.Cmd != proto3.ProtoCmd_CMD_FrameSyncResp && msg.Cmd != proto3.ProtoCmd_CMD_HisFSPFrameResp && msg.Cmd != proto3.ProtoCmd_CMD_HeartBeatResp {
					logger.Log.Infof("------- userId = %v, cmd = %v,  pbData = %v\n", w.userId, msg.Cmd, msg.PbData)
				}

				SendWsMessage(ws, msg.Cmd, msg.PbData)
			//if err != nil {
			//	_ = ws.Close()
			//	return
			//}
			default:
				logger.Log.Errorln("tcp_server send packet err111: ", msg)
			}
		}
	}
}

var sendMsgSize int64

func SendWsMessage(ws *websocket.Conn, cmd proto3.ProtoCmd, pbData interface{}) error {
	msgData, err := proto.Marshal(pbData.(proto.Message))
	if err != nil {
		fmt.Println("proto marshal err : ", err)
	}
	msg := make([]byte, len(msgData)+2)
	binary.LittleEndian.PutUint16(msg, uint16(cmd))
	copy(msg[2:], msgData)
	sendMsgSize += int64(len(msg))
	err = ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		logger.Log.Info("websocket send message err :", err)
		return err
	}
	return nil
}

func getClientIp(ws *websocket.Conn) (ip string) {
	addrStr := ws.RemoteAddr().String()
	index := strings.Index(addrStr, ":")
	if index > 0 {
		ip = addrStr[:index]
	}
	return
}

func (w *wsClient) close() {
	logger.Log.Infof("------------------- c.userPid stop:%v ", w.userPid)
	if w.userPid != nil {
		w.userPid.CastStop()
	}
}

func deleteClient(ws *websocket.Conn) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("websocket_server deleteClient---: ", err, ws, string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
			debug.PrintStack()
		}
	}()
	clientMutex.Lock()
	delete(handler.Clients, ws)
	clientMutex.Unlock()
}
