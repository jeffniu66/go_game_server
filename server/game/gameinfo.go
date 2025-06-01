package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/sdk"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"time"
)

type Roommate = proto3.Roommate

type PlayerGameInfo struct {
	Roommate
	PlayerAutoOut    int32                // 0-正常游戏 1-死亡主动退出 2-游戏中主动退出
	IsFoundDead      bool                 // 死亡是否被发现，发现后不能重新发现
	CallVoteNum      int32                // 发起投票次数
	VoteTotal        int32                // 参与投票次数
	ChatNum          int32                // 发言字数
	VoteCorrectTotal int32                // 投票命中次数 命中狼人
	VoteFailedTotal  int32                // 投票失败次数 命中平民
	SettleData       *include.SettletData // 结算数据
	UrgencyVoteNum   int32                // 紧急投票次数
	nextAttackTime   int32                // 下一次攻击时间
	Player           *Player              `json:"-"` // player和room循环 so json -
}

type roomGamingInfo struct {
	GameWaitTime        int32                 // 候场时长
	GameVoteTime        int32                 // 投票时长
	GameDuration        int32                 // 游戏时长
	FirstKillUserID     int32                 // 第一个被杀害
	proto3.RoomInfoResp                       // 房间信息信息
	VoteInfo            voteInfo              // 投票信息-游戏处在投票过程中使用
	VoteInfoHis         []*voteInfo           // 投票历史
	VoteCount           int32                 // 投票次数
	AllPlayerGameInfo   []*PlayerGameInfo     // 玩家角色
	GameWinRole         proto3.PlayerRoleEnum // 胜利方
	AiTask              []*proto3.TaskInfo    // AI任务
	RobotUserList       []*PlayerGameInfo     // 机器人列表
	LiveUserList        []*PlayerGameInfo     // 活着的玩家列表 key userID,value userID
	LampUserMap         map[int32]int32       // 监控室开关
}

func (r *roomGamingInfo) GetPlayerGameInfo(userID int32) *PlayerGameInfo {
	for _, v := range r.AllPlayerGameInfo {
		if v.UserId == userID {
			return v
		}
	}
	return nil
}

func (r *roomGamingInfo) GetALiveWolfUserID() int32 {
	var robotId, userID int32
	for _, v := range r.AllPlayerGameInfo {
		if v.PlayerRole == proto3.PlayerRoleEnum_wolf_man && v.PlayerGameStatus == proto3.PlayerGameStatus_normal {
			if v.IsRobot == 1 {
				robotId = v.UserId
			} else {
				userID = v.UserId
			}
		}
	}
	if robotId <= 0 {
		return userID
	}
	return robotId
}

func (r *roomGamingInfo) GetALiveRankUser() *PlayerGameInfo {
	i := util.RandInt(0, int32(len(r.LiveUserList))-1)
	return r.LiveUserList[i]
}

func (r *roomGamingInfo) GetALiveAiRankUserNoMe(userID int32) *PlayerGameInfo {
	for i := 0; i < 10; i++ {
		voteUser := r.GetALiveRankUser()
		if voteUser.UserId != userID {
			return voteUser
		}
	}
	return r.GetALiveRankUser()
}

func (r *roomGamingInfo) GetALiveAiRankUser() *PlayerGameInfo {
	index := util.RandInt(0, int32(len(r.RobotUserList))-1)
	n := 0
	for i := index; i < int32(len(r.RobotUserList)); i++ {
		if r.RobotUserList[i].PlayerGameStatus == proto3.PlayerGameStatus_normal {
			index = int32(i)
			break
		} else {
			if i == int32(len(r.RobotUserList))-1 {
				i = 0
			}
		}
		n++
		if n >= len(r.RobotUserList) {
			break
		}
	}
	return r.RobotUserList[index]
}

