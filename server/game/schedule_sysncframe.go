package game

import (
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"github.com/gorilla/websocket"
	"go_game_server/proto3"
	"log"
	"time"
)

var ConnMap = make(map[string]*websocket.Conn) // 玩家连接集合

var FrameList []*proto3.FSPMsg
var TickerNum int32
var frameIndex int64

func SchedulePushClient() {
	ticker := time.NewTicker(50 * time.Millisecond)

	for {
		<-ticker.C
		fmt.Println("=================ticker=================")
		// 从列表里取数据
		frameData := FrameList[frameIndex:]
		frameIndex = int64(len(FrameList))
		fmt.Println("===========frameData=============", frameData)
		if frameData == nil {
			fmt.Println("没有新增帧数据")
			continue
		}
		fmt.Println("ConnMap len: ", len(ConnMap))

		TickerNum++

		for _, con := range ConnMap {
			//binary.LittleEndian.PutUint16(msg[:2], uint16(proto3.ProtoCmd_CMD_FrameSyncResp))

			//con.WriteMessage(websocket.BinaryMessage, msg)

			//cmd := proto3.ProtoCmd_CMD_FrameSyncResp
			//pbData := &proto3.FSPS2CData{}
			//pbData.Type = proto3.FSPS2CDataType_TYPE_FRAME
			//fspFrame := &proto3.FSPFrame{}
			//fspFrame.FrameId = frameId
			//
			//fspc2sData := &proto3.FSPC2SData{}
			//proto.Unmarshal(msg[2:], fspc2sData)
			//fspFrame.Msgs = fspc2sData.Msgs
			//
			//pbData.Msgs = []*proto3.FSPFrame{fspFrame}

			cmd := proto3.ProtoCmd_CMD_FrameSyncResp
			pbData := &proto3.FSPFrameResp{}
			pbData.FrameId = uint32(TickerNum)
			pbData.Msgs = frameData

			sendWsMessage(con, cmd, pbData)
		}
	}
}

func sendWsMessage(ws *websocket.Conn, cmd proto3.ProtoCmd, pbData interface{}) error {
	msgData, err := proto.Marshal(pbData.(proto.Message))
	if err != nil {
		fmt.Println("proto marshal err : ", err)
	}
	msg := make([]byte, len(msgData)+2)
	fmt.Println("==============cmd=========", cmd)
	binary.LittleEndian.PutUint16(msg, uint16(cmd))
	copy(msg[2:], msgData)
	fmt.Println("===========msg=============", msg)
	err = ws.WriteMessage(websocket.BinaryMessage, msg)
	if err != nil {
		log.Println("websocket send message err :", err)
		return err
	}
	return nil
}
