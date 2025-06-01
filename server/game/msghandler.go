package game

import (
	"go_game_server/proto3"
)

var Handler *MsgData

type Message struct {
	Cmd    proto3.ProtoCmd
	PbData interface{}
}

type callBackFunc func(interface{}, *Player) interface{}

type MsgData struct {
	msgInfo map[proto3.ProtoCmd]*MsgInfo
}

type MsgInfo struct {
	pb       interface{}
	callback callBackFunc
}

func init() {
	Handler = newRegister()
}

func newRegister() *MsgData {
	msgData := new(MsgData)
	msgData.msgInfo = make(map[proto3.ProtoCmd]*MsgInfo)
	return msgData
}

func (msgData *MsgData) RegistHandler(cmd proto3.ProtoCmd, pbData interface{}, callback callBackFunc) {
	msgInfo := &MsgInfo{pb: pbData, callback: callback}
	msgData.msgInfo[cmd] = msgInfo
}

func (msgData *MsgData) GetPbData(cmd proto3.ProtoCmd) (pbData interface{}) {
	if msgInfo, ok := msgData.msgInfo[cmd]; ok {
		pbData = msgInfo.pb
	}
	return
}

// 可以实现不同的callback
func (msgData *MsgData) Callback(cmd proto3.ProtoCmd, pbData interface{}, player *Player) {
	msgInfo := msgData.msgInfo[cmd]
	callBack := msgInfo.callback
	callBack(pbData, player)
}