func (r *roomGamingInfo) GetALiveRankPeopleUser() *PlayerGameInfo {
	var ret *PlayerGameInfo
	i := util.RandInt(0, int32(len(r.LiveUserList))-1)
	n := 0
	for {
		n++
		if r.LiveUserList[i].PlayerRole == proto3.PlayerRoleEnum_comm_people {
			ret = r.LiveUserList[i]
		} else {
			if i < int32(len(r.LiveUserList))-1 {
				i++
			} else {
				i = 0
			}
		}
		if n >= len(r.LiveUserList) {
			break
		}
	}
	return ret
}

func (r *roomGamingInfo) GetRoommates() []*proto3.Roommate {
	ret := make([]*proto3.Roommate, 0)
	for i := range r.AllPlayerGameInfo {
		v := r.AllPlayerGameInfo[i]
		ret = append(ret, &v.Roommate)
	}
	return ret
}
func (r *roomGamingInfo) IsWolfMan(userID int32) bool {
	for _, v := range r.AllPlayerGameInfo {
		if v.UserId == userID && v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			return true
		}
	}
	return false
}
func (r *roomGamingInfo) GetProto3Roommate() []*proto3.Roommate {
	ret := make([]*proto3.Roommate, 0)
	for i := range r.AllPlayerGameInfo {
		ret = append(ret, &r.AllPlayerGameInfo[i].Roommate)
	}
	return ret
}

// isOnly -no use, for only model
func (r *roomGamingInfo) InitRoomGamingInfo(playerIDList, birthList []int32, langArr []int, isOnly bool) {
	if switchRole := global.MyConfig.ReadInt32("switch", "switch_role_wolf"); switchRole != 1 {
		// 不允许机器人狼人
		var robotIndex, restIndex []int32
		for index, userID := range playerIDList {
			p := global.GloInstance.GetPlayer(userID)
			if p != nil {
				player := p.(*Player)
				if player.IsRobot == 1 {
					robotIndex = append(robotIndex, int32(index))
					continue
				}
			} else {
				robotIndex = append(robotIndex, int32(index))
				continue
			}
			isWolf := false
			for _, v := range langArr {
				if index == v {
					isWolf = true
				}
			}
			if !isWolf {
				restIndex = append(restIndex, int32(index))
			}
		}
		logger.Log.Infof("robotIndex:%v, restIndex:%v langArr:%v", robotIndex, restIndex, langArr)
		for _, v := range robotIndex {
			if len(restIndex) <= 0 {
				break
			}
			for i, vv := range langArr {
				if int(v) == vv {
					langArr[i] = int(restIndex[0])
					if len(restIndex) > 0 {
						restIndex = restIndex[1:]
					}
					break
				}
			}
		}
	}
	r.LampUserMap = make(map[int32]int32, 0)
	for i, v := range playerIDList {
		tmp := PlayerGameInfo{}
		p := global.GloInstance.GetPlayer(v)
		tmp.PlayerRole = r.playerRole(langArr, i)
		tmp.PlayerGameStatus = proto3.PlayerGameStatus_normal
		tmp.UserId = v
		randNum := util.RandInt(0, int32(len(birthList)-1))
		tmp.BirthPoint = birthList[randNum]
		birthList = append(birthList[:randNum], birthList[randNum+1:]...)
		// 玩家进入匹配，但已退出
		if p != nil {
			player := p.(*Player)
			tmp.Username = player.Attr.Username
			tmp.UserSkin = player.Attr.UseSkin
			tmp.UserPhoto = player.Attr.UserPhoto
			tmp.UserTitle = player.Attr.UserTitle.UseTitle
			tmp.OpenId = sdk.TokenMap[tmp.UserId]
			tmp.Player = player
			tmp.IsRobot = player.IsRobot
		} else {
			tmp.Username = db.GetUsername(tmp.UserId)
			if tmp.Username == "" {
				username, _, _ := GetRankName(tableconfig.NameZhConfigs)
				tmp.Username = username
			}
			tmp.IsRobot = 1
		}
		if tmp.IsRobot == 1 {
			r.RobotUserList = append(r.RobotUserList, &tmp)
		}
		tmp.UrgencyVoteNum = tableconfig.ConstsConfigs.GetIdValue(constant.UrgencyVoteNum)
		r.AllPlayerGameInfo = append(r.AllPlayerGameInfo, &tmp)
		r.LiveUserList = append(r.LiveUserList, &tmp)
		r.LampUserMap[v] = 0
	}
	// 初始化房间信息
	r.GameWaitTime = tableconfig.ConstsConfigs.GetGameWaitTime()
	r.GameVoteTime = tableconfig.ConstsConfigs.GetGameVoteTime()
	r.TurnIndex = 1
	r.RoomStatus = proto3.RoomStatus_wait_game
	if isOnly {
		r.IsOnly = proto3.CommonStatusEnum_true
	}
	r.VoteCount = 0
}

