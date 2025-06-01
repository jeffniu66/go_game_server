package game

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"runtime/debug"
	"strconv"
)

type Room struct {
	roomPid *global.PidObj // 房间协程pid
	// birthList             []int32                // 出生点列表
	frameList             []*proto3.FSPMsg       // 客户端逻辑帧列表
	currentFrame          []*proto3.FSPMsg       // 当前帧列表
	hisFrameList          []*proto3.FSPFrameResp // 历史帧数据
	tickerNum             int                    // ticker次数
	totalScore            int                    // 总积分
	aiScore               int32                  // 机器人完成总积分
	taskMap               map[int32]*Task        // 房间里所有玩家任务数据
	roomId                int32                  // 房间id
	roomInfo              *roomGamingInfo        // 房间内玩家数据
	urgencyTask           *UrgencyTask           // 紧急任务
	skill                 *Skill                 // 技能数据
	closeDoorTime         int64                  // 关门时间 秒
	normalManSkillPoolMap map[int32]int32        // 平民技能池 key: 技能id value: 使用数量
	wolfManSkillPoolMap   map[int32]int32        // 狼人技能池 key: 技能id value: 使用数量
	deathSkillPoolMap     map[int32]int32        // 死亡技能池 key: 技能id value: 使用数量
	isOnlyGame            bool                   // 是否单人模式
	auction               *include.Auction       // 竞拍
	nextUrgencyVoteTime   int32                  // 紧急解冻时间
}

func InitRoomMgr(roomId int32, playerIdList []int32, isOnly bool) *Room {
	frameList := make([]*proto3.FSPMsg, 0)

	roomObj := &Room{
		frameList:  frameList,
		roomId:     roomId,
		roomInfo:   &roomGamingInfo{},
		isOnlyGame: isOnly,
	}
	// 狼人数量以及狼人所在位置
	langNum, _ := tableconfig.ConstsConfigs.GetMatchGameParam()
	matchNum := len(playerIdList)
	langArr := util.RandNoRepeatIntN(0, matchNum, langNum)
	logger.Log.Info("langArr:", langArr)

	// 出身点位置
	birthList := make([]int32, matchNum)
	for i := 0; i < matchNum; i++ {
		birthList[i] = int32(i) + 1
	}
	// fmt.Println("birth list >>>> ", birthList)
	// roomObj.birthList = birthList
	roomObj.roomInfo.InitRoomGamingInfo(playerIdList, birthList, langArr, isOnly)
	roomObj.RecordKillNum(0)
	roomObj.normalManSkillPoolMap = InitSkillPool(constant.NormalManSkillPool)
	roomObj.wolfManSkillPoolMap = InitSkillPool(constant.WolfManSkillPool)
	roomObj.deathSkillPoolMap = InitSkillPool(constant.DeathDropSumPool)

	pidName := "roomPid_" + strconv.Itoa(int(roomId))
	roomPid := global.RegisterPid(pidName, 2048, roomObj)
	roomObj.roomPid = roomPid
	// logger.Log.Info(">>>>>>>room_ticker starting")
	// t := tableconfig.ConstsConfigs.GetFPSNumParam()
	// ticketKey := include.TimerRoomRickerPrefix + strconv.Itoa(int(roomId))
	// roomPid.SendAfter("room_ticker", ticketKey, t, nil)
	return roomObj
}

func (m *Room) Start() {
	logger.Log.Infof("room:%d start ... ", m.roomId)
}

func (m *Room) HandleCall(req global.GenReq) global.Reply {
	return nil
}

