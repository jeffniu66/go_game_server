package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"strconv"
	"strings"
	"time"
)

func (r *Room) GetRoomInfoResp() *proto3.RoomInfoResp {
	r.roomInfo.ServerCurrentTime = int64(util.UnixTime())
	return &r.roomInfo.RoomInfoResp
}
func (r *Room) GetRoommate() []*PlayerGameInfo {
	return r.roomInfo.AllPlayerGameInfo
}

func (r *Room) UpdateCurrentFrame(fspMsg *proto3.FSPMsg) {
	for i := 0; i < len(r.currentFrame); i++ {
		if r.currentFrame[i].UId == fspMsg.UId {
			*r.currentFrame[i] = *fspMsg
			logger.Log.Debugf("currentFram:%v, fmsMsg:%v", *r.currentFrame[i], *fspMsg)
		}
	}
}

func (m *Room) initAiTask() {
	m.roomInfo.AiTask = make([]*proto3.TaskInfo, 0)
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		if v.IsRobot == 1 {
			positions := m.RankTaskIds(v.UserId)
			tmp := &proto3.TaskInfo{TaskPoint: positions}
			tmp.UserId = v.UserId
			m.roomInfo.AiTask = append(m.roomInfo.AiTask, tmp)
		}
	}
	return
}

func (m *Room) enterRoom(userID int32) {
	// 保留玩家出身点信息，防止匹配断线后重连无数据
	positions := m.RankTaskIds(userID)
	taskInfo := &proto3.TaskInfo{TaskPoint: positions, UrgencyTaskPoint: GetAssignUrgencyTaskPoints()}
	taskInfo.UserId = userID

	p := global.GloInstance.GetPlayer(userID) // TODO 判断失效的socket
	if p != nil {
		player := p.(*Player)
		m.PlayerEnterRoom(player)
		roommate := m.roomInfo.GetProto3Roommate()
		m.DealSkillTab(player, roommate) // 扣除玩家身上的2号或3号技能选项卡
		pbData := &proto3.EnterRoomResp{Roommate: roommate, TaskInfo: taskInfo, RoomInfoResp: &m.roomInfo.RoomInfoResp}
		pbData.RoomId = m.roomId
		pbData.UrgencyVoteNum = m.roomInfo.GetUrgencyVoteNum(userID)
		pbData.AiTask = m.roomInfo.AiTask
		if player.IsRobot != 1 {
			player.Pid.Cast("enterRoomSuccess", pbData)
		}
	}
}

func (m *Room) ninjaAttack(attReq *include.AttackReq) {
	m.roomInfo.ChangePlayerGameStatus(attReq.SufferUserId, proto3.PlayerGameStatus_killed)
	m.roomInfo.AddKillNum(attReq.UserId, attReq.SufferUserId)
	// 判断游戏是否结束
	_, b := m.roomInfo.JudgeGameEndByRoleNum()
	if b {
		m.roomPid.Cast("gameEnd", nil)
		logger.Log.Infof("room:%d gameEnd, because ninjaAttack JudgeGameEndByRoleNum is true", m.roomId)
	} else {
		playerInfo := m.roomInfo.GetPlayerGameInfo(attReq.UserId)
		if playerInfo.IsRobot != int32(proto3.CommonStatusEnum_true) {
			// 处理使用过自爆符的情况
			m.DealBoom(attReq.UserId, attReq.SufferUserId)
			m.KillPeopleDropItem(attReq.UserId)
			m.AttackedClearSkills(attReq.SufferUserId)
		}
	}

	logger.Log.Infof("ninja attack roommate:%v, player:%d, sufferId:%d", m.roomInfo.AllPlayerGameInfo, attReq.UserId, attReq.SufferUserId)
	cmd := proto3.ProtoCmd_CMD_NinjaAttackResp
	pbData := &proto3.NinjaAttackResp{
		NinjaStatus: proto3.CommonStatusEnum_true,
		UserId:      attReq.UserId,
		SuffererId:  attReq.SufferUserId,
	}
	player := global.GloInstance.GetPlayer(attReq.SocketUser).(*Player)
	m.sendMsgPlayer(player, cmd, pbData)
	// player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	m.RecordKillNum(1)
}

func (m *Room) startVote(userID, sufferID int32, voteType proto3.VoteTypeEnum) bool {
	if !m.roomInfo.CheckPlayerGameStatus(userID, proto3.PlayerGameStatus_normal) {
		return false
	}
	b := m.roomInfo.CheckPlayerGameStatus(userID, proto3.PlayerGameStatus_normal) && m.roomInfo.CheckPlayerGameStatus(sufferID, proto3.PlayerGameStatus_killed) && !m.roomInfo.FoundPlayeDead(sufferID)
	if voteType == proto3.VoteTypeEnum_normal_vote && !b {
		return false
	}
	if voteType == proto3.VoteTypeEnum_urgent_vote {
		n := m.roomInfo.GetUrgencyVoteNum(userID)
		if n <= 0 {
			logger.Log.Infof("userID urgencyVoteNum < 0, urgencyVoteNum:%d", userID, n)
			return false
		}
		now := util.UnixTime()
		if now < m.nextUrgencyVoteTime {
			logger.Log.Infof("userID urgency frozen, urgencyVoteNum:%d", n)
			return false
		}
	}

	m.roomInfo.SetInitiateVote(userID, sufferID, voteType)

	m.roomInfo.VoteStep = proto3.VoteStepEnum_vote_chat
	cmd := proto3.ProtoCmd_CMD_StartVoteResp
	pbData := &proto3.StartVoteResp{
		Roommate:     m.roomInfo.GetProto3Roommate(),
		RoomInfoResp: &m.roomInfo.RoomInfoResp,
		UserId:       userID,
		SuffererId:   sufferID,
		VoteType:     voteType,
		TurnVoteNum:  m.roomInfo.VoteCount,
	}
	m.spreadPlayers(cmd, pbData)
	return true
}