// func (r *roomGamingInfo) SetBirthPoint(userID, birthPoint int32) {
// 	for i := range r.AllPlayerGameInfo {
// 		if r.AllPlayerGameInfo[i].UserId == userID {
// 			r.AllPlayerGameInfo[i].BirthPoint = birthPoint
// 		}
// 	}
// }

func (r *roomGamingInfo) FoundPlayeDead(userID int32) bool {
	for i := 0; i < len(r.AllPlayerGameInfo); i++ {
		if r.AllPlayerGameInfo[i].UserId == userID {
			b := r.AllPlayerGameInfo[i].IsFoundDead
			r.AllPlayerGameInfo[i].IsFoundDead = true
			return b
		}
	}
	return false
}

func (r *roomGamingInfo) CheckPlayerGameStatus(userID int32, gameStatus proto3.PlayerGameStatus) bool {
	userInfo := r.GetPlayerGameInfo(userID)
	if userInfo == nil {
		return false
	}
	if userInfo.UserId == userID && userInfo.PlayerGameStatus == gameStatus {
		return true
	}
	return false
}
func (r *roomGamingInfo) ChangePlayerRole(userID int32, status proto3.PlayerRoleEnum) bool {
	userInfo := r.GetPlayerGameInfo(userID)
	if userInfo == nil {
		return false
	}
	return status == userInfo.PlayerRole
}

func (r *roomGamingInfo) ChangePlayerGameStatus(userID int32, status proto3.PlayerGameStatus) {
	userInfo := r.GetPlayerGameInfo(userID)
	if userInfo == nil {
		return
	}
	userInfo.PlayerGameStatus = status
	userInfo.TurnIndex = r.TurnIndex
	if status != proto3.PlayerGameStatus_normal {
		r.deleteLiveUserList(userID)
	}
}

func (r *roomGamingInfo) deleteLiveUserList(userID int32) {
	for i := 0; i < len(r.LiveUserList); i++ {
		if r.LiveUserList[i].UserId == userID {
			r.LiveUserList = append(r.LiveUserList[:i], r.LiveUserList[i+1:]...)
		}
	}
}

func (r *roomGamingInfo) GetLivePlayerNum() int32 {
	return int32(len(r.LiveUserList))
}

// param1-people num, param2-wolf num
func (r *roomGamingInfo) GetLiveRoleNum() (int32, int32) {
	var peopleNum, wolfNum int32
	for i := 0; i < len(r.LiveUserList); i++ {
		if r.LiveUserList[i].PlayerRole == proto3.PlayerRoleEnum_comm_people {
			peopleNum++
		}
		if r.LiveUserList[i].PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			wolfNum++
		}
	}
	return peopleNum, wolfNum
}

func (r *roomGamingInfo) VoteKillPlayer(killUserID int32) {
	r.ChangePlayerGameStatus(killUserID, proto3.PlayerGameStatus_killed_vote)
}

func (r *roomGamingInfo) playerRole(langArr []int, index int) proto3.PlayerRoleEnum {
	for _, v := range langArr {
		if index == v {
			return proto3.PlayerRoleEnum_wolf_man
		}
	}
	return proto3.PlayerRoleEnum_comm_people
}

func (r *roomGamingInfo) SetGameWaitInfo() {
	r.ServerCurrentTime = time.Now().Unix()
	r.GameWaitEndTime = r.ServerCurrentTime + int64(r.GameWaitTime)
}

