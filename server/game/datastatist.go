package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/util"
	"time"
)

// 记录游戏局数
func (r *Room) RecordGameData() {
	if r.roomInfo == nil {
		return
	}
	if len(r.roomInfo.RobotUserList) >= constant.AIMaxNum { // 机器人数量
		return
	}
	dateStr := util.GetDateStr(time.Now().Unix())
	for i := 0; i < len(r.roomInfo.AllPlayerGameInfo); i++ {
		v := r.roomInfo.AllPlayerGameInfo[i]
		if v.IsRobot == 1 {
			continue
		}
		if v.PlayerAutoOut == autoOutNormal && v.Player != nil && v.Player.Attr.RoomID == r.roomId {
			gameDate := db.SelectGameData(v.UserId, dateStr)
			if gameDate == nil {
				gameDate = &include.GameData{UserId: v.UserId, Username: v.Username, GameNum: 1, LoginDate: dateStr, RegisterTime: v.Player.Attr.RegisterTime}
				db.InsertGameData(gameDate)
			} else {
				gameDate.GameNum++
				db.UpdateGameData(gameDate)
			}
		}
	}
}

// 记录登录数据
func RecordLogin(p *Player) {
	dateStr := util.GetDateStr(time.Now().Unix())
	loginData := db.SelectLoginData(p.Attr.UserID, dateStr)
	if loginData != nil {
		return
	}
	loginData = &include.LoginData{UserId: p.Attr.UserID, Username: p.Attr.Username, LoginDate: dateStr, RegisterTime: p.Attr.RegisterTime}
	db.InsertLoginData(loginData)
}

// 记录广告数据
func RecordAdData(p *Player, adType int32) {
	dateStr := util.GetDateStr(time.Now().Unix())
	adData := db.SelectAdData(p.Attr.UserID, adType, dateStr)
	if adData == nil {
		adData = &include.AdData{UserId: p.Attr.UserID, AdType: adType, AdNum: 1, RegisterTime: p.Attr.RegisterTime, StatDate: dateStr}
		db.InsertAdData(adData)
	} else {
		adData.AdNum++
		db.UpdateAdData(adData)
	}
}

// 记录阵营胜率数据
func (r *Room) RecordGroupWinData() {
	if r.roomInfo == nil {
		return
	}
	if len(r.roomInfo.RobotUserList) >= constant.AIMaxNum { // 机器人数量
		return
	}
	var wolfManWin, normalManWin int32
	if r.roomInfo.GameWinRole == proto3.PlayerRoleEnum_wolf_man {
		wolfManWin = 1
	} else {
		normalManWin = 1
	}
	dateStr := util.GetDateStr(time.Now().Unix())
	groupWinData := db.SelectGroupWinData(dateStr)
	if groupWinData == nil {
		groupWinData = &include.GroupWinData{
			WolfMan:   wolfManWin,
			NormalMan: normalManWin,
			StatData:  dateStr,
		}
		db.InsertGroupWinData(groupWinData)
	} else {
		if wolfManWin == 1 {
			groupWinData.WolfMan++
		} else {
			groupWinData.NormalMan++
		}
		db.UpdateGroupWinData(groupWinData)
	}
}

// 记录狼人杀人数据
func (r *Room) RecordKillNum(num int32) {
	if r.roomInfo == nil {
		return
	}
	if len(r.roomInfo.RobotUserList) >= constant.AIMaxNum {
		return
	}
	allPlayerInfo := r.roomInfo.AllPlayerGameInfo
	if len(allPlayerInfo) == 0 {
		return
	}
	dateStr := util.GetDateStr(time.Now().Unix())
	for _, v := range allPlayerInfo {
		if v.IsRobot == int32(proto3.CommonStatusEnum_true) {
			continue
		}
		if v.PlayerRole == proto3.PlayerRoleEnum_wolf_man {
			killData := db.SelectKillData(v.UserId, dateStr)
			if killData == nil {
				killData = &include.KillData{
					UserId:   v.UserId,
					RankId:   v.Player.Attr.RankID,
					KillNum:  num,
					StatDate: dateStr,
				}
				db.InsertKillData(killData)
			} else {
				if num == 0 {
					continue
				}
				killData.KillNum += num
				db.UpdateKillData(killData)
			}
		}
	}
}

// 记录会议(投票)次数
func (r *Room) RecordConferNum() {
	if r.roomInfo == nil {
		return
	}
	if len(r.roomInfo.RobotUserList) >= constant.AIMaxNum {
		return
	}
	conferMap := make(map[int32]int32) // key:userId value:num
	votes := r.roomInfo.VoteInfoHis
	for _, v := range votes {
		for _, voteUserId := range v.UserVoteList {
			conferMap[voteUserId]++
		}
		for _, loseUserId := range v.LoseVoteList {
			conferMap[loseUserId]++
		}
	}
	dateStr := util.GetDateStr(time.Now().Unix())
	for k, v := range conferMap {
		p := global.GloInstance.GetPlayer(k)
		if p == nil {
			continue
		}
		player := p.(*Player)
		if player.IsRobot == 1 {
			continue
		}
		conferData := &include.ConferData{
			UserId:       player.Attr.UserID,
			RankId:       player.Attr.RankID,
			ConferNum:    v,
			StatDate:     dateStr,
			RegisterDate: player.Attr.RegisterTime,
		}
		db.InsertConferData(conferData)
	}
}

// 记录道具获得与消耗
func RecordItem(itemDatas []*include.ItemData) {
	if len(itemDatas) == 0 {
		return
	}
	db.BatchInsertItem(itemDatas)
}
