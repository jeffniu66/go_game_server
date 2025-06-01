package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/include"
	"go_game_server/server/logger"
)

type PlayerVote struct {
	userId  int32 // 被投票人
	voteNum int32 // 得票数
}

type voteInfo struct {
	// votedUserIds []int32         // 已投票列表
	voteNum       int32                // 投票人数包含弃票
	userVoteMap   map[int32]PlayerVote // key-投票者 value-被投票者 一个玩家的投票 key-userID 无器票
	maxVoteUserId int32                // 最多票玩家
	maxNum        int32                // 得票
	voteType      proto3.VoteTypeEnum  // 会议类型
	turnIndex     int32                // 当前投票第几轮
	proto3.VoteResultResp
}

func (v *voteInfo) GetVoteResult() *proto3.VoteResultResp {
	return &v.VoteResultResp
}

func (v *voteInfo) CheckUserVoted(userID int32) bool {
	if _, ok := v.userVoteMap[userID]; ok {
		return ok
	}
	for _, v := range v.LoseVoteList {
		if userID == v {
			return true
		}
	}
	return false
}

// param1 参与投票， param2 投票命中
func (v *voteInfo) JudgePlayerVoteWolf(playerId int32) (bool, bool) {
	// 判断用户是否投对
	u, ok := v.userVoteMap[playerId]
	if !ok {
		b := v.CheckUserVoted(playerId)
		return b, false
	}

	return true, u.userId == v.maxVoteUserId
}

// 判断狼人是否闪避
func (v *voteInfo) JudgeWolfHide(wolfUserId int32) bool {
	for _, vv := range v.userVoteMap {
		if vv.userId == wolfUserId {
			return false
		}
	}
	return true
}

func (v *voteInfo) GetVotedResult(liveNum int32) (killUserId int32, voteNum int32, status proto3.VoteStatus) {
	if v.maxNum > 0 {
		return v.maxVoteUserId, v.maxNum, v.VoteStatus
	}
	logger.Log.Infof("vote:%v", v.userVoteMap)

	tmpVoteMap := make(map[int32]int32) // key-userID value-得票数
	logger.Log.Infof("userVoteMap:%v", v.userVoteMap)
	for userID, playerVote := range v.userVoteMap {
		// playerVote 玩家得票
		if n, ok := tmpVoteMap[playerVote.userId]; ok {
			n += playerVote.voteNum
			tmpVoteMap[playerVote.userId] = n
		} else {
			tmpVoteMap[playerVote.userId] = playerVote.voteNum
		}
		// 计算玩家获得的投票
		// 注意指针之间的赋值
		for i := range v.PlayerGainVoteList {
			// gainVote *PlayerGainVote
			gainVote := v.PlayerGainVoteList[i]
			if gainVote.UserId == playerVote.userId {
				// 使用地址直接赋值
				tmp := proto3.VotedUserList{
					UserId:  userID,
					VoteNum: playerVote.voteNum,
				}
				v.PlayerGainVoteList[i].VoteUserList = append(v.PlayerGainVoteList[i].VoteUserList, &tmp)
				break
			}
		}
	}
	// 过半为有效票
	if int32(len(v.userVoteMap)) <= liveNum/2 {
		v.VoteStatus = proto3.VoteStatus_invalid_no_helf
		return -1, -1, proto3.VoteStatus_invalid_no_helf
	}

	// 查找最大票人
	for k, vv := range tmpVoteMap {
		if int32(vv) > voteNum {
			killUserId = k
			voteNum = int32(vv)
		}
	}
	logger.Log.Infof("tmpVoteMap:%v", tmpVoteMap)
	// 判断是否平票
	for k, vv := range tmpVoteMap {
		if killUserId != k && voteNum == int32(vv) {
			v.VoteStatus = proto3.VoteStatus_invalid_equal_vote
			return -1, -1, proto3.VoteStatus_invalid_equal_vote
		}
	}
	v.maxVoteUserId = killUserId
	v.maxNum = voteNum
	status = proto3.VoteStatus_valid_vote
	v.VoteStatus = status

	RecordAction(killUserId, constant.VoteKill)
	return
}

func (v *voteInfo) AddVoteInfo(msg *include.VoteMsg) {
	v.voteNum++
	if msg.Status == proto3.CommonStatusEnum_false {
		// 弃票
		v.LoseVoteList = append(v.LoseVoteList, msg.UserId)

		RecordAction(msg.UserId, constant.GiveUpVote)
		return
	}
	playerVote := PlayerVote{userId: msg.VoteUserId, voteNum: msg.VoteNum}
	v.userVoteMap[msg.UserId] = playerVote
	v.UserVoteList = append(v.UserVoteList, msg.UserId)

	RecordAction(msg.UserId, constant.Vote)
}

func (v *voteInfo) GetVoteNum() int32 {
	return v.voteNum
}