func (r *roomGamingInfo) SetInitiateVote(userID, sufferID int32, voteType proto3.VoteTypeEnum) {
	r.RoomStatus = proto3.RoomStatus_voting
	r.ServerCurrentTime = time.Now().Unix()
	// 投票前聊天结束时间
	voteChatTime := tableconfig.ConstsConfigs.GetChatVoteTime()
	r.VoteChatEndTime = r.ServerCurrentTime + int64(voteChatTime)

	// 投票结束时间
	r.VoteEndTime = r.VoteChatEndTime + int64(r.GameVoteTime)

	// 汇总结束时间
	voteSumTime := tableconfig.ConstsConfigs.GetVoteSumEndTime()
	r.VoteSumEndTime = r.VoteEndTime + int64(voteSumTime)

	// 初始化投票信息
	r.VoteInfo = voteInfo{
		userVoteMap: make(map[int32]PlayerVote),
	}
	r.VoteInfo.VoteResultResp = proto3.VoteResultResp{}
	r.VoteInfo.VoteResultResp.UserId = userID
	r.VoteInfo.VoteResultResp.SuffererId = sufferID
	r.VoteInfo.VoteResultResp.PlayerGainVoteList = make([]*proto3.PlayerGainVote, 0)
	r.VoteInfo.voteType = voteType
	r.VoteInfo.voteNum = 0
	r.VoteInfo.turnIndex = r.TurnIndex
	for i := 0; i < len(r.AllPlayerGameInfo); i++ {
		v := r.AllPlayerGameInfo[i]
		tmp := &proto3.PlayerGainVote{
			UserId: v.UserId,
		}
		if userID == v.UserId {
			v.CallVoteNum++
			if voteType == proto3.VoteTypeEnum_urgent_vote {
				v.UrgencyVoteNum--
			}
		}
		r.VoteInfo.PlayerGainVoteList = append(r.VoteInfo.PlayerGainVoteList, tmp)
	}
	logger.Log.Infof("init vote info:%v playerGainVote:%v", r.AllPlayerGameInfo, r.VoteInfo.PlayerGainVoteList)
	r.VoteCount++
}

func (r *roomGamingInfo) GameEndWolfSettle(wolfID, turnIndex int32, ret *proto3.GameEndResp) (int32, int32) {
	ret = &proto3.GameEndResp{}
	var voteCount, hidVote int32 = 0, 0
	for _, v := range r.VoteInfoHis {
		if v.JudgeWolfHide(wolfID) && v.turnIndex < turnIndex {
			hidVote++
		}
		if _, ok := v.userVoteMap[wolfID]; ok || (v.CheckUserVoted(wolfID) && v.turnIndex < turnIndex) {
			voteCount++
		}
	}
	return voteCount, hidVote
}

func (r *roomGamingInfo) GameEndPeopleSettle(userID, turnIndex int32, ret *proto3.GameEndResp) (int32, int32) {
	ret = &proto3.GameEndResp{}
	voteCount, successVote, _ := r.GetVoteNum(userID, turnIndex)
	return voteCount, successVote

}

func (r *roomGamingInfo) GetVoteNum(userID, turnIndex int32) (voteCount, voteSuccess, voteFailed int32) {
	for _, v := range r.VoteInfoHis {
		b1, b2 := v.JudgePlayerVoteWolf(userID)
		if b1 && v.turnIndex < turnIndex {
			voteCount++
		}
		if b2 && r.IsWolfMan(v.maxVoteUserId) {
			voteSuccess++
		}
		if b2 && !r.IsWolfMan(v.maxVoteUserId) {
			voteFailed++
		}
	}
	return
}

func (r *roomGamingInfo) GetKillNum(userID int32) int32 {
	for i, v := range r.AllPlayerGameInfo {
		if userID == v.UserId {
			return r.AllPlayerGameInfo[i].KillNum
		}
	}
	return 0
}
func (r *roomGamingInfo) AddKillNum(attackID int32, sufferID int32) {
	if r.FirstKillUserID <= 0 {
		r.FirstKillUserID = sufferID
	}
	for i := range r.AllPlayerGameInfo {
		v := r.AllPlayerGameInfo[i]
		if attackID == v.UserId {
			r.AllPlayerGameInfo[i].KillNum++
		}
		if sufferID == v.UserId {
			r.AllPlayerGameInfo[i].TurnIndex = r.TurnIndex
		}
		if attackID == v.UserId && v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			nowTime := util.UnixTime()
			v.nextAttackTime = nowTime + tableconfig.ConstsConfigs.GetAttackFrozen() - constant.WolfManCdAllowRange // 减去两秒比前端提前cd
		}
	}
}

