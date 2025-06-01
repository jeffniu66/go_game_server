package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
)

type Skill struct {
	forbidChatMap     map[int32]int32 // 禁言符 key: 使用者 value: 被使用者
	doubleVoteUserIds []int32         // 双倍票玩家使用者
	guardUserIds      []int32         // 守护符使用者
	boomMap           map[int32]int32 // 自爆符 key: 使用者 value: 被使用者
	togetherMap       map[int32]int32 // 比翼鸟 key: 使用者 value: 被使用者
}

// 使用技能
func (r *Room) UseSkill(itemId int32, attackUserId int32, attackedUserId int32) (role proto3.PlayerRoleEnum, randRoomId int32) {
	skill := r.skill
	switch itemId {
	case constant.Skill10008: // 现行符
		if r.roomInfo.IsWolfMan(attackUserId) {
			return proto3.PlayerRoleEnum_wolf_man, randRoomId
		}
		return proto3.PlayerRoleEnum_comm_people, randRoomId
	case constant.Skill10009: // 妖刀杀人
		if r.roomInfo.CheckPlayerGameStatus(attackedUserId, proto3.PlayerGameStatus_normal) {
			r.roomInfo.ChangePlayerGameStatus(attackedUserId, proto3.PlayerGameStatus_killed)
			r.roomInfo.AddKillNum(attackUserId, attackedUserId)
			// 判断游戏是否结束
			_, b := r.roomInfo.JudgeGameEndByRoleNum()
			if b {
				r.roomPid.Cast("gameEnd", nil)
			} else { // 开启投票
				r.roomPid.Cast("startVote", attackUserId)
			}
		}
	case constant.Skill10010: // 屏蔽聊天
		if skill == nil {
			tempMap := make(map[int32]int32)
			skill = &Skill{forbidChatMap: tempMap}
		} else {
			skill.forbidChatMap[attackUserId] = attackedUserId
		}
	case constant.Skill10011: // 自爆符
		if skill == nil {
			tempMap := make(map[int32]int32)
			tempMap[attackUserId] = attackedUserId
			skill = &Skill{boomMap: tempMap}
		} else {
			skill.boomMap[attackUserId] = attackedUserId
		}
	case constant.Skill10013: // 双倍票数
		if skill == nil {
			skill = &Skill{doubleVoteUserIds: []int32{attackUserId}}
		} else {
			skill.doubleVoteUserIds = append(skill.doubleVoteUserIds, attackUserId)
		}
	case constant.Skill10014: // 守护符
		if skill == nil {
			skill = &Skill{guardUserIds: []int32{attackUserId}}
		} else {
			skill.guardUserIds = append(skill.guardUserIds, attackUserId)
		}
	case constant.Skill10102: // 无影针
		if r.roomInfo.CheckPlayerGameStatus(attackedUserId, proto3.PlayerGameStatus_normal) {
			r.roomInfo.ChangePlayerGameStatus(attackedUserId, proto3.PlayerGameStatus_killed)
			r.roomInfo.AddKillNum(attackUserId, attackedUserId)
			// 判断游戏是否结束
			_, b := r.roomInfo.JudgeGameEndByRoleNum()
			if b {
				r.roomPid.Cast("gameEnd", nil)
			} else {
				// 客户端要求去掉下发攻击返回
				//player := global.GloInstance.GetPlayer(attackUserId).(*Player)
				//cmd := proto3.ProtoCmd_CMD_NinjaAttackResp
				//pbData := &proto3.NinjaAttackResp{
				//	NinjaStatus: proto3.CommonStatusEnum_true,
				//	UserId:      attackUserId,
				//	SuffererId:  attackedUserId,
				//}
				//player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
			}
		}
	case constant.Skill10105: // 比翼鸟
		if skill == nil {
			tempMap := make(map[int32]int32)
			tempMap[attackUserId] = attackedUserId
			skill = &Skill{togetherMap: tempMap}
		} else {
			skill.togetherMap[attackUserId] = attackedUserId
		}
	case constant.Skill10003: // 烟雾弹
		randRoomId = util.RandInt(1, util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.SmokeEgg)))
	default:
		break
	}
	r.skill = skill
	return 0, randRoomId
}

// 处理双倍票符
func (r *Room) DealDoubleVote() {
	for i := 0; i < len(r.roomInfo.AllPlayerGameInfo); i++ {
		userId := r.roomInfo.AllPlayerGameInfo[i].UserId
		if r.isUseDoubleVoteSkill(userId) {
			userVoteMap := r.roomInfo.VoteInfo.userVoteMap
			playerVote := userVoteMap[userId]
			if playerVote.userId == 0 {
				continue
			}
			playerVote.voteNum++
			userVoteMap[userId] = playerVote
			r.roomInfo.VoteInfo.userVoteMap = userVoteMap
		}
	}
}

// 是否使用过双倍票技能
func (r *Room) isUseDoubleVoteSkill(userId int32) bool {
	skill := r.skill
	if skill == nil {
		return false
	}
	doubleVoteUserIds := skill.doubleVoteUserIds
	if doubleVoteUserIds == nil {
		return false
	}
	for i, doubleVoteUserId := range doubleVoteUserIds {
		if doubleVoteUserId == userId {
			doubleVoteUserIds = append(doubleVoteUserIds[:i], doubleVoteUserIds[i+1:]...)
			r.skill.doubleVoteUserIds = doubleVoteUserIds
			return true
		}
	}
	return false
}