func (m *Room) playerVote(voteMsg *include.VoteMsg) {
	m.roomInfo.VoteInfo.AddVoteInfo(voteMsg)

	cmd := proto3.ProtoCmd_CMD_PlayerVoteResp
	pbData := &proto3.PlayerVoteResp{
		UserId:  voteMsg.UserId,
		ErrCode: proto3.ErrEnum_Error_Pass,
		Status:  voteMsg.Status,
	}
	logger.Log.Infof("player vote msg:%v", voteMsg)

	m.spreadPlayers(cmd, pbData)

	// 判断是否全部进行了投票
	liveNum := m.roomInfo.GetLivePlayerNum()
	logger.Log.Infof("liveNum:%d get voteNum:%v", liveNum, m.roomInfo.VoteInfo.GetVoteNum())
	if liveNum <= m.roomInfo.VoteInfo.GetVoteNum() {
		m.roomPid.Cast("sumVote", 0)
	}
}

func (m *Room) voteEnd() {
	if m.roomInfo.RoomStatus == proto3.RoomStatus_voting {
		cmd := proto3.ProtoCmd_CMD_VoteResultResp
		// 处理使用过双倍票情况
		m.DealDoubleVote()

		killUserID, voteNum, voteStatus := m.roomInfo.VoteInfo.GetVotedResult(m.roomInfo.GetLivePlayerNum())
		logger.Log.Infof("userID:%d, voteNum:%d, votesStatus:%v", killUserID, voteNum, voteStatus)
		if m.DealGuard(killUserID) {
			killUserID = -1
			voteNum = -1
			voteStatus = proto3.VoteStatus_invalid_vote
		}
		m.DealTogether(killUserID)

		m.roomInfo.VoteInfo.VoteStatus = voteStatus
		m.roomInfo.VoteInfo.KillUserId = killUserID
		m.roomInfo.VoteInfo.VoteNumber = voteNum

		logger.Log.Infof("userID:%d, voteNum:%d, votesStatus:%v", killUserID, voteNum, voteStatus)
		pbData := &m.roomInfo.VoteInfo.VoteResultResp
		logger.Log.Infof("vote end send msg:%v", pbData)

		// 修改被票杀人的信息
		m.roomInfo.VoteKillPlayer(killUserID)
		m.spreadPlayers(cmd, pbData)

	}
}

func (m *Room) voteClose() {
	// 游戏开始
	m.roomInfo.RoomStatus = proto3.RoomStatus_gameing
	m.nextUrgencyVoteTime = util.UnixTime() + tableconfig.ConstsConfigs.GetIdValue(constant.UrgVoteFrozenTime)
	// 保持投票历史信息
	m.roomInfo.ArchiveVoteInfo()
	// 清理帧数据
	m.hisFrameList = []*proto3.FSPFrameResp{}
	m.frameList = []*proto3.FSPMsg{}
	m.tickerNum = 0

	// 设置默认的用户帧数据
	m.initCurrentFrame()

	// 判断游戏是否结束
	_, b := m.roomInfo.JudgeGameEndByRoleNum()
	if b {
		m.roomPid.Cast("gameEnd", nil)
		logger.Log.Infof("room:%d gameEnd, because voting JudgeGameEndByRoleNum is true", m.roomId)
	}
}

func (m *Room) PlayerOffline(userID int32) {
	userInfo := m.roomInfo.GetPlayerGameInfo(userID)
	userInfo.Player = nil
	// for i, v := range m.roomInfo.AllPlayerGameInfo {
	// 	if v.UserId == userID {
	// 		m.roomInfo.AllPlayerGameInfo[i].Player = nil
	// 	}
	// }
	m.roomInfo.LampUserMap[userID] = 0
	pbdata := new(proto3.LampSwitchResp)
	pbdata.UserId = userID
	pbdata.Status = 0
	m.spreadPlayers(proto3.ProtoCmd_CMD_LampSwitchResp, pbdata)
}