func (m *Room) HandleCast(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("room SOCKET_EVENT HandleCast error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			logger.Log.Errorf("room HandleCast panic-------------roomInfo:%v", m)
			debug.PrintStack()
		}
	}()
	switch req.Method {
	case "enterRoom": // 进入房间
		playerIdList := req.MsgData.([]int32)
		// 初始化房间后边游戏开始了
		m.initAiTask()
		m.roomInfo.SetGameWaitInfo()
		for _, userID := range playerIdList {
			m.enterRoom(userID)
		}
		m.roomInfo.RoomStatus = proto3.RoomStatus_wait_game
		t := tableconfig.ConstsConfigs.GetRobotChatTime() * 1000
		m.roomPid.SendAfter("robotChat", "robotChat", t, nil)
	case "gameWait":
		if m.roomInfo.RoomStatus == proto3.RoomStatus_wait_game {
			// 设置默认的用户帧数据
			m.initCurrentFrame()
			waitTime := m.roomInfo.GameWaitTime * 1000
			if switchWait := global.MyConfig.ReadInt32("switch", "switch_game_wait"); switchWait != 1 {
				waitTime = 1 // 1 毫秒 相当于关闭
			}
			m.roomPid.SendAfter("gameBegin", "gameBegin", waitTime, 1)
			auctionDurationTime := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.AuctionDurationTime))
			m.roomPid.SendAfter("auctionEnd", "auctionEnd", auctionDurationTime*1000, nil)
		} else {
			logger.Log.Infof("room status is not wait game, roomStatus:%d", m.roomInfo.RoomStatus)
		}
	case "useLuckyCard":
		msg := req.MsgData.([]int32)
		m.UseLuckyCard(msg[0], msg[1], msg[2])
	case "auction":
		userId := req.MsgData.(int32)
		m.DealAuction(userId)
	case "ninjaAttack":
		req := req.MsgData.(*include.AttackReq)
		if m.CheckAttackFrozen(req.UserId) {
			logger.Log.Warnf("Attack Frozen roomID:%d, attack info:%v", m.roomId, req)
			// 技能冷却
			return
		}
		if m.roomInfo.ChangePlayerRole(req.UserId, proto3.PlayerRoleEnum_wolf_man) && m.roomInfo.CheckPlayerGameStatus(req.UserId, proto3.PlayerGameStatus_normal) &&
			m.roomInfo.CheckPlayerGameStatus(req.SufferUserId, proto3.PlayerGameStatus_normal) && m.roomInfo.ChangePlayerRole(req.SufferUserId, proto3.PlayerRoleEnum_comm_people) {
			m.ninjaAttack(req)
		} else {
			logger.Log.Infof("player is wolf:%v, player status:%v, suffer is people:%v", m.roomInfo.ChangePlayerRole(req.UserId, proto3.PlayerRoleEnum_wolf_man),
				m.roomInfo.CheckPlayerGameStatus(req.SufferUserId, proto3.PlayerGameStatus_normal), m.roomInfo.ChangePlayerRole(req.SufferUserId, proto3.PlayerRoleEnum_comm_people))
		}
	case "startVote":
		if !m.CheckRoomGameStatus(proto3.RoomStatus_gameing) {
			return
		}
		msg := req.MsgData.(*proto3.StartVoteReq)
		b := m.startVote(msg.UserId, msg.SuffererId, msg.VoteType)
		if b {
			// 投票聊天倒计时
			// voteChatTime := 1000 * tableconfig.ConstsConfigs.GetChatVoteTime()
			// 去除投票前聊天功能
			// voteChatTime = 1 // 毫秒
			m.roomPid.SendAfter("voteProChatEnd", "voteProChatEnd", 1, 1)
			t := tableconfig.ConstsConfigs.GetRobotChatTime() * 1000
			m.roomPid.SendAfter("robotChat", "robotChat", t, nil)
			m.RobotStartVote()

		} else {
			logger.Log.Infof("mate:%v, msg:%d isn't match game status:%d", m.roomInfo.AllPlayerGameInfo, msg, proto3.PlayerGameStatus_killed)
		}
	case "playerVote":
		if m.roomInfo.RoomStatus != proto3.RoomStatus_voting {
			return
		}
		voteMsg := req.MsgData.(*include.VoteMsg)
		// 投票算一票
		voteMsg.VoteNum = 1
		if voteMsg != nil {
			if m.roomInfo.VoteInfo.CheckUserVoted(voteMsg.UserId) {
				logger.Log.Warnf("user:%v voted", voteMsg.UserId)
				return
			}
			//// voteMsg.VoteNum 是否使用双倍投票符
			//if m.DealDoubleVote(voteMsg.UserId) {
			//	voteMsg.VoteNum++
			//}
			m.playerVote(voteMsg)
		}
	case "sumVote":
		f := req.MsgData.(int)
		bec := "vote end because player is voted"
		if f == 1 {
			bec = "vote end because vote timeout"
		}
		logger.Log.Infof("room_id:%v %s, voteStep:%v", m.roomId, bec, m.roomInfo.VoteStep)

		if m.roomInfo.VoteStep != proto3.VoteStepEnum_vote_in {
			return
		}
		m.roomInfo.VoteStep = proto3.VoteStepEnum_vote_end
		now := int64(util.UnixTime())
		stepTime := now + int64(tableconfig.ConstsConfigs.GetVoteSumEndTime())
		pbData := &proto3.VoteStepResp{VoteStep: m.roomInfo.VoteStep, ServerCurrentTime: now, StepEndTime: stepTime}
		m.spreadPlayers(proto3.ProtoCmd_CMD_VoteStepResp, pbData)

		m.voteEnd()

		sumTime := 1000 * tableconfig.ConstsConfigs.GetVoteSumEndTime()
		m.roomPid.StopPidTimer("sumVote")
		m.roomPid.SendAfter("voteClose", "voteClose", sumTime, 1)
	case "playerExit":
		userID := req.MsgData.(int32)
		m.playerExit(userID)
	case "playerOffline":
		userID := req.MsgData.(int32)
		m.PlayerOffline(userID)

	case "fpsFrame":
		msg := req.MsgData.(*proto3.FSPC2SDataReq)
		m.frameList = append(m.frameList, msg.Msgs...)
	case "finishTask":
		msgData := req.MsgData.([]int32)
		userId := msgData[0]
		pointId := msgData[1]
		m.FinishTask(userId, pointId)
	case "useItem":
		msgMap := req.MsgData.(map[int32]*proto3.UseItemReq)
		m.UseItem(msgMap)
	case "urgencyTask":
		msgData := req.MsgData.([]int32)
		userId := msgData[0]
		triggerPointId := msgData[1]
		m.TriggerUrgencyTask(userId, triggerPointId)
	case "finishUrgencyTask":
		msgData := req.MsgData.([]int32)
		userId := msgData[0]
		pointId := msgData[1]
		m.FinishUrgencyTask(userId, pointId)
	case "choiceItem":
		msgDataMap := req.MsgData.(map[int32]*proto3.ChoiceItemReq)
		m.ChoiceItem(msgDataMap)
	case "closeDoor":
		msgData := req.MsgData.([]int32)
		userId := msgData[0]
		doorId := msgData[1]
		m.CloseDoor(userId, doorId)
	case "openWind":
		msgData := req.MsgData.([]int32)
		userId := msgData[0]
		doorId := msgData[1]
		m.OpenWind(userId, doorId)
	case "gameEnd":
		if !m.CheckRoomGameStatus(proto3.RoomStatus_game_end) {
			m.roomInfo.RoomStatus = proto3.RoomStatus_game_settle
			logger.Log.Infof("game end roomId:%s", m.roomPid.PidName)
			m.gameEnd()
			m.roomInfo.RoomStatus = proto3.RoomStatus_game_end
			closeTime := 1000 * tableconfig.ConstsConfigs.GetIdValue(constant.RoomCloseTime)
			m.roomPid.SendAfter("roomClose", "roomClose", closeTime, nil)
		}
	case "lampSwitch":
		msg, ok := req.MsgData.(*proto3.LampSwitchReq)
		if !ok {
			logger.Log.Error("lamp switch is error")
			return
		}
		m.roomInfo.LampUserMap[msg.UserId] = int32(msg.Status)
		pbdata := new(proto3.LampSwitchResp)
		pbdata.UserId = msg.UserId
		pbdata.Status = msg.Status
		m.spreadPlayers(proto3.ProtoCmd_CMD_LampSwitchResp, pbdata)
	default:
		fmt.Println("err room handle call method")
	}
}