func (r *roomGamingInfo) JudgeGameEndByRoleNum() (proto3.PlayerRoleEnum, bool) {
	if util.UnixTime() >= int32(r.GameWaitEndTime)+global.MyConfig.ReadInt32("room", "game_length_time") {
		r.GameWinRole = proto3.PlayerRoleEnum_wolf_man
		return proto3.PlayerRoleEnum_wolf_man, true
	}
	wolfLiveNum, peopleLiveNum := 0, 0
	voteKillNum := 0
	for _, v := range r.AllPlayerGameInfo {
		if v.PlayerGameStatus == proto3.PlayerGameStatus_normal {
			if v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
				wolfLiveNum++
			} else {
				peopleLiveNum++
			}
		} else if v.PlayerGameStatus == proto3.PlayerGameStatus_killed_vote && v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			voteKillNum++
		}
	}

	// 当杀手被全部投票出局后，平民玩家获得胜利
	if wolfLiveNum <= 0 {
		r.GameWinRole = proto3.PlayerRoleEnum_comm_people
		logger.Log.Info("is stop because wolf num <= 0")
		return proto3.PlayerRoleEnum_comm_people, true
	}

	if peopleLiveNum <= 0 {
		r.GameWinRole = proto3.PlayerRoleEnum_wolf_man
		logger.Log.Info("is stop because people num <= 0")
		return proto3.PlayerRoleEnum_wolf_man, true
	}

	return proto3.PlayerRoleEnum_invalid_role, false
}

func (r *roomGamingInfo) ArchiveVoteInfo() {
	vi := new(voteInfo)
	*vi = r.VoteInfo
	r.VoteInfoHis = append(r.VoteInfoHis, vi)
	r.VoteInfo = voteInfo{}
	r.TurnIndex++
}

func (r *roomGamingInfo) IsInRoom(userID int32) bool {
	if r.RoomStatus == proto3.RoomStatus_game_end {
		return false
	}
	for i := range r.AllPlayerGameInfo {
		v := r.AllPlayerGameInfo[i]
		if v.UserId == userID && v.PlayerAutoOut == autoOutNormal {
			return true
		}
	}
	return false
}

func (r *roomGamingInfo) GetUrgencyVoteNum(userID int32) int32 {
	userInfo := r.GetPlayerGameInfo(userID)
	if userInfo == nil {
		return -1
	}
	return userInfo.UrgencyVoteNum

	// for i := range r.AllPlayerGameInfo {
	// 	v := r.AllPlayerGameInfo[i]
	// 	if v.UserId == userID {
	// 		return v.UrgencyVoteNum
	// 	}
	// }
	// return 0
}

// 获取狼人杀人CD
func (r *roomGamingInfo) GetWolfManCd() (wolfManCds []*proto3.WolfManCd) {
	var cdEndTime int32
	for _, v := range r.AllPlayerGameInfo {
		if v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			if v.nextAttackTime != 0 {
				cdEndTime = v.nextAttackTime + constant.WolfManCdAllowRange
			}
			wolfManCd := &proto3.WolfManCd{UserId: v.UserId, CdEndTime: cdEndTime}
			wolfManCds = append(wolfManCds, wolfManCd)
		}
	}
	return wolfManCds
}

func (r *roomGamingInfo) GetLampList() (ret []*proto3.UserLamp) {
	ret = make([]*proto3.UserLamp, 0)
	for k, v := range r.LampUserMap {
		tmp := new(proto3.UserLamp)
		tmp.UserId = k
		tmp.Status = proto3.CommonStatusEnum(v)
		ret = append(ret, tmp)
	}
	return
}
