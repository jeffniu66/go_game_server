package handler

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

func init() {
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_FrameSyncReq, &proto3.FSPC2SDataReq{}, dealFrameData)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_EnterMatchReq, &proto3.EnterMatchReq{}, handleEnterMatch)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_FinishMissionReq, &proto3.FinishMissionReq{}, handleFinishMission)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_StartVoteReq, &proto3.StartVoteReq{}, handleStartVote)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_PlayerVoteReq, &proto3.PlayerVoteReq{}, handlePlayerVote)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_NinjaAttackReq, &proto3.NinjaAttackReq{}, handleNinjaAttack)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UseItemReq, &proto3.UseItemReq{}, handleUseItem)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_PlayerExitReq, &proto3.PlayerExitReq{}, handlePlayerExit)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UrgencyTaskReq, &proto3.UrgencyTaskReq{}, handleUrgencyTask)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_FinishUrgencyTaskReq, &proto3.FinishUrgencyTaskReq{}, handleFinishUrgencyTask)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_ChoiceItemReq, &proto3.ChoiceItemReq{}, handleChoiceItem)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_CloseDoorReq, &proto3.CloseDoorReq{}, handleCloseDoor)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_AuctionReq, &proto3.AuctionReq{}, handleAuction)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UseLuckyCardReq, &proto3.UseLuckyCardReq{}, handleUseLuckyCard)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_OpenWindReq, &proto3.OpenWindReq{}, handleOpenWind)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_LampSwitchReq, &proto3.LampSwitchReq{}, handleLampSwitch)
}

func handleEnterMatch(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.EnterMatchReq)
	if msg.Key == proto3.MatchEnum_enter_match {
		if room := player.IsInRoom(); room != nil {
			room.PlayerEnterRoom(player)
			cmd := proto3.ProtoCmd_CMD_EnterRoomResp
			pbData := room.GetEnterRoomResp(player)
			player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
			if room.CheckRoomGameStatus(proto3.RoomStatus_voting) || room.CheckRoomGameStatus(proto3.RoomStatus_gameing) {
				room.SendHisFSPFrame(player)
			}
			return nil
		}
		if !global.MemIsHealthy() || !global.CpuIsHealthy() {
			cmd := proto3.ProtoCmd_CMD_EnterMatchErrorResp
			pbData := &proto3.EnterMatchErrorResp{
				ErrCode: proto3.ErrEnum_Error_System_busy,
				ErrMsg:  proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_System_busy)],
			}
			player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
			return nil
		}
		if !global.RoomNumIsHealthy() {
			// todo
		}
		if msg.IsOnly == proto3.CommonStatusEnum_true {
			game.MatchMgrPid.Cast("onlyGame", player)
		} else {
			game.MatchMgrPid.Cast("enterMatch", player)

			game.RecordAction(player.Attr.UserID, constant.ClickMatch)
		}
	} else if msg.Key == proto3.MatchEnum_quit_match {
		game.MatchMgrPid.Cast("quitMatch", player)

		game.RecordAction(player.Attr.UserID, constant.ClickCancelMatch)
	}
	return nil
}

func enterRoom(req interface{}, player *game.Player) interface{} {
	cmd := proto3.ProtoCmd_CMD_HeartBeatResp
	pbData := &proto3.HeartBeatResp{ServerSec: util.UnixTime()}
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}

func dealFrameData(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.FSPC2SDataReq)
	if player.CheckWaiting() || player.CheckGaming() || player.CheckVoting() {
		logger.Log.Debugf("接收到user:%v, 的帧数据 msg.Msgs: %v", player.Attr, msg.Msgs)
		player.RoomPid.Cast("fpsFrame", req)
	} else {
		logger.Log.Debugf("游戏状态不匹配，不能进行同步, waiting:%v, gaming:%v, voting:%v", player.CheckWaiting(), player.CheckGaming(), player.CheckVoting())
	}
	return nil
}

func handleFinishMission(req interface{}, player *game.Player) interface{} {
	if player.CheckGaming() || player.CheckVoting() {
		msg := req.(*proto3.FinishMissionReq)
		player.RoomPid.Cast("finishTask", []int32{player.Attr.UserID, msg.PointId})

		game.RecordAction(player.Attr.UserID, constant.DoTask)
	} else {
		logger.Log.Errorln("handleFinishMission 不在游戏或投票状态")
		if player.Room == nil || player.RoomPid == nil {
			logger.Log.Info("handleFinishMission 房间为空!")
		} else {
			logger.Log.Infof("handleFinishMission 玩家的游戏状态: %d", player.Room.GetRoomInfoResp().RoomStatus)
		}
	}
	return nil
}