func (m *Room) HandleInfo(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("room SOCKET_EVENT HandleInfo error:%v, req: %v, stack: %s\n ", err, req, string(debug.Stack()))
			debug.PrintStack()
			logger.Log.Errorf("room Handleinfo panic-------------roomInfo:%v", m)
		}
	}()
	switch req.Method {
	case "room_ticker":
		m.roomTicker()
	case "gameBegin":
		logger.Log.Info(">>>>>>>room_ticker starting")
		t := tableconfig.ConstsConfigs.GetFPSNumParam()
		ticketKey := constant.TimerRoomRickerPrefix + strconv.Itoa(int(m.roomId))
		m.roomPid.SendAfter("room_ticker", ticketKey, t, nil)
		v := req.MsgData.(int)
		if v == 1 && m.roomInfo.RoomStatus == proto3.RoomStatus_wait_game {
			m.roomInfo.RoomStatus = proto3.RoomStatus_gameing

			pbData := &proto3.BeginGameResp{}
			pbData.RoomInfoResp = &m.roomInfo.RoomInfoResp
			pbData.Roommate = m.roomInfo.GetRoommates()
			m.spreadPlayers(proto3.ProtoCmd_CMD_BeginGameResp, pbData)

			for i := 0; i < len(m.roomInfo.AllPlayerGameInfo); i++ {
				v := m.roomInfo.AllPlayerGameInfo[i]
				if v.Player != nil && v.Player.Attr != nil {
					RecordAction(v.Player.Attr.UserID, constant.IDClose)
				}
			}
		}
		logger.Log.Infof("game begin: room:%v, roomMate:%v", m.roomPid.PidName, m.roomInfo)

		lengthTime := global.MyConfig.ReadInt32("room", "game_length_time") * 1000
		m.roomPid.SendAfter("room_stop", "room_stop", int32(lengthTime), nil)
		scheduleKey := constant.TimerAddScorePrefix + strconv.Itoa(int(m.roomId))
		m.roomPid.SendAfter("scheduleAddScore", scheduleKey, 20*1000, nil)
	case "room_stop":
		logger.Log.Infof("roomId:%d is stop because the time is up", m.roomId)
		m.roomInfo.JudgeGameEndByRoleNum()
		m.gameEnd()
		m.roomPid.CastStop()
	case "auctionEnd":
		m.AuctionEnd()
	case "scheduleAddScore":
		m.scheduleAddScore()
	case "validGameEnd":
		m.validGameEnd()
	case "voteProChatEnd":
		m.roomInfo.VoteStep = proto3.VoteStepEnum_vote_in // 投票中
		now := int64(util.UnixTime())
		stepTime := now + int64(m.roomInfo.GameVoteTime)
		m.roomInfo.VoteEndTime = stepTime
		pbData := &proto3.VoteStepResp{VoteStep: m.roomInfo.VoteStep, ServerCurrentTime: now, StepEndTime: stepTime}
		m.spreadPlayers(proto3.ProtoCmd_CMD_VoteStepResp, pbData)
		// 投票倒计时
		m.roomPid.SendAfter("sumVote", "sumVote", 1000*m.roomInfo.GameVoteTime, 1)
	case "sumVote":
		f := req.MsgData.(int)
		bec := "vote end because player is voted"
		if f == 1 {
			bec = "vote end because vote timeout"
		}
		logger.Log.Infof("room_id:%v %s, voteStep:%v", m.roomId, bec, m.roomInfo.VoteStep)

		if m.roomInfo.VoteStep != proto3.VoteStepEnum_vote_in {
			return
		}
		m.roomInfo.VoteStep = proto3.VoteStepEnum_vote_end
		now := int64(util.UnixTime())
		stepTime := now + int64(tableconfig.ConstsConfigs.GetVoteSumEndTime())
		m.roomInfo.VoteSumEndTime = stepTime
		pbData := &proto3.VoteStepResp{VoteStep: m.roomInfo.VoteStep, ServerCurrentTime: now, StepEndTime: stepTime}
		m.spreadPlayers(proto3.ProtoCmd_CMD_VoteStepResp, pbData)

		m.voteEnd()

		sumTime := 1000 * tableconfig.ConstsConfigs.GetVoteSumEndTime()
		m.roomPid.SendAfter("voteClose", "voteClose", sumTime, 1)
	case "voteClose":
		m.spreadPlayers(proto3.ProtoCmd_CMD_VoteCloseResp, &proto3.VoteCloseResp{Status: proto3.CommonStatusEnum_true})
		m.voteClose()
	case "roomClose":
		m.roomPid.CastStop()
	case "robotVote":
		robotID := req.MsgData.(int32)
		m.RobotVote(robotID)

	case "robotChat":
		if m.roomInfo.RoomStatus == proto3.RoomStatus_wait_game || m.roomInfo.RoomStatus == proto3.RoomStatus_voting {
			m.robotChat()
		} else {
			m.roomPid.StopPidTimer("robotChat")
		}
	default:
		fmt.Println("err room handle call method")
	}
}

