package client_robot

import (
	"encoding/binary"
	"fmt"
	"github.com/gogo/protobuf/proto"
	"io"
	"net"
	"reflect"
	"go_game_server/proto3"
)

// 包头8个字节
const (
	CMDSIZE  = 4 // 消息ID 4个字节
	BODYSIZE = 2 // 包体大小 2个字节
	TEMPSIZE = 2 // 预留2个字节
)

// 接收Length-Type-Value格式的封包流程
func RecvLTVPacket(reader net.Conn, robot *ClientRobot) (cmd proto3.ProtoCmd, tempID uint16, msg interface{}, err error) {
	// 8个字节缓冲区
	headBuff := make([]byte, CMDSIZE+BODYSIZE+TEMPSIZE)

	_, err = io.ReadFull(reader, headBuff)

	if err != nil {
		fmt.Println("------------------err ", err, robot.CurIndex, robot.LordName)
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
		return
	}

	// 最终获取消息体内容
	cmd = proto3.ProtoCmd(msgID)
	pbData := Handler.GetPbData(cmd)
	if pbData != nil {
		msg = reflect.New(reflect.TypeOf(pbData).Elem()).Interface()
		proto.Unmarshal(body, msg.(proto.Message))
	}
	return cmd, tempID, msg, nil
}

// 发送Length-Type-Value格式的封包流程
func SendLTVPacket(writer net.Conn, cmd proto3.ProtoCmd, pbData interface{}) error {
	msgData, _ := proto.Marshal(pbData.(proto.Message))
	pkt := make([]byte, CMDSIZE+BODYSIZE+TEMPSIZE+len(msgData))

	binary.LittleEndian.PutUint32(pkt, uint32(cmd))
	binary.LittleEndian.PutUint16(pkt[CMDSIZE:], uint16(len(msgData)))
	binary.LittleEndian.PutUint16(pkt[CMDSIZE+BODYSIZE:], uint16(255))

	copy(pkt[CMDSIZE+BODYSIZE+TEMPSIZE:], msgData)
	writer.Write(pkt)
	return nil
}