func (m *Room) playerExit(userID int32) {
	logger.Log.Infof("room player exit userId:%d, roomId:%v roomInfo:%v", userID, m.roomId, m.roomInfo.AllPlayerGameInfo)
	roomPlayerNum := 0
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		if v.UserId == userID {
			if v.PlayerGameStatus == proto3.PlayerGameStatus_normal {
				m.roomInfo.AllPlayerGameInfo[i].PlayerAutoOut = autoOutGaming
			} else {
				m.roomInfo.AllPlayerGameInfo[i].PlayerAutoOut = autoOutGameEnd
			}
			if m.roomInfo.RoomStatus == proto3.RoomStatus_game_end {
				m.roomInfo.AllPlayerGameInfo[i].PlayerAutoOut = autoOutGameEnd
			}
		}
		if m.roomInfo.AllPlayerGameInfo[i].PlayerAutoOut == autoOutNormal {
			roomPlayerNum++
		}
	}
	if roomPlayerNum <= 0 {
		m.roomPid.Cast("gameEnd", nil)
		logger.Log.Infof("room:%d gameEnd, because player num = 0", m.roomId)
	}
	if m.isOnlyGame {
		m.roomPid.Cast("gameEnd", nil)
	}
}
func (m *Room) gameEnd() {
	// 设置游戏时长
	m.roomInfo.GameDuration = util.UnixTime() - int32(m.roomInfo.GameWaitEndTime)
	cmd := proto3.ProtoCmd_CMD_GameEndResp
	spreadPbCmd := proto3.ProtoCmd_CMD_GameDetailResp
	spreadPbData := &proto3.GameDetailResp{}
	m.settleGameInfo()
	m.saveGameInfo()

	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		if v.PlayerAutoOut == autoOutGaming {
			// 主动退出不在进入房间
			continue
		}
		userID := v.UserId
		pbData := &proto3.GameEndResp{}
		pbData.GameWinRole = m.roomInfo.GameWinRole
		if v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			voteCount, hidVote := m.roomInfo.GameEndWolfSettle(userID, v.TurnIndex, pbData)
			pbData.VoteCount = voteCount
			pbData.HideNum = hidVote
			pbData.Rewards = v.SettleData.Rewards
			pbData.KillNum = v.KillNum
			pbData.BreakNum = m.GetTriggerUrgencyNum(userID)
			pbData.Score += m.getWolfScore(v.UserId, v.SettleData)
			// 狼人评分-逃票系数
			pbData.Score += hidVote * tableconfig.ConstsConfigs.GetIdValue(constant.HidScoreRatio)

			// m.sendMsgPlayer(player, cmd, pbData)
			logger.Log.Infof("game end send data:%v", pbData)
		} else {
			voteCount, successVote := m.roomInfo.GameEndPeopleSettle(userID, v.TurnIndex, pbData)
			pbData.VoteCount = voteCount
			pbData.VoteNum = successVote
			pbData.Rewards = v.SettleData.Rewards
			pbData.AllTaskNum = tableconfig.ConstsConfigs.GetIdValue(constant.TaskTotalScoreId)
			pbData.AllTaskedNum = int32(m.totalScore)
			userTask := m.taskMap[v.UserId]
			pbData.PrivateTaskNum = userTask.GetTaskNum()
			pbData.PrivateTaskedNum = GetPlayerTaskProgressNum(m.taskMap[userID])
			// 好人评分--任务系数
			pbData.Score += pbData.PrivateTaskedNum * tableconfig.ConstsConfigs.GetIdValue(constant.TaskScoreRatio)
			// 好人评分--任务占比系数
			if m.totalScore != 0 {
				pbData.Score += int32(float32(pbData.PrivateTaskedNum) / float32(m.totalScore) * float32(tableconfig.ConstsConfigs.GetIdValue(constant.TaskInScoreRatio)))
			}
			// 好人评分--投凶成功系数
			pbData.Score += successVote * tableconfig.ConstsConfigs.GetIdValue(constant.VoteScoreRatio)
			// 好人评分--发现尸体加分
			pbData.Score += v.CallVoteNum * tableconfig.ConstsConfigs.GetIdValue(constant.CallVoteScoreRatio)
			// player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
			// m.sendMsgPlayer(player, cmd, pbData)

			logger.Log.Infof("game end send data:%v", pbData)
		}
		pbData.UserId = v.UserId
		pbData.ChatNum = v.ChatNum
		pbData.PlayerGameStatus = v.PlayerGameStatus

		// player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
		p := global.GloInstance.GetPlayer(userID)
		if v.PlayerAutoOut != autoOutGaming && p != nil {
			player := p.(*Player)
			if player.IsRobot != 1 {
				m.sendMsgPlayer(player, cmd, pbData)
				go player.GetFreshGiftStep()
			}
		}

		spreadPbData.PlayerDetail = append(spreadPbData.PlayerDetail, pbData)

		if v.PlayerRole == m.roomInfo.GameWinRole {
			RecordAction(userID, constant.Victory)
		} else {
			RecordAction(userID, constant.Failure)
		}
	}
	m.spreadPlayers(spreadPbCmd, spreadPbData)
	// m.roomPid.CastStop()
	// 记录游戏局数
	m.RecordGameData()
	// 记录阵营胜率
	m.RecordGroupWinData()
	// 记录会议次数
	m.RecordConferNum()
}

func (m *Room) DoubleRewards(userID int32) {
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		if v.UserId == userID {
			p := global.GloInstance.GetPlayer(userID)
			if p != nil {
				player := p.(*Player)
				player.AddRewards(v.SettleData.Rewards)
				break
			}
		}
	}
}