func (m *Room) Terminate() {
	logger.Log.Infof("room pid terminate roomPidName:%v", m.roomPid.PidName)
	// 清理玩家
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		v.PlayerAutoOut = autoOutGameEnd
		// 玩家可能主动退出，也可能游戏结束了房间主动清理
		p := global.GloInstance.GetPlayer(v.UserId)
		if p != nil {
			player := p.(*Player)
			player.GameEnd(m.roomId)
		}
	}
	// 清理房间
	// m.birthList = []int32{}
	m.currentFrame = []*proto3.FSPMsg{}
	m.hisFrameList = []*proto3.FSPFrameResp{}
	m.frameList = []*proto3.FSPMsg{}
	m.roomInfo = nil
	m.taskMap = nil
	m.totalScore = 0
	m.aiScore = 0
	m.tickerNum = 0
	m.urgencyTask = nil
	m.skill = nil
	m.closeDoorTime = 0
	m.normalManSkillPoolMap = nil
	m.wolfManSkillPoolMap = nil
	m.deathSkillPoolMap = nil
	m.auction = nil

	// 清理线程
	m.roomPid.StopAllTimer()
	m.roomPid = nil

	// 修改房间可用
	global.GloInstance.ChangeRoomIdUnused(m.roomId)

	GMSkillId = 0 // 每一局结束后恢复
}

