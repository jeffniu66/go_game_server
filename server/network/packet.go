package network

import (
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"go_game_server/proto3"
	"go_game_server/server/game"
	"go_game_server/server/logger"
	"io"
	"net"
	"reflect"
	"runtime/debug"
)

// 包头8个字节
const (
	CMDSIZE  = 4 // 消息ID 4个字节
	BODYSIZE = 2 // 包体大小 2个字节
	TEMPSIZE = 2 // 预留2个字节
)

// 接收Length-Type-Value格式的封包流程
func RecvLTVPacket(reader net.Conn) (cmd proto3.ProtoCmd, tempID uint16, msg interface{}, err error) {
	// 8个字节缓冲区
	headBuff := make([]byte, CMDSIZE+BODYSIZE+TEMPSIZE)
	//_, err = io.ReadFull(reader, headBuff)

	_, err = io.ReadFull(reader, headBuff)

	if err != nil {
		fmt.Println("------------------err ", err)
		return // 这里并非跳出循环，而是返回值
	}

	// 先读取4字节的消息id
	msgID := int32(binary.LittleEndian.Uint32(headBuff[0:CMDSIZE]))
	// 分配包体缓冲区,再读取2字节的包体长度
	bodyLen := binary.LittleEndian.Uint16(headBuff[CMDSIZE:])
	// 分配包体缓冲区,再读取2字节的预留id
	tempID = binary.LittleEndian.Uint16(headBuff[CMDSIZE+BODYSIZE:])

	// 分配包体大小
	body := make([]byte, bodyLen)
	// 这里的ReadFull，是指将reader字节流，读取到缓冲区buff，而且必须按照包体读满！否则会一直在循环等待字节进来，所以我们必须设置心跳
	// 以免它长时间在等待解包，导致socket一直被占用，close_wait
	_, err = io.ReadFull(reader, body)

	if err != nil {
		//reader.Close()		// client disconnect
		fmt.Println("------------------client error ", err)
		return
	}

	// 最终获取消息体内容
	cmd = proto3.ProtoCmd(msgID)
	pbData := game.Handler.GetPbData(cmd)
	if pbData != nil {
		msg = reflect.New(reflect.TypeOf(pbData).Elem()).Interface()
		err = proto.Unmarshal(body, msg.(proto.Message))
	}
	//logger.Log.Infof("receive cmd = %v", cmd)

	return cmd, tempID, msg, err
}

// 发送Length-Type-Value格式的封包流程
func SendLTVPacket(writer net.Conn, cmd proto3.ProtoCmd, pbData interface{}) error {
	if err := recover(); err != nil {
		logger.Log.Errorf("error SendLTVPacket cmd = %v, len = %v \n", cmd, pbData) // 这里的err其实就是panic传入的内容，55
		debug.PrintStack()
	}
	msgData, err := proto.Marshal(pbData.(proto.Message))
	if len(msgData) > 65535 {
		logger.Log.Errorf("SendLTVPacket len too long cmd = %v, len = %v \n", cmd, len(msgData))
		cmd = proto3.ProtoCmd_CMD_ErrResp
		//pbData = &proto3.ErrResp{ErrCode: proto3.ErrEnum_Error_Packet_limit, ErrMsg: "packet too big"}
		_ = SendLTVPacket(writer, cmd, pbData)
		return nil
	}
	if err != nil {
		logger.Log.Errorf("eSendLTVPacket encode error cmd = %v, err = %v \n", cmd, err) // 这里的err其实就是panic传入的内容，55
		return nil
	}
	pkt := make([]byte, CMDSIZE+BODYSIZE+TEMPSIZE+len(msgData))

	binary.LittleEndian.PutUint32(pkt, uint32(cmd))
	binary.LittleEndian.PutUint16(pkt[CMDSIZE:], uint16(len(msgData)))
	binary.LittleEndian.PutUint16(pkt[CMDSIZE+BODYSIZE:], uint16(255))

	copy(pkt[CMDSIZE+BODYSIZE+TEMPSIZE:], msgData)
	n, err := writer.Write(pkt)
	if err != nil {
		fmt.Printf("send cmd = %v, len = %v\n, pbData = %v\n", cmd, len(msgData), pbData)
		logger.Log.Error(cmd, n, err)
	}

	return nil
}
