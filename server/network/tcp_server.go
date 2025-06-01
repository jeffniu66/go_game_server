package network

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"net"
	"runtime/debug"
	"strings"
	"time"
)

var Server ServerType

const (
	logout  = iota // 未登录
	login          // 已登陆
	timeout = 60   // 心跳超时
)

type ServerType struct {
	Done    chan struct{}
	Connect net.Listener
}

type client struct {
	readClose bool
	login     int
	userID    int32
	socket    net.Conn
	userPid   *global.PidObj
}

func StartTcpServer() {
	listener, err := net.Listen("tcp", *global.ServerPort)
	logger.Log.Infof(">>>>>>>>>>>> tcp网络启动成功,端口%v \n\n", *global.ServerPort)
	Server.Connect = listener
	defer listener.Close()

	if err != nil {
		logger.Log.Errorln("Error listening", err.Error())
		return
	}

	//listen and receive client's connection
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				continue
			}
			client := &client{socket: conn}
			client.handleServer(conn) // 是否需要开协程
		}
	}()

	// go Signals() // 用于接收信号量

	Server.Done = make(chan struct{})
	logger.Log.Infof(">>>>>>>>>>>> 服务器启动成功 !!! <<<<<<<<<<<<<<\n\n")
	<-Server.Done // 用来阻塞main主协程,以免应用结束
}

func (c *client) handleServer(conn net.Conn) {
	writeChan := make(chan interface{}, 1024)
	go c.readServer(conn, writeChan)
	go c.writeServer(writeChan) // 是否创建玩家成功后才开这个协程

}

// 阻塞式读，通过read的标志位来退出阻塞式循环
func (c *client) readServer(conn net.Conn, writeChan chan interface{}) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("tcp_server receive packet err: ", err, string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
			debug.PrintStack()
		}
	}()
	for {
		// 网络问题：心跳包设置必须放在解包前面，以防解包出错还可以通过心跳断开socket
		_ = conn.SetDeadline(time.Now().Add(time.Duration(time.Second * timeout))) // 设置tcp心跳时间
		cmd, _, pbData, err := RecvLTVPacket(conn)
		if err != nil {
			if c.readClose { // 服务器踢人最后一步就是conn关闭socket后，通知读协程退出
				logger.Log.Warnln(" server kickplayer close: ", c.userID, err)
			} else {
				logger.Log.Warnln(" client disconnect: ", c.userID, err, cmd)
				c.close()
			}
			return
		}
		if pbData == nil {
			logger.Log.Warnln(" client disconnect: ", c.userID, err)
			c.close()
			return
		}

		if c.login == login && cmd == proto3.ProtoCmd_CMD_LogoutReq { // 客户端发送退出协议
			c.userPid.Stop()
			return
		}

		if c.login == login && cmd == proto3.ProtoCmd_CMD_LoginReq { // 未断开socket就重复发登录协议
			return
		}

		if c.login == logout && cmd == proto3.ProtoCmd_CMD_LoginReq { // 第一次登陆
			msg := pbData.(*proto3.LoginReq)
			ip := c.getClientIp()
			ret, userPid, userID := game.CreatePid(writeChan, msg, ip)
			if ret != 0 {
				continue
			}
			c.userPid = userPid
			c.userID = userID
			c.login = login
		}
		// 异步向玩家协程发送消息
		msgData := &game.Message{Cmd: cmd, PbData: pbData}
		c.userPid.Cast("SOCKET_EVENT", msgData)
	}
}

// writeChan 阻塞式读，所以不会一直死循环,有数据才循环
func (c *client) writeServer(writeChan chan interface{}) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("tcp_server send packet err000---: ", err, len(writeChan), string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
		}
	}()
	for {
		msgData, ok := <-writeChan
		if !ok {
			c.readClose = true
			c.socket.Close() // 通知read协程退出
			return
		}
		if msgData != nil {
			switch msg := msgData.(type) {
			case *game.Message:
				fmt.Printf("send userId = %d, cmd = %v,  pbData = %v\n", c.userID, msg.Cmd, msg.PbData)
				SendLTVPacket(c.socket, msg.Cmd, msg.PbData)
			default:
				logger.Log.Errorln("tcp_server send packet err111: ", msg)
			}
		}
	}
}

func (c *client) close() {
	_ = c.socket.Close() // 网络问题：必须先断开socket，因为c.userPid.Stop可能会报错
	if c.userPid != nil {
		c.userPid.Stop()
	}
}

func (c *client) getClientIp() (ip string) {
	addrStr := c.socket.RemoteAddr().String()
	index := strings.Index(addrStr, ":")
	if index > 0 {
		ip = addrStr[:index]
	}

	return
}
