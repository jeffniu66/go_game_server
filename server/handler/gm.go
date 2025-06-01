package handler

import (
	"encoding/json"
	"go_game_server/proto3"
	"go_game_server/server/db"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/util"
	"strings"

	"github.com/gorilla/websocket"
)

var Clients = make(map[*websocket.Conn]bool) // connected clients

func GmCommand(p *game.Player, msg string, httpRes *string) {
	//if !global.OrderConfig.CanUseGm(p.Attr.Username) {
	//	logger.Log.Warnln("can't use gm, because no in gm order list!", p.Attr.Username)
	//	return
	//}
	cmd := proto3.ProtoCmd_CMD_GmChatResp
	pbData := &proto3.GmChatResp{}
	retMsg := &game.Message{Cmd: cmd, PbData: pbData}
	chatMsg := strings.Split(msg, "^")
	switch chatMsg[0] {
	case include.GMAllRoom:
		rooms := global.GloInstance.GetUsedRoomFaceList()
		num := len(rooms)
		pbData.Response = util.ToStr(int32(num)) + ","
		for _, v := range rooms {
			ro := v.(*game.Room)
			pbData.Response += util.ToStr(ro.GetRoomID()) + ","
		}
	case include.GMGetRoom:
		roomID := util.ToInt(chatMsg[1])
		r := global.GloInstance.GetRoom(roomID)
		if r == nil {
			return
		}
		room, ok := r.(*game.Room)
		if !ok {
			return
		}
		data, err := json.Marshal(room.GetRoomInfoResp())
		if err != nil {
			logger.Log.Errorf("gm roomInfo failed, roomID:%d err:%s", roomID, err.Error())
			return
		}
		da, err := json.Marshal(room.GetRoommate())
		if err != nil {
			logger.Log.Errorf("gm roomInfo failed, roomID:%d err:%s", roomID, err.Error())
			return
		}
		pbData.Response += string(data) + "\n" + string(da)
	case include.GMAddItem:
		p.AddItems(chatMsg[1])
	case include.GMAddSkill:
		game.GMSkillId = util.ToInt(chatMsg[1])
	case include.GMGetOnlineNum:
		if util.ToInt(chatMsg[1]) == 1 {
			num := len(Clients)
			pbData.Response = util.ToStr(int32(num))
		} else {
			pbData.Response = util.ToStr(db.SelectOnlineNum())
		}
	case include.GMGetRegisterNum:
		pbData.Response = "wx: " + util.ToStr(db.SelectRegisterNum("wx", chatMsg[1], chatMsg[2])) +
			"\n" + "qq: " + util.ToStr(db.SelectRegisterNum("qq", chatMsg[1], chatMsg[2])) +
			"\n" + "oppo: " + util.ToStr(db.SelectRegisterNum("oppo", chatMsg[1], chatMsg[2]))
	case include.GMUptRank:
		game.InitTopBoard()
	default:
	}
	if p != nil && p.Pid != nil {
		p.SendMessage(retMsg)
	}

	if httpRes != nil {
		*httpRes = pbData.Response
	}
}