func (m *Room) getWolfScore(userID int32, settle *include.SettletData) int32 {
	// 狼人评分-破坏系数
	score := m.GetTriggerUrgencyNum(userID) * tableconfig.ConstsConfigs.GetIdValue(constant.BreakScoreRatio)
	// 狼人评分-击破系数
	score += settle.KillTotal * tableconfig.ConstsConfigs.GetIdValue(constant.KillScoreRatio)

	return score
}

func (m *Room) roomTicker() {
	// 游戏结束停止发送帧
	ticketKey := constant.TimerRoomRickerPrefix + strconv.Itoa(int(m.roomId))
	if m.roomInfo == nil || m.roomInfo.RoomStatus == proto3.RoomStatus_game_end || m.roomInfo.RoomStatus == proto3.RoomStatus_game_settle {
		m.roomPid.StopPidTimer(ticketKey)
		return
	}
	m.tickerNum++
	// 发送默认的帧数据 没有收到客户端帧数据，服务器也进行定时发送
	t := tableconfig.ConstsConfigs.GetFPSNumParam()
	m.roomPid.SendAfter("room_ticker", ticketKey, t, nil) // 自己掉自己进行帧同步

	pbData := &proto3.FSPFrameResp{}
	// tmp := make([]proto3.FSPMsg, len(m.frameList))
	clientFrameMsg := m.frameList[:]
	if len(clientFrameMsg) > 0 {
		// 默认发送最新帧数据
		for _, v := range clientFrameMsg {
			if v.Cmd == proto3.FSPCmd_CMD_IDLE || v.Cmd == proto3.FSPCmd_CMD_WALK {
				logger.Log.Debug(v)
				m.UpdateCurrentFrame(v)
			} else {
				pbData.Msgs = append(pbData.Msgs, v)
			}

			// if clientFrameMsg[i].Cmd == proto3.FSPCmd_CMD_ATTACK {
			// 	logger.Log.Info("attack:", clientFrameMsg)
			// }
			// m.UpdateCurrentFrame(v)
		}
	}
	pbData.FrameId = uint32(m.tickerNum)
	pbData.Msgs = append(pbData.Msgs, m.currentFrame...)
	msgs := make([]*proto3.FSPMsg, 0)
	for i := range pbData.Msgs {
		tmp := *pbData.Msgs[i]
		msgs = append(msgs, &tmp)
	}
	pbData.Msgs = msgs //tmpPb.Msgs
	cmd := proto3.ProtoCmd_CMD_FrameSyncResp
	m.spreadPlayers(cmd, pbData)

	m.frameList = []*proto3.FSPMsg{}

	m.hisFrameList = append(m.hisFrameList, pbData)
	logger.Log.Debugf("roomID:%v send fms:%v currentFms:%v", m.roomId, pbData, m.currentFrame)
}

func addRewards(t, num string, ret *string) {
	items := strings.Split(*ret, "|")
	for i, v := range items {
		item := strings.Split(v, ",")
		if item[0] == t {
			value := util.ToStr(util.ToInt(item[1]) + util.ToInt(num))
			items[i] = item[0] + "," + value
			break
		}
		if i == len(items)-1 {
			tmp := t + "," + num
			items = append(items, tmp)
		}
	}
	*ret = strings.Join(items, "|")
}

