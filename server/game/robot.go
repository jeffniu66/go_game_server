package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
)

func (m *Room) RobotStartVote() {
	logger.Log.Info("robotVote Start>>>>>")
	for _, v := range m.roomInfo.RobotUserList {
		if v.PlayerGameStatus != proto3.PlayerGameStatus_normal { // 死亡不进行投票
			continue
		}
		var t2 int32 = 15 * 1000         // 机器人投票[15,45]
		t2 += util.RandInt(0, 30) * 1000 // ms
		m.roomPid.SendAfter("robotVote", "robotVote_"+util.ToStr(v.UserId), t2, v.UserId)
	}
}

func (m *Room) RobotVote(robotID int32) {
	var votedId int32
	var voteType proto3.CommonStatusEnum = 2 // 0-跟投，1-随机投，2-弃票 // 默认弃票
	voteMsg := new(include.VoteMsg)

	playerInfo := m.roomInfo.GetPlayerGameInfo(robotID)
	if playerInfo.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
		votedUser := m.roomInfo.GetALiveRankPeopleUser()
		votedId = votedUser.UserId
		voteType = proto3.CommonStatusEnum_pass // 0
	} else {
		var voteTList []int32
		var index int32
		peopleNum, _ := m.roomInfo.GetLiveRoleNum()
		if peopleNum >= 3 {
			if peopleNum >= 6 {
				voteTList = []int32{0, 1, 2, 2} // 投票概率 跟头0-25% 随机投1-25% 弃票2-50%
				index = util.RandInt(0, 3)
			}
			if peopleNum < 6 && peopleNum >= 3 {
				voteTList = []int32{0, 1} // 投票概率 跟头0-50% 随机投1-50%
				index = util.RandInt(0, 1)
			}
			voteType = proto3.CommonStatusEnum(voteTList[index])
			if voteType == proto3.CommonStatusEnum_true { // 随机投
				voteUser := m.roomInfo.GetALiveAiRankUserNoMe(robotID)
				votedId = voteUser.UserId
			} else if voteType == proto3.CommonStatusEnum_pass { // 跟投
				if len(m.roomInfo.VoteInfo.userVoteMap) <= 0 {
					voteUser := m.roomInfo.GetALiveAiRankUserNoMe(robotID)
					votedId = voteUser.UserId
				} else {
					for _, v := range m.roomInfo.VoteInfo.userVoteMap {
						votedId = v.userId
						break
					}
				}
			} else {
				// 弃票
			}
			playerInfo := m.roomInfo.GetPlayerGameInfo(votedId)
			if (playerInfo != nil && playerInfo.IsRobot != int32(proto3.CommonStatusEnum_true)) || votedId == robotID {
				voteType = 2 // 弃票
				votedId = 0
			}
		}
		if peopleNum < 3 {
			voteType = proto3.CommonStatusEnum_true // 投票
			votedId = m.roomInfo.GetALiveWolfUserID()
		}
	}

	voteMsg.UserId = robotID
	voteMsg.VoteUserId = votedId
	voteMsg.Status = voteType
	go m.roomPid.Cast("playerVote", voteMsg)
}

func (m *Room) robotChat() {
	if len(m.roomInfo.RobotUserList) <= 0 {
		return
	}

	t := tableconfig.ConstsConfigs.GetRobotChatTime() * 1000
	m.roomPid.SendAfter("robotChat", "robotChat", t, nil)

	// i := util.RandInt(0, int32(len(m.roomInfo.RobotUserList))-1)
	// userID := m.roomInfo.RobotUserList[i]
	userInfo := m.roomInfo.GetALiveAiRankUser()
	userID := userInfo.UserId
	if !m.roomInfo.CheckPlayerGameStatus(userID, proto3.PlayerGameStatus_normal) {
		return
	}
	if m.roomInfo.RoomStatus == proto3.RoomStatus_wait_game {
		quickConfig := tableconfig.QuickTextCols.GetRandText(constant.TextTypeWait, m.roomInfo.GetLivePlayerNum())
		m.SendWaitChat(userID, "", 1, quickConfig.Id)
	}
	if m.roomInfo.RoomStatus == proto3.RoomStatus_voting {
		quickConfig := tableconfig.QuickTextCols.GetRandText(constant.TextTypeVote, m.roomInfo.GetLivePlayerNum())
		subsText := make([]string, 0)
		for i := 0; i < int(quickConfig.DataNumber); i++ {
			info := m.roomInfo.GetALiveAiRankUserNoMe(userID)
			subsText = append(subsText, util.ToStr(info.BirthPoint))
		}
		if len(subsText) <= 0 {
			m.SendChat(userID, "", 1, quickConfig.Id, nil)
		} else {
			m.SendChat(userID, "", 1, quickConfig.Id, subsText)
		}
	}
}
