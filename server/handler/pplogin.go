package handler

import (
	"go_game_server/proto3"
	"go_game_server/server/db"
	"go_game_server/server/game"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
)

func init() {
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_HeartBeatReq, &proto3.HeartBeatReq{}, handleHeartbeat)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_LoginReq, &proto3.LoginReq{}, handleLogin)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_RegisterReq, &proto3.RegisterReq{}, handleRegister)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_RandNameReq, &proto3.RandNameReq{}, handleRandName)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_CreateUserReq, &proto3.CreateUserReq{}, handleCreateUser)
}

func handleHeartbeat(req interface{}, player *game.Player) interface{} {
	cmd := proto3.ProtoCmd_CMD_HeartBeatResp
	pbData := &proto3.HeartBeatResp{ServerSec: util.UnixTime()}
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}

func handleLogin(req interface{}, player *game.Player) interface{} {
	logger.Log.Infof("login success player_id = %v, login_time = %v, logout_time = %v, LoginRespStatus = %v ",
		player.Attr.UserID, player.Attr.LoginTime, player.Attr.RegistTime, player.LoginRespStatus)

	playerAttr := player.GetProto3PlayerAttr()

	room := player.IsInRoom()
	pbData := &proto3.LoginResp{ServerSec: util.UnixTime(), PlayerAttr: playerAttr, Ret: player.LoginRespStatus, SessionKey: player.SessionKey}
	pbData.ItemResp = &proto3.ItemResp{Items: game.GetItems(player.Attr.UserID)}
	pbData.Guides = player.GetGuides()
	pbData.RedDot = player.GetRedData()
	if room == nil {
		pbData.IsRoom = proto3.CommonStatusEnum_false
		player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_LoginResp, PbData: pbData})
		// 新手礼包
		go player.GetFreshGiftStep()
		go player.CheckUpTop()
		// 记录登录数据
		game.RecordLogin(player)
	} else {
		pbData.IsRoom = proto3.CommonStatusEnum_true
		player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_LoginResp, PbData: pbData})

		// 如果已经在房间，重新进入房间
		room.PlayerEnterRoom(player)
		cmd := proto3.ProtoCmd_CMD_EnterRoomResp
		pbData := room.GetEnterRoomResp(player)
		player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
		if room.CheckRoomGameStatus(proto3.RoomStatus_voting) || room.CheckRoomGameStatus(proto3.RoomStatus_gameing) {
			room.SendHisFSPFrame(player)
		}
	}
	logger.Log.Infof("player is in room, bool:%v, playerAttr:%v", room != nil, player.Attr)
	return nil
}

func handleRegister(req interface{}, tourist *game.Player) interface{} {
	msg := req.(*proto3.RegisterReq)

	cmd := proto3.ProtoCmd_CMD_RegisterResp
	pbData := &proto3.RegisterResp{}
	if msg.RegisterType == proto3.RegisterTypeEnum_register_nomal {
		if db.JudgeAcctName(msg.AcctName) {
			tourist.ErrorResponse(proto3.ErrEnum_Error_ExistName, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_ExistName)])
			pbData.Status = proto3.CommonStatusEnum_false
		} else {
			acct, b := tourist.Register(msg.AcctName, msg.Pw)
			if !b {
				tourist.ErrorResponse(proto3.ErrEnum_Error_RegisterFailed, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_RegisterFailed)])
				return nil
			}
			pbData.Acct = acct.AcctName
			pbData.Status = proto3.CommonStatusEnum_true
		}
	} else if msg.RegisterType == proto3.RegisterTypeEnum_register_quick {
		var acctName, pw string = "", ""
		for {
			acctName = util.RandStr(8)
			if !db.JudgeAcctName(acctName) {
				break
			}
		}
		for {
			pw = util.RandStr(8)
			if !db.JudgeAcctName(pw) {
				break
			}
		}
		acct, b := tourist.Register(acctName, pw)
		if !b {
			tourist.ErrorResponse(proto3.ErrEnum_Error_RegisterFailed, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_RegisterFailed)])
			return nil
		}
		pbData.Acct = acct.AcctName
		pbData.Pw = pw
		pbData.Status = proto3.CommonStatusEnum_true
	}
	tourist.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}

func handleRandName(req interface{}, tourist *game.Player) interface{} {
	randName, sex, nameIndex := game.GetRankName(tableconfig.NameZhConfigs)
	cmd := proto3.ProtoCmd_CMD_RandNameResp
	pbData := &proto3.RandNameResp{}
	pbData.Username = randName
	pbData.Sex = sex
	pbData.NameIndex = nameIndex
	tourist.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}

func handleCreateUser(req interface{}, tourist *game.Player) interface{} {
	//msg := req.(*proto3.CreateUserReq)
	//attr := new(db.Attr)
	//nameCount := db.CountUsername(msg.NameIndex)
	//if nameCount > 0 {
	//	msg.Username += util.ToStr(nameCount + 1)
	//}
	//attr.InitData(msg.Username, msg.AcctName, msg.NameIndex, msg.Sex)
	//cmd := proto3.ProtoCmd_CMD_CreateUserResp
	//pbData := &proto3.CreateUserResp{}
	//pbData.Username = msg.Username
	//tourist.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}
