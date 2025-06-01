package main

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/network"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	c, _, err := websocket.DefaultDialer.Dial("ws://localhost:8005/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		cmd0 := proto3.ProtoCmd_CMD_LoginReq
		pbData0 := &proto3.LoginReq{}
		// pbData0.AcctName = "hello"

		network.SendWsMessage(c, cmd0, pbData0)

		time.Sleep(time.Second)

		cmd := proto3.ProtoCmd_CMD_FrameSyncReq
		pbData := &proto3.FSPC2SDataReq{}
		fspMsg := &proto3.FSPMsg{}
		fspMsg.UId = 1
		pbData.Msgs = []*proto3.FSPMsg{fspMsg}

		network.SendWsMessage(c, cmd, pbData)

		time.Sleep(1000 * time.Second)
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
