package main

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/network"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main1() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8005/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		cmd := proto3.ProtoCmd_CMD_LoginReq
		pbData := &proto3.LoginReq{}
		// pbData.AcctName = "hello3"

		network.SendWsMessage(c, cmd, pbData)

		time.Sleep(time.Second)

		cmd = proto3.ProtoCmd_CMD_EnterMatchReq
		pbData1 := &proto3.EnterMatchReq{}
		pbData1.Key = proto3.MatchEnum_enter_match
		network.SendWsMessage(c, cmd, pbData1)

		////////////////////////////// test frame ////////////////////////////

		time.Sleep(10 * time.Second)

		cmd = proto3.ProtoCmd_CMD_StoreReq
		pbData3 := &proto3.StoreReq{}
		network.SendWsMessage(c, cmd, pbData3)

		//cmd = proto3.ProtoCmd_CMD_UrgencyTaskReq
		//pbData3 := &proto3.UrgencyTaskReq{}
		//pbData3.TriggerPoint = 21
		//network.SendWsMessage(c, cmd, pbData3)

		//cmd = proto3.ProtoCmd_CMD_FrameSyncReq
		//pbData2 := &proto3.FSPC2SDataReq{}
		//fspMsg := &proto3.FSPMsg{}
		//fspMsg.UId = 1
		//pbData2.Msgs = []*proto3.FSPMsg{fspMsg}
		//
		//network.SendWsMessage(c, cmd, pbData2)
		//
		//time.Sleep(5 * time.Second)
		//network.SendWsMessage(c, cmd, pbData2)
		//time.Sleep(5 * time.Second)
		//network.SendWsMessage(c, cmd, pbData2)
		//time.Sleep(5 * time.Second)
		//network.SendWsMessage(c, cmd, pbData2)

	}()

	// 接收服务器回复的数据
	for {
		t, msg, err := c.ReadMessage() // 接收服务器的请求
		if err != nil {
			fmt.Println("c.ReadMessage err = ", err, t)
			return
		}
		fmt.Println("rec data: ", msg) // 打印接收到的请求
	}

	//msgData, err := proto.Marshal(pbData)
	//msg := make([]byte, len(msgData)+2)
	//binary.LittleEndian.PutUint16(msg, uint16(cmd))
	//copy(msg[2:], msgData)
	//fmt.Println("client 1111111 ", msg, msgData, msg[2:])
	//err = c.WriteMessage(websocket.BinaryMessage, msg)
	//if err != nil {
	//	log.Println("llllllll:", err)
	//	return
	//}
	//t, msg, err := c.ReadMessage()
	//if err != nil {
	//	log.Println("read111:", err, t)
	//	return
	//}
	//log.Printf("receive11: %s\n", msg)

}