func (m *Room) settleGameInfo() {
	if m.isOnlyGame {
		// return
	}
	logger.Log.Infof("roomID::%v game end start save game info", m.roomId)
	for i := range m.roomInfo.AllPlayerGameInfo {
		// v := m.roomInfo.AllPlayerGameInfo[i]
		settleData := &include.SettletData{}
		m.roomInfo.AllPlayerGameInfo[i].SettleData = settleData
		userInfo := m.roomInfo.AllPlayerGameInfo[i]
		p := global.GloInstance.GetPlayer(userInfo.UserId)
		if p == nil {
			// TODO
			continue
		}
		player := p.(*Player)
		settleData.GameDuration = m.roomInfo.GameDuration // 游戏时长
		settleData.MatchGameNum = 1                       // 游戏盘数
		// 判断掉线
		if userInfo.PlayerAutoOut == autoOutGaming {
			// 主动退出算掉线
			settleData.OffLine = 1
			continue
		}
		if userInfo.PlayerRole == proto3.PlayerRoleEnum_comm_people {
			if v, ok := m.taskMap[userInfo.UserId]; ok {
				if v.GetFinishTaskNum() < 1 {
					// 主动退出算掉线
					settleData.OffLine = 1
				}
			}
		} else {
			// 掉线逻辑
			// if userInfo.KillNum < 1 {
			// 	// 主动退出算掉线
			// 	settleData.OffLine = 1
			// }
		}

		if settleData.OffLine == 1 {
			continue
		}

		if userInfo.PlayerRole == m.roomInfo.GameWinRole { // 胜利场数与段位
			settleData.MatchWinNum = 1
			settleData.StarCount = 1
			settleData.Star = 1
			if userInfo.PlayerRole == m.roomInfo.GameWinRole {
				settleData.WolfWinNum = 1
			} else {
				settleData.PoorWinNum = 1
			}
		} else {
			if m.roomInfo.GameWinRole != proto3.PlayerRoleEnum_invalid_role {
				settleData.StarCount = -1
				settleData.Star = -1
			}
		}

		// 获取金币和经验
		// gold, exp := tableConfig.QuaLevelConf.GetGoldAndExp(player.Attr.RankID)
		// settleData.Exp += exp
		// settleData.Gold += gold
		quaLevelConfig := tableconfig.QuaLevelConfs.GetQuaConfig(player.Attr.RankID)
		if quaLevelConfig == nil {
			logger.Log.Errorf("qualifyinglevels.excel don't set please set rankID:%d", player.Attr.RankID)
		}
		if quaLevelConfig != nil {
			settleData.Rewards = quaLevelConfig.Rewards
		}

		if userInfo.PlayerRole == m.roomInfo.GameWinRole {
			settleData.Exp += quaLevelConfig.WinExpReward   // tableConfig.ConstsConfigs.GetIdValue(include.WinAnnxExpReward)
			settleData.Gold += quaLevelConfig.WinGoldReward // tableConfig.ConstsConfigs.GetIdValue(include.WinAnnxGlodReward)
		}

		// 狼人杀人数
		settleData.KillTotal = userInfo.KillNum
		if userInfo.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			settleData.WolfKillTotal += userInfo.KillNum
			settleData.MatchWolfNum++
		} else {
			settleData.PoorKillTotal += userInfo.KillNum
			settleData.Gold += int32(m.totalScore) * quaLevelConfig.AllFinReward //tableConfig.ConstsConfigs.GetIdValue(include.TaskAllReward)
			if task, ok := m.taskMap[player.Attr.UserID]; ok {
				settleData.Gold += GetPlayerTaskProgressNum(task) * quaLevelConfig.TaskFinReward
			}
		}
		settleData.Gold += userInfo.KillNum * quaLevelConfig.KillReward // tableConfig.ConstsConfigs.GetIdValue(include.KillReward)

		// 被票杀次数
		if userInfo.PlayerGameStatus == proto3.PlayerGameStatus_killed_vote {
			settleData.BevoteedTotal = 1
		} else if userInfo.PlayerGameStatus == proto3.PlayerGameStatus_killed {
			settleData.BekilledTotal = 1
		}
		if userInfo.PlayerGameStatus == proto3.PlayerGameStatus_normal {
			userInfo.TurnIndex = m.roomInfo.TurnIndex + 1 // 未死亡
		}
		userInfo.VoteTotal, userInfo.VoteCorrectTotal, userInfo.VoteFailedTotal = m.roomInfo.GetVoteNum(userInfo.UserId, userInfo.TurnIndex)
		settleData.VoteTotal = userInfo.VoteTotal
		settleData.VoteCorrectTotal = userInfo.VoteCorrectTotal
		settleData.VoteFailedTotal = userInfo.VoteFailedTotal

		logger.Log.Infof("UpdateGameData settleData:%v usreID:%v roomId:%d", *settleData, userInfo.UserId, m.roomId)

		// 设置称号数据
		userTitle := player.Attr.UserTitle
		if userInfo.UserId == m.roomInfo.FirstKillUserID {
			settleData.KeepFirstOut = 1
		} else {
			settleData.KeepFirstOut = 0
		}
		if userInfo.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			settleData.KeepWolf = 1
			settleData.KeepPoor = 0
		} else if userInfo.PlayerRole == proto3.PlayerRoleEnum_comm_people {
			settleData.KeepPoor = 1
			settleData.KeepWolf = 0
		}

		// TODO 未得到道具
		settleData.KeepNoItem = m.GetPlayerSkillNum(userInfo.UserId)

		// TotalWolfDay
		nowDate := util.UnixTime()
		if !util.IsSameDay(nowDate, userTitle.WolfTimestamp) && userInfo.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			userTitle.WolfTimestamp = nowDate
			settleData.TotalWolfDay = 1
		} else {
			settleData.TotalWolfDay = 0
		}

		// 获取玩家完成任务
		if task, ok := m.taskMap[userInfo.UserId]; ok {
			l := strings.Split(task.PointInfo, ",")
			settleData.TotalTask = int32(len(l))
		}

		addRewards(util.ToStr(constant.ItemIdExp), util.ToStr(settleData.Exp), &settleData.Rewards)
		addRewards(util.ToStr(constant.ItemIdGold), util.ToStr(settleData.Gold), &settleData.Rewards)
	}
}

