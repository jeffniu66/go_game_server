package handler

import (
	"go_game_server/proto3"
	"go_game_server/server/db"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/tableconfig"
	"strings"
)

func init() {
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_GameChatReq, &proto3.GameChatReq{}, handleGameChat)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_WaitGameChatReq, &proto3.WaitGameChatReq{}, handleWaitGameChat)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_HomeChatReq, &proto3.HomeChatReq{}, handleHomeChat)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_GmChatReq, &proto3.GmChatReq{}, handleGmChat)
}

func handleGameChat(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.GameChatReq)
	if len(msg.ChatData) <= 0 && msg.QuickId < 0 {
		return nil
	}
	room := player.IsInRoom()
	if room != nil {
		if player.CheckVoting() && room.CheckPlayerGameStatus(player.Attr.UserID, proto3.PlayerGameStatus_normal) {
			room.AddChatNum(player.Attr.UserID, msg.ChatData)
			msg.ChatData = tableconfig.SenWordConfigs.ReplaceSenWord(msg.ChatData)
			room.SendChat(player.Attr.UserID, msg.ChatData, msg.QuickType, msg.QuickId, nil)
		}
	}
	return nil
}

func handleWaitGameChat(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.WaitGameChatReq)
	if len(msg.ChatData) <= 0 && msg.QuickId < 0 {
		return nil
	}
	room := player.IsInRoom()
	if room != nil {
		if player.CheckWaiting() && room.CheckPlayerGameStatus(player.Attr.UserID, proto3.PlayerGameStatus_normal) {
			msg.ChatData = tableconfig.SenWordConfigs.ReplaceSenWord(msg.ChatData)
			room.SendWaitChat(player.Attr.UserID, msg.ChatData, msg.QuickType, msg.QuickId)
		}
	}
	return nil
}

// 主页聊天
func handleHomeChat(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.HomeChatReq)
	if len(msg.ChatData) <= 0 {
		return nil
	}
	if strings.Contains(msg.ChatData, "gm") {
		if !global.OrderConfig.OpenGM {
			return nil
		}
		GmCommand(player, msg.ChatData, nil)
		game.SpreadWorldPlayer(player, proto3.ChatTypeEnum_chat_merge, msg.ChatData, 0, 0, nil)
	} else {
		attr := player.Attr
		msg.ChatData = tableconfig.SenWordConfigs.ReplaceSenWord(msg.ChatData)
		go db.SaveWorldChat(attr.UserID, attr.UserPhoto, attr.Username, msg.ChatData)
		game.SpreadWorldPlayer(player, proto3.ChatTypeEnum_chat_merge, msg.ChatData, 0, 0, nil)
	}
	return nil
}

func handleGmChat(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.GmChatReq)
	if len(msg.Msg) <= 0 {
		return nil
	}
	GmCommand(player, msg.Msg, nil)
	return nil
}
