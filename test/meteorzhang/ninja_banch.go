package main

import (
	"encoding/binary"
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/game"
	"go_game_server/server/network"
	"go_game_server/server/util"
	"log"
	"reflect"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
)

var roomNum int = 100
var matchNum int = 10
var timeOut int = 2
var linehost string = "ws://42.192.186.40:80/ws"
var localhost string = "ws://localhost:80/ws"

func AiLogin(roomL, mateNum int) {
	roomNum = roomL
	matchNum = mateNum
	log.Printf("------start branch room:%d, userNum:%d\n", roomNum, 10)
	for i := 0; i < roomNum; i++ {
		for u := 0; u < matchNum; u++ {
			p := make(chan interface{})
			c, _, err := websocket.DefaultDialer.Dial(localhost, nil)
			defer c.Close()
			if err != nil {
				log.Fatal("dial:", err)
			}
			username := "AI" + util.ToStr(int32(i)) + "-" + util.ToStr(int32(u))
			go func(username string, cc *websocket.Conn, pt chan interface{}) {
				cmd := proto3.ProtoCmd_CMD_LoginReq
				pbData := &proto3.LoginReq{}
				pbData.Username = username
				// pbData.AcctName = "hello"
				network.SendWsMessage(c, cmd, pbData)

				msg := <-pt
				req := msg.(*game.Message)
				reqPb := req.PbData.(*proto3.LoginResp)
				player := reqPb.PlayerAttr
				if reqPb.IsRoom != proto3.CommonStatusEnum_true {
					_ = player
				}
			}(username, c, p)
			go func(cc *websocket.Conn, pt chan interface{}) {
				// 接收服务器回复的数据
				for {
					_, msgBin, err := cc.ReadMessage()
					if err != nil {
						log.Printf("banch test read server error~~~~: %v\n", err)
						_ = cc.Close()
						// 移除连接
						return
					}
					msgID := binary.LittleEndian.Uint16(msgBin[0:2])
					cmd := proto3.ProtoCmd(msgID)
					if cmd < 0 {
						return
					}
					if cmd == proto3.ProtoCmd_CMD_LoginResp {
						pbData := &proto3.LoginResp{}
						msgPb := reflect.New(reflect.TypeOf(pbData).Elem()).Interface()
						err = proto.Unmarshal(msgBin[2:], msgPb.(proto.Message))
						if err != nil {
							fmt.Println("proto unmarshal err: ", err)
							_ = cc.Close()
						}
						msgData := &game.Message{Cmd: cmd, PbData: msgPb}
						pt <- msgData
					}
				}
			}(c, p)
		}
	}

	log.Printf("------end branch room:%d, userNum:%d\n", roomNum, 10)
	select {}
}

func AiLoginAndMatch(one, mateNum int, tag, addr string) {
	roomNum = one
	matchNum = mateNum
	localhost = addr
	log.Printf("------start branch room:%d, userNum:%d\n", roomNum, 10)
	for i := 0; i < roomNum; i++ {
		time.Sleep(time.Second * 1)
		for u := 0; u < matchNum; u++ {
			p := make(chan interface{})
			c, _, err := websocket.DefaultDialer.Dial(localhost, nil)
			defer c.Close()
			if err != nil {
				log.Fatal("dial:", err)
			}
			username := tag + util.ToStr(int32(i)) + "-" + util.ToStr(int32(u))
			go func(username string, cc *websocket.Conn, pt chan interface{}) {
				defer func() { // 必须要先声明defer，否则不能捕获到panic异常
					if err := recover(); err != nil {
						fmt.Println(err)
					}
				}()
				cmd := proto3.ProtoCmd_CMD_LoginReq
				pbData := &proto3.LoginReq{}
				pbData.Username = username
				// pbData.AcctName = "hello"
				network.SendWsMessage(c, cmd, pbData)

				msg := <-pt
				req := msg.(*game.Message)
				reqPb := req.PbData.(*proto3.LoginResp)
				player := reqPb.PlayerAttr
				if reqPb.IsRoom != proto3.CommonStatusEnum_true {
					cmd = proto3.ProtoCmd_CMD_EnterMatchReq
					pbData1 := &proto3.EnterMatchReq{
						Key: proto3.MatchEnum_enter_match,
					}
					log.Printf("enterMatch:%v", player.Username)
					network.SendWsMessage(c, cmd, pbData1)
				}
				////////////////////////////// test frame ////////////////////////////
				msg = <-pt
				re := msg.(*game.Message)
				_, ok := re.PbData.(*proto3.BeginGameResp)
				if ok {
					cmd = proto3.ProtoCmd_CMD_FrameSyncReq
					pbData2 := &proto3.FSPC2SDataReq{}
					fspMsg := &proto3.FSPMsg{
						UId: uint32(player.UserId),
						Cmd: proto3.FSPCmd_CMD_IDLE,
						Args: &proto3.FSPCmdArgs{
							X:   0,
							Y:   0,
							Dir: 0,
						},
					}
					pbData2.Msgs = []*proto3.FSPMsg{fspMsg}
					for {
						time.Sleep(30 * time.Second)
						network.SendWsMessage(c, cmd, pbData2)
					}
				}
			}(username, c, p)
			go func(cc *websocket.Conn, pt chan interface{}) {
				defer func() { // 必须要先声明defer，否则不能捕获到panic异常
					if err := recover(); err != nil {
						fmt.Println("receive: ", err)
					}
				}()
				// 接收服务器回复的数据
				for {
					_, msgBin, err := cc.ReadMessage()
					if err != nil {
						log.Printf("banch test read server error~~~~: %v\n", err)
						_ = cc.Close()
						// 移除连接
						return
					}
					msgID := binary.LittleEndian.Uint16(msgBin[0:2])
					cmd := proto3.ProtoCmd(msgID)
					if cmd < 0 {
						return
					}
					if cmd == proto3.ProtoCmd_CMD_LoginResp {
						pbData := &proto3.LoginResp{}
						msgPb := reflect.New(reflect.TypeOf(pbData).Elem()).Interface()
						err = proto.Unmarshal(msgBin[2:], msgPb.(proto.Message))
						if err != nil {
							fmt.Println("proto unmarshal err: ", err)
							_ = cc.Close()
						}
						msgData := &game.Message{Cmd: cmd, PbData: msgPb}
						pt <- msgData
					}
					// 游戏开始
					if cmd == proto3.ProtoCmd_CMD_BeginGameResp {
						pbData := &proto3.BeginGameResp{}
						msgPb := reflect.New(reflect.TypeOf(pbData).Elem()).Interface()
						err = proto.Unmarshal(msgBin[2:], msgPb.(proto.Message))
						if err != nil {
							fmt.Println("proto unmarshal err: ", err)
							_ = cc.Close()
						}
						msgData := &game.Message{Cmd: cmd, PbData: msgPb}
						pt <- msgData
					}
				}
			}(c, p)
		}
	}

	log.Printf("------end branch room:%d, userNum:%d\n", roomNum, 10)
	select {}
}