func (m *Room) saveGameInfo() {
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		p := global.GloInstance.GetPlayer(v.UserId)
		if p == nil {
			// TODO
			continue
		}
		player := p.(*Player)
		levelData := &player.Attr.UserData
		userTitle := player.Attr.UserTitle
		settleData := v.SettleData
		nowDate := util.UnixTime()
		if !util.IsSameDay(nowDate, userTitle.WolfTimestamp) {
			userTitle.WolfTimestamp = nowDate
		}

		// 段位数据
		levelData.GameDuration += settleData.GameDuration
		levelData.MatchGameNum += settleData.MatchGameNum
		levelData.MatchWolfNum += settleData.MatchWolfNum
		levelData.OfflineNum += settleData.OffLine
		if settleData.OffLine == 1 {
			continue
		}
		levelData.MatchWinNum += settleData.MatchWinNum
		levelData.WolfWinNum += settleData.WolfWinNum
		levelData.PoorWinNum += settleData.PoorWinNum

		protectRank := tableconfig.ConstsConfigs.GetIdValue(constant.ProtectRankID)
		if settleData.StarCount > 0 || levelData.RankID >= protectRank {
			levelData.StarCount += settleData.StarCount
			if levelData.StarCount < 0 {
				levelData.StarCount = 0
			}
			levelData.Star += settleData.Star
			levelData.RankID, levelData.Star = tableconfig.QuaLevelConfs.GetChangeRank(player.Attr.RankID, levelData.Star)
		}

		levelData.KillTotal += settleData.KillTotal
		levelData.WolfKillTotal += settleData.WolfKillTotal
		levelData.PoorKillTotal += settleData.PoorKillTotal
		levelData.BekilledTotal += settleData.BekilledTotal
		levelData.BevoteedTotal += settleData.BevoteedTotal
		levelData.VoteTotal += settleData.VoteTotal
		levelData.VoteFailedTotal += settleData.VoteFailedTotal

		// 称号数据
		userTitle.TotalWolfDay += settleData.TotalWolfDay
		if settleData.KeepFirstOut == 0 {
			userTitle.KeepFirstOut = 0
		} else {
			userTitle.KeepFirstOut += settleData.KeepFirstOut
		}
		if settleData.KeepWolf == 0 {
			userTitle.KeepWolf = 0
		} else {
			userTitle.KeepWolf += settleData.KeepWolf
		}
		if settleData.KeepPoor == 0 {
			userTitle.KeepPoor = 0
		} else {
			userTitle.KeepPoor += settleData.KeepPoor
		}
		// 连续未得到道具
		if settleData.KeepNoItem != 0 {
			userTitle.KeepNoItem = 0
		} else {
			userTitle.KeepNoItem++
		}
		userTitle.TotalTask += settleData.TotalTask

		player.AddRewards(settleData.Rewards)

		// 设置类型值
		player.initAchievementMap()
		// 判断是否达到称号
		player.AlterUserTitle()

		cmd := proto3.ProtoCmd_CMD_NewRankResp
		pbData := &proto3.NewRankResp{RankId: levelData.RankID, Star: levelData.Star, Level: levelData.Level, Exp: levelData.Exp, MaxExp: levelData.MaxExp}
		m.sendMsgPlayer(player, cmd, pbData)

		// 更新排名
		topBoard.UpdateBoard(player)
	}
}

func (m *Room) SendChat(userID int32, chatData string, quickType, quickId int32, subsText []string) {
	if len(chatData) <= 0 && quickId < 0 {
		return
	}
	cmd := proto3.ProtoCmd_CMD_GameChatResp
	pbData := &proto3.GameChatResp{
		UserId:    userID,
		ChatData:  chatData,
		QuickType: quickType,
		QuickId:   quickId,
		SubsText:  subsText,
	}

	m.spreadPlayers(cmd, pbData)
}

func (m *Room) SendWaitChat(userID int32, chatData string, quickType, quickId int32) {
	if len(chatData) <= 0 && quickId < 0 {
		return
	}
	cmd := proto3.ProtoCmd_CMD_WaitGameChatResp
	pbData := &proto3.WaitGameChatResp{
		UserId:    userID,
		ChatData:  chatData,
		QuickType: quickType,
		QuickId:   quickId,
	}

	m.spreadPlayers(cmd, pbData)
}

func (m *Room) AddChatNum(userID int32, chatData string) {
	for i := range m.roomInfo.AllPlayerGameInfo {
		v := m.roomInfo.AllPlayerGameInfo[i]
		if userID == v.UserId {
			v.ChatNum += util.CountStrNum(chatData)
		}
	}
}

func (m *Room) CheckPlayerGameStatus(userID int32, gameStatus proto3.PlayerGameStatus) bool {
	return m.roomInfo.CheckPlayerGameStatus(userID, gameStatus)
}

func (m *Room) ChangeSkin(userID, skinID int32) {
	userInfo := m.roomInfo.GetPlayerGameInfo(userID)
	if userInfo.PlayerAutoOut != autoOutGaming {
		logger.Log.Infof("user：%d is in room, the room info:%v", userID, m.roomInfo.AllPlayerGameInfo)
		userInfo.UserSkin = skinID
	}
	// for i := range m.roomInfo.AllPlayerGameInfo {
	// 	v := m.roomInfo.AllPlayerGameInfo[i]
	// 	if v.UserId == userID && v.PlayerAutoOut != autoOutGaming {
	// 		v.UserSkin = skinID
	// 		logger.Log.Infof("user：%d is in room, the room info:%v", userID, m.roomInfo.AllPlayerGameInfo)
	// 		return
	// 	}
	// }
	return
}

func (m *Room) spreadPlayers(cmd proto3.ProtoCmd, pbData interface{}) {
	for i := 0; i < len(m.roomInfo.AllPlayerGameInfo); i++ {
		v := m.roomInfo.AllPlayerGameInfo[i]
		// per := global.GloInstance.GetPlayer(p.UserId)
		if v.PlayerAutoOut == autoOutNormal && v.Player != nil && v.Player.Attr.RoomID == m.roomId {
			// player := per.(*Player)
			logger.Log.Debugf("cmd:%d pbData:%v", cmd, pbData)
			v.Player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
		}
	}
}