func handleStartVote(req interface{}, player *game.Player) interface{} {
	if player.CheckGaming() {
		msg := req.(*proto3.StartVoteReq)
		if msg.UserId <= 0 { // 兼容客户端不传userID
			msg.UserId = player.Attr.UserID
		}
		logger.Log.Info("initiate vote!", msg)
		player.RoomPid.Cast("startVote", msg)

		game.RecordAction(player.Attr.UserID, constant.CallPolice)
	}
	return nil
}

func handlePlayerVote(req interface{}, player *game.Player) interface{} {
	if player.CheckVoting() && player.CheckPlayerGameStatus(proto3.PlayerGameStatus_normal) {
		logger.Log.Info("player vote!")
		msg := req.(*proto3.PlayerVoteReq)
		voteMsg := &include.VoteMsg{UserId: player.Attr.UserID, VoteUserId: msg.VoteUserId, Status: msg.Status}
		player.RoomPid.Cast("playerVote", voteMsg)

	}
	return nil
}

func handleNinjaAttack(req interface{}, player *game.Player) interface{} {
	if player.CheckGaming() {
		msg := req.(*proto3.NinjaAttackReq)
		var userID int32
		if msg.UserId <= 0 {
			userID = player.Attr.UserID // 兼容旧版本
		} else {
			userID = msg.UserId
		}
		// Todo check userId for security
		attackReq := &include.AttackReq{
			SocketUser:   player.Attr.UserID,
			UserId:       userID,
			SufferUserId: msg.SuffererId,
		}
		logger.Log.Infof("ninja attack! reqMsg:%v", msg)
		player.RoomPid.Cast("ninjaAttack", attackReq)

		game.RecordAction(msg.UserId, constant.KillPeople)
		game.RecordAction(msg.SuffererId, constant.Killed)
	}
	return nil
}

func handleUseItem(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UseItemReq)
	msgMap := make(map[int32]*proto3.UseItemReq)
	msgMap[player.Attr.UserID] = msg
	player.RoomPid.Cast("useItem", msgMap)
	return nil
}

func handlePlayerExit(req interface{}, player *game.Player) interface{} {
	player.Pid.Cast("quitRoom", player.Attr.UserID)
	return nil
}

func handleUrgencyTask(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UrgencyTaskReq)
	player.RoomPid.Cast("urgencyTask", []int32{player.Attr.UserID, msg.TriggerPoint})
	return nil
}

func handleFinishUrgencyTask(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.FinishUrgencyTaskReq)
	userID := msg.UserId
	if userID <= 0 {
		userID = player.Attr.UserID
	}
	player.RoomPid.Cast("finishUrgencyTask", []int32{userID, msg.PointId})
	return nil
}

func handleChoiceItem(req interface{}, player *game.Player) interface{} {
	msgMap := make(map[int32]*proto3.ChoiceItemReq)
	msgMap[player.Attr.UserID] = req.(*proto3.ChoiceItemReq)
	player.RoomPid.Cast("choiceItem", msgMap)
	return nil
}

func handleCloseDoor(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.CloseDoorReq)
	player.RoomPid.Cast("closeDoor", []int32{player.Attr.UserID, msg.DoorId})
	return nil
}

func handleAuction(req interface{}, player *game.Player) interface{} {
	if player.CheckWaiting() {
		player.RoomPid.Cast("auction", player.Attr.UserID)

		game.RecordAction(player.Attr.UserID, constant.ClickAuction)
	}
	return nil
}

func handleUseLuckyCard(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UseLuckyCardReq)
	player.RoomPid.Cast("useLuckyCard", []int32{player.Attr.UserID, msg.ItemId, msg.ItemNum})

	game.RecordAction(player.Attr.UserID, constant.ClickLuckyCard)
	return nil
}

func handleOpenWind(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.OpenWindReq)
	player.RoomPid.Cast("openWind", []int32{player.Attr.UserID, msg.WindId})
	return nil
}

func handleLampSwitch(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.LampSwitchReq)
	player.RoomPid.Cast("lampSwitch", msg)
	return nil
}