func (m *Room) initCurrentFrame() {
	m.hisFrameList = []*proto3.FSPFrameResp{}
	m.defaultFram()
}

func (m *Room) defaultFram() {
	tmpFms := make([]*proto3.FSPMsg, 0)
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		tmp := &proto3.FSPMsg{
			UId:  uint32(v.UserId),
			Args: &proto3.FSPCmdArgs{},
		}
		logger.Log.Debugf("init current frame:%v", tmp)
		tmpFms = append(tmpFms, tmp)
	}
	m.currentFrame = tmpFms
}

func (r *Room) CheckRoomGameStatus(status proto3.RoomStatus) bool {
	if r.roomInfo == nil {
		return false
	}
	return status == r.roomInfo.RoomStatus
}

func (r *Room) IsInRoom(userId int32) bool {
	return r.roomInfo.IsInRoom(userId)
}

// 进入房间 清理旧的playerPid
func (r *Room) PlayerEnterRoom(player *Player) {
	player.Room = r
	player.RoomPid = r.roomPid
	player.Attr.RoomID = r.roomId
	userInfo := r.roomInfo.GetPlayerGameInfo(player.Attr.UserID)
	userInfo.Player = player
	// for i := range r.roomInfo.AllPlayerGameInfo {
	// 	v := r.roomInfo.AllPlayerGameInfo[i]
	// 	if v.UserId == player.Attr.UserID {
	// 		v.Player = player
	// 	}
	// }
}

func (r *Room) SendHisFSPFrame(player *Player) {
	logger.Log.Infof("--------start----SendHisFSPFrame playerID:%d", player.Attr.UserID)
	cmd := proto3.ProtoCmd_CMD_HisFSPFrameResp
	pbData := &proto3.HisFSPFrameResp{}
	count := 0
	for i := 0; i < len(r.hisFrameList); i++ {
		pbData.HisFSPFrameList = append(pbData.HisFSPFrameList, r.hisFrameList[i])
		count++
		if i%100 == 99 || i >= len(r.hisFrameList)-1 {
			// player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
			r.sendMsgPlayer(player, cmd, pbData)
			logger.Log.Infof("--------sending----SendHisFSPFrame startFmsId:%v, count:%d endFmsId:%v", pbData.HisFSPFrameList[0].FrameId, count, pbData.HisFSPFrameList[count-1].FrameId)
			pbData = &proto3.HisFSPFrameResp{} // 注意异步发送，可能会清理指针内容，需要重新开辟空间
			count = 0
		}
	}
	logger.Log.Infof("--------end----SendHisFSPFrame playerID:%d, count:%v", player.Attr.UserID, len(r.hisFrameList))
}

func (r *Room) GetEnterRoomResp(player *Player) *proto3.EnterRoomResp {
	userID := player.Attr.UserID
	ret := &proto3.EnterRoomResp{}
	ret.Roommate = r.roomInfo.GetProto3Roommate()
	task := r.taskMap[player.Attr.UserID]
	assignTaskPoint := ""
	pointInfo := ""
	var skill int32
	if task != nil {
		assignTaskPoint = task.AssignTaskPoint
		pointInfo = task.PointInfo
		if task.Skill != nil {
			skill = task.Skill.Id
		}
	}
	ret.TaskInfo = &proto3.TaskInfo{TaskPoint: GetAssignTaskPoints(assignTaskPoint), FinishPoint: GetFinishPoints(pointInfo),
		TotalScore: int32(r.totalScore), UrgencyTaskPoint: GetAssignUrgencyTaskPoints(), UrgencyTaskResp: r.GetUrgencyTaskInfoResp(), Skill: skill}
	ret.RoomInfoResp = r.GetRoomInfoResp()
	ret.VoteResult = r.roomInfo.VoteInfo.GetVoteResult()
	ret.SkillInfo = &proto3.SkillInfo{BoomUserId: r.GetBoomUserId(player.Attr.UserID)}
	ret.RoomId = r.roomId
	ret.UrgencyVoteNum = r.roomInfo.GetUrgencyVoteNum(userID)
	ret.UrgencyVoteTime = r.nextUrgencyVoteTime
	ret.AiTask = r.roomInfo.AiTask
	ret.RoomInfoResp.IsOnly = r.roomInfo.IsOnly
	ret.WolfManCd = r.roomInfo.GetWolfManCd()
	ret.LampList = r.roomInfo.GetLampList()
	return ret
}