func (m *Room) sendMsgPlayer(player *Player, cmd proto3.ProtoCmd, pbData interface{}) {
	if player != nil && player.Room != nil {
		if player.Room.roomId == m.roomId {
			logger.Log.Debugf("cmd:%d pbData:%v", cmd, pbData)
			player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
		}
	}
}

// 关门
func (m *Room) CloseDoor(userId int32, doorId int32) {
	if !m.roomInfo.IsWolfMan(userId) { // 非狼人不能关门
		return
	}
	cd := util.ToInt64(tableconfig.ConstsConfigs.GetValueById(constant.CloseDoorCD))
	player := global.GloInstance.GetPlayer(userId)
	if player == nil {
		return
	}
	p := player.(*Player)
	cdEndTime := m.closeDoorTime + cd
	curTime := time.Now().Unix()
	if cdEndTime > curTime { // 冷却时间未结束
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	m.closeDoorTime = curTime
	cmd := proto3.ProtoCmd_CMD_CloseDoorResp
	pbData := &proto3.CloseDoorResp{DoorId: doorId, CdEndTime: curTime + cd}
	// p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	m.sendMsgPlayer(p, cmd, pbData)
}

func (r *Room) GetRoomID() int32 {
	return r.roomId
}

// 竞拍
func (r *Room) DealAuction(userId int32) {
	player := global.GloInstance.GetPlayer(userId)
	if player == nil {
		logger.Log.Errorln("DealAuction player is nil")
		return
	}
	p := player.(*Player)
	auction := r.auction
	if auction == nil {
		auction = &include.Auction{}
		r.auction = auction
	}
	auctionPrice := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.OneAuctionPrice))
	auction.CurPrice += auctionPrice
	if p.Attr.Gold < auction.CurPrice {
		p.ErrorResponse(proto3.ErrEnum_Error_Gold_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Gold_Not_Enough)])
		return
	}
	if auction.AuctionTimeMap == nil {
		auction.AuctionTimeMap = make(map[int32]int32)
	}
	curTime := util.UnixTime()
	auctionCD := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.AuctionCD))
	if curTime < auction.AuctionTimeMap[userId]+auctionCD { // CD未结束
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	auction.AuctionTimeMap[userId] = curTime

	// 竞拍次数
	clickNumMap := auction.ClickNumMap
	if clickNumMap == nil {
		clickNumMap = make(map[int32]int32)
		auction.ClickNumMap = clickNumMap
	}
	auction.ClickNumMap[userId]++

	// 总数加1
	auction.ClickSum++
	auction.UserId = userId

	//p.Attr.Gold -= auctionPrice

	pbData := &proto3.AuctionResp{ErrNum: proto3.ErrEnum_Error_Pass, UserId: userId, UserName: p.Attr.Username, CurPrice: r.auction.CurPrice, Gold: p.Attr.Gold}
	r.spreadPlayers(proto3.ProtoCmd_CMD_AuctionResp, pbData)
}

// 竞拍结束
func (r *Room) AuctionEnd() {
	auction := r.auction
	if auction == nil {
		return
	}
	userId := auction.UserId
	player := global.GloInstance.GetPlayer(userId)
	if player == nil {
		logger.Log.Errorln("AuctionEnd player is nil")
		return
	}
	p := player.(*Player)
	var dropGroupId int32
	if r.roomInfo.IsWolfMan(userId) { // 狼人
		dropGroupId = util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.AuctionDropIdWolfMan))
	} else { // 平民
		dropGroupId = util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.AuctionDropIdNormalMan))
	}
	dropGroups := tableconfig.DropGroupConfigs.GetDropGroupMap(int(dropGroupId))
	if len(dropGroups) == 0 {
		logger.Log.Errorln("AuctionEnd is null")
		return
	}
	taskMap := r.taskMap
	if taskMap == nil {
		logger.Log.Errorln("AuctionEnd taskMap is nil")
		return
	}
	task := taskMap[userId]
	if task == nil { // 狼人可能为nil
		task = &Task{UserId: userId}
	}
	dropGroupConf := GetDropGroupConfWithBuf(task, dropGroups)
	if dropGroupConf == nil {
		logger.Log.Info("AuctionEnd dropGroupConfig is nil")
		return
	}
	r.addDropItem(p.Attr.UserID, dropGroupConf)
	// 扣对应金币
	p.Attr.Gold -= auction.CurPrice

	// 下发客户端
	dropItemMap := make(map[int32][]*proto3.Item)
	items := []*proto3.Item{{Id: int32(dropGroupConf.ItemId), Num: int32(dropGroupConf.Num)}}
	dropItemMap[int32(dropGroupConf.ItemType)] = items
	p.Pid.Cast("dropItemResp", dropItemMap)

	// 返还金币: 返还数量 =（结算金额/2）*（玩家点击次数/总点击次数）
	returnGoldRatio := tableconfig.ConstsConfigs.GetValueById(constant.ReturnGoldRatio)
	clickNumMap := auction.ClickNumMap
	clickSum := auction.ClickSum
	if clickNumMap != nil && clickSum != 0 {
		for k, v := range clickNumMap {
			if k == userId {
				continue
			}
			returnGold := float32(auction.CurPrice/util.ToInt(returnGoldRatio)) * (float32(v) / float32(clickSum))
			temPlayer := global.GloInstance.GetPlayer(k)
			if temPlayer == nil {
				logger.Log.Errorf("AuctionEnd player %v is nil", k)
				continue
			}
			tempP := temPlayer.(*Player)
			tempP.AddItems(util.ToStr(constant.ItemIdGold) + "," + util.ToStr(int32(returnGold)))
		}
	}
}