// 处理使用过自爆符情况 反杀
func (r *Room) DealBoom(attackUserId int32, attackedUserId int32) {
	// 判断反杀
	skill := r.skill
	if skill == nil {
		return
	}
	boomUserId, ok := skill.boomMap[attackedUserId]
	if !ok {
		return
	}
	var flag bool
	if boomUserId == attackUserId {
		flag = true
	}
	if flag {
		r.roomInfo.ChangePlayerGameStatus(attackUserId, proto3.PlayerGameStatus_killed)
		r.roomInfo.AddKillNum(attackedUserId, attackUserId)
		// 判断游戏是否结束
		_, b := r.roomInfo.JudgeGameEndByRoleNum()
		if b {
			r.roomPid.Cast("gameEnd", nil)
		}
	}
}

// 获取被自爆符使用过的玩家id
func (r *Room) GetBoomUserId(userId int32) int32 {
	skill := r.skill
	if skill == nil {
		return 0
	}
	boomUserId, ok := skill.boomMap[userId]
	if !ok {
		return 0
	}
	return boomUserId
}

// 处理使用过守护符的玩家
func (r *Room) DealGuard(voteKilledUserId int32) bool {
	skill := r.skill
	if skill == nil {
		return false
	}
	guards := skill.guardUserIds
	if len(guards) == 0 {
		return false
	}
	for i, v := range guards {
		if v == voteKilledUserId {
			guards = append(guards[:i], guards[i+1:]...)
			skill.guardUserIds = guards
			r.skill = skill

			player := global.GloInstance.GetPlayer(voteKilledUserId)
			if player != nil {
				p := player.(*Player)
				skillConfig := tableconfig.SkillConfigs.GetSkillById(constant.Skill10014)
				pbData := &proto3.CommonTextNotice{TextKey: skillConfig.Text, Param: p.Attr.Username}
				r.spreadPlayers(proto3.ProtoCmd_CMD_CommonTextNoticeResp, pbData)
			}
			return true
		}
	}
	return false
}

// 处理比翼鸟
func (r *Room) DealTogether(voteKillUserId int32) {
	if voteKillUserId == -1 {
		return
	}
	skill := r.skill
	if skill == nil {
		return
	}
	togetherMap := skill.togetherMap
	if togetherMap == nil {
		return
	}
	for k, v := range togetherMap {
		if k == voteKillUserId { // 使用者被票杀 被使用者跟着出局
			r.roomInfo.ChangePlayerGameStatus(v, proto3.PlayerGameStatus_killed)
			r.roomInfo.AddKillNum(k, v)
		}
		if v == voteKillUserId { // 被使用者被票杀 使用者跟着出局
			r.roomInfo.ChangePlayerGameStatus(k, proto3.PlayerGameStatus_killed)
			r.roomInfo.AddKillNum(v, k)
		}
	}
}

func (r *Room) reduceSkill(userId int32, skillId int32) proto3.ErrEnum {
	if task, ok := r.taskMap[userId]; ok {
		skills := task.Skill
		if !existsSkill(skills, skillId) {
			// 道具不存在
			return proto3.ErrEnum_Error_Goods_NotExists
		}
		//// 减道具
		//for k, v := range skills {
		//	if v.Id == skillId {
		//		v.Num = v.Num - 1
		//		skills[k] = v
		//	}
		//}
		task.Skill = nil
		r.taskMap[userId] = task
	}
	return proto3.ErrEnum_Error_Pass
}

func existsSkill(skill *proto3.Item, skillId int32) bool {
	if skill == nil {
		return false
	}
	if skill.Id == skillId && skill.Num > 0 {
		return true
	}
	return false
}

// 选择了一种技能，则清除其它技能
func (r *Room) ClearSkills(userId int32) {
	skill := r.skill
	if skill == nil {
		return
	}
	forbidChatMap := skill.forbidChatMap
	if forbidChatMap != nil {
		if _, ok := forbidChatMap[userId]; ok {
			delete(forbidChatMap, userId)
			skill.forbidChatMap = forbidChatMap
		}
	}

	doubleVoteUserIds := skill.doubleVoteUserIds
	if doubleVoteUserIds != nil {
		for i, v := range doubleVoteUserIds {
			if v == userId {
				doubleVoteUserIds = append(doubleVoteUserIds[:i], doubleVoteUserIds[i+1:]...)
				skill.doubleVoteUserIds = doubleVoteUserIds
			}
		}
	}

	guardUserIds := skill.guardUserIds
	if guardUserIds != nil {
		for i, v := range guardUserIds {
			if v == userId {
				guardUserIds = append(guardUserIds[:i], guardUserIds[i+1:]...)
				skill.guardUserIds = guardUserIds
			}
		}
	}

	boomMap := skill.boomMap
	if boomMap != nil {
		if _, ok := boomMap[userId]; ok {
			delete(boomMap, userId)
			skill.boomMap = boomMap
		}
	}
}

// 获取玩家在游戏内获得的技能数
func (r *Room) GetPlayerSkillNum(userId int32) int32 {
	taskMap := r.taskMap
	if taskMap == nil {
		return 0
	}
	task, ok := taskMap[userId]
	if !ok {
		return 0
	}
	return int32(len(task.TotalGetSkills))
}

// 被击杀后清除玩家身上技能
func (r *Room) AttackedClearSkills(userId int32) {
	if r.taskMap == nil {
		return
	}
	task := r.taskMap[userId]
	if task == nil {
		return
	}
	task.Skill = nil
}