func (r *Room) CheckAttackFrozen(userID int32) bool {
	userInfo := r.roomInfo.GetPlayerGameInfo(userID)
	if userInfo == nil {
		return false
	}
	nowTime := util.UnixTime()
	if userInfo.nextAttackTime != 0 && userInfo.nextAttackTime > nowTime {
		return true
	}
	return false
}

// 获取AI平民数
func (r *Room) getAiComPeopleNum() int32 {
	if r.roomInfo == nil {
		return 0
	}
	players := r.roomInfo.AllPlayerGameInfo
	if players == nil {
		return 0
	}
	var count int32
	for _, v := range players {
		if v.IsRobot != 1 { // 非机器人
			continue
		}
		if v.PlayerRole == proto3.PlayerRoleEnum_comm_people {
			count++
		}
	}
	return count
}

// 定时给机器人加积分
func (r *Room) scheduleAddScore() {
	ticketKey := constant.TimerAddScorePrefix + strconv.Itoa(int(r.roomId))
	if r == nil {
		r.roomPid.StopPidTimer(ticketKey)
		return
	}
	interTime := tableconfig.ConstsConfigs.GetValueById(constant.AIAddScoreInterTime)
	interTimes := strings.Split(interTime, ",")
	randTime := util.RandInt(util.ToInt(interTimes[0]), util.ToInt(interTimes[1]))
	// 投票环节 不加积分
	if r.CheckRoomGameStatus(proto3.RoomStatus_voting) {
		r.roomPid.SendAfter("scheduleAddScore", ticketKey, randTime*1000, nil) // 自己掉自己
		return
	}
	// 如果达到了（机器人数*6）任务数 则停止加积分
	if r.aiScore >= r.getAiComPeopleNum()*constant.TaskPoints {
		r.roomPid.StopPidTimer(ticketKey)
		return
	}
	r.totalScore++
	r.aiScore++

	// 机器人增加任务
	var robotId, n int32 = -1, 0
	for {
		n++
		if n > 100 {
			break
		}
		// rankID := util.RandInt(0, int32(len(r.roomInfo.RobotUserList))-1)
		// robotId = r.roomInfo.RobotUserList[rankID]

		if userInfo := r.roomInfo.GetALiveAiRankUser(); userInfo != nil && userInfo.IsRobot == 1 && userInfo.PlayerGameStatus == proto3.PlayerGameStatus_normal && userInfo.PlayerRole == proto3.PlayerRoleEnum_comm_people {
			task := r.taskMap[robotId]
			if task != nil {
				point := strings.Split(task.PointInfo, ",")
				if len(point) > constant.TaskPoints {
					continue
				}
				if len(task.PointInfo) <= 0 {
					task.PointInfo = "1"
				} else {
					task.PointInfo += ",1" // 临时增加机器人完成任务，便于计算
				}
				r.taskMap[robotId] = task
				break
			}
		}
	}

	// 判断是否达到最大任务分数
	maxTotalScoreStr := tableconfig.ConstsConfigs.GetValueById(constant.TaskTotalScoreId)
	maxTotalScore, _ := strconv.Atoi(maxTotalScoreStr)
	if r.totalScore >= maxTotalScore {
		r.roomInfo.GameWinRole = proto3.PlayerRoleEnum_comm_people
		// 游戏结束逻辑调用
		r.roomPid.Cast("gameEnd", nil)
		logger.Log.Infof("room:%d gameEnd, because totalScore >= maxTotalScore:%v", r.roomId, maxTotalScore)
		return
	}
	// 广播分数
	cmd := proto3.ProtoCmd_CMD_TotalScoreResp
	pbData := &proto3.TotalScoreResp{}
	pbData.TotalScore = int32(r.totalScore)
	r.spreadPlayers(cmd, pbData)

	r.roomPid.SendAfter("scheduleAddScore", ticketKey, randTime*1000, nil) // 自己掉自己
}

// 通风口
func (r *Room) OpenWind(userId int32, windId int32) {
	if windId < 0 || windId > 99 {
		return
	}
	player := global.GloInstance.GetPlayer(userId)
	if player == nil {
		return
	}
	p := player.(*Player)
	if !r.roomInfo.IsWolfMan(userId) { // 非狼人不能通风
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	if !p.CheckPlayerGameStatus(proto3.PlayerGameStatus_normal) {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	cmd := proto3.ProtoCmd_CMD_OpenWindResp
	pbData := &proto3.OpenWindResp{WindId: windId}
	//p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	r.sendMsgPlayer(p, cmd, pbData)
}
