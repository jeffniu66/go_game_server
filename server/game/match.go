package game

import (
	"go_game_server/server/logger"
)

type userRank struct {
	UserID    int32
	StarLevel int32 // 玩家星级
	IsRobot   int32
	Channel   string // 通道
}

type chanMatch struct {
	userRankList []userRank
	maxStarLevel int32 // 本段位最大星数
	minStarLevel int32 // 本段位最小星数
	mergeNum     int32 // 合并次数
}

type rankMatchUser struct {
	userMap map[string]*chanMatch // key-channel

}

func newRankMatchUser(minStarLevel, maxStarLevel int32) *rankMatchUser {
	ret := new(rankMatchUser)
	userMap := make(map[string]*chanMatch)
	userMap["dev"] = &chanMatch{minStarLevel: minStarLevel, maxStarLevel: maxStarLevel}
	userMap["wx"] = &chanMatch{minStarLevel: minStarLevel, maxStarLevel: maxStarLevel}
	userMap["qq"] = &chanMatch{minStarLevel: minStarLevel, maxStarLevel: maxStarLevel}
	userMap["oppo"] = &chanMatch{minStarLevel: minStarLevel, maxStarLevel: maxStarLevel}
	ret.userMap = userMap
	return ret
}

func (r *rankMatchUser) GetRobotNum(ch string) (ret int32) {
	mat, ok := r.userMap[ch]
	if !ok {
		return 0
	}
	for _, v := range mat.userRankList {
		if v.IsRobot == 1 {
			ret++
		}
	}
	return
}

func (r *rankMatchUser) SetMergeZero(ch string) {
	mat := r.userMap[ch]
	mat.mergeNum = 0
	r.userMap[ch] = mat
}

func (r *rankMatchUser) IsExist(userID int32, c string) bool {
	if v, ok := r.userMap[c]; ok {
		return v.IsExist(userID)
	}
	return false
}

func (r *rankMatchUser) GetOtherChannel(c string) string {
	if v, ok := r.userMap[c]; ok && len(v.userRankList) == 0 {
		for k, vv := range r.userMap {
			if k != c && len(vv.userRankList) > 0 {
				return k
			}
		}
	}
	return ""
}

func (r *rankMatchUser) addPlayerID(userID int32, StarLevel, isRobot int32, c string) bool {
	// userR := userRank{UserID: userID, StarLevel: StarLevel, IsRobot: isRobot, Channel: c}
	if !r.IsExist(userID, c) {
		r.userMap[c].addPlayerID(userID, StarLevel, isRobot, c)
		return true
	}
	return false
}

func (r *rankMatchUser) GetPlayerIDList(c string) (ret []int32) {
	userList := r.userMap[c]
	ret = userList.GetPlayerIDList()
	return
}

func (r *rankMatchUser) ClearRankList(c string) {
	userList := r.userMap[c]
	userList.ClearRankList()
}

func (r *rankMatchUser) deletePlayerID(userID int32, c string) {
	chanMatch := r.userMap[c]
	chanMatch.deletePlayerID(userID)
}

func (r *rankMatchUser) GetPlayerIDListByStarLevel(starLevel int32, needNum int32, isNext bool, ch string) []int32 {
	if uL, ok := r.userMap[ch]; ok {
		return uL.GetPlayerIDListByStarLevel(starLevel, needNum, isNext)
	}
	return []int32{}
}

func (r *rankMatchUser) GetPlayerIDListByChannel(ch string, matchNum int32) ([]int32, bool) {
	logger.Log.Info("start merge channel:", ch)

	ret := make([]int32, 0)
	mat := r.userMap[ch]
	for _, vv := range mat.userRankList {
		ret = append(ret, vv.UserID) // 取一个返回
		break
	}
	needNum := matchNum - int32(len(mat.userRankList))
	for k, v := range r.userMap {
		if k == ch {
			continue
		}
		if needNum == 0 {
			return ret, true
		}
		tmp := v.PopUser(needNum)
		needNum -= int32(len(tmp))
		mat.userRankList = append(mat.userRankList, tmp...)
	}
	r.userMap[ch] = mat
	return ret, false
}

func (c *chanMatch) IsExist(userID int32) bool {
	for i := range c.userRankList {
		if c.userRankList[i].UserID == userID {
			return true
		}
	}
	return false
}

func (c *chanMatch) GetPlayerIDListByStarLevel(starLevel int32, needNum int32, isNext bool) []int32 {
	ret := []int32{}
	for i := range c.userRankList {
		if isNext {
			if c.userRankList[i].StarLevel <= starLevel {
				ret = append(ret, c.userRankList[i].UserID)
				if int32(len(ret)) == needNum {
					break
				}
			}
		} else {
			if c.userRankList[i].StarLevel >= starLevel {
				ret = append(ret, c.userRankList[i].UserID)
				if int32(len(ret)) == needNum {
					break
				}
			}
		}
	}
	return ret
}

func (c *chanMatch) addPlayerID(userID int32, StarLevel, isRobot int32, ch string) bool {
	logger.Log.Infof("add playerID u:%v playerID:%d exist:%v", c, userID, c.IsExist(userID))
	userR := userRank{UserID: userID, StarLevel: StarLevel, IsRobot: isRobot, Channel: ch}
	if !c.IsExist(userID) {
		c.userRankList = append(c.userRankList, userR)
		logger.Log.Infof("add playerID u:%v playerID:%d exist:%v", c, userID, c.IsExist(userID))
		return true
	}
	return false
}

func (c *chanMatch) deletePlayerID(userID int32) {
	tmp := []userRank{}
	robotNum := 0
	for i := range c.userRankList {
		if c.userRankList[i].UserID != userID {
			tmp = append(tmp, c.userRankList[i])
		}
		if c.userRankList[i].IsRobot == 1 {
			robotNum++
		}
	}
	c.userRankList = tmp
	// 全是机器人情况房间
	if robotNum == len(c.userRankList) {
		c.ClearRankList()
	}
}

func (c *chanMatch) ClearRankList() {
	c.mergeNum = 0
	c.userRankList = make([]userRank, 0)
}

func (c *chanMatch) GetPlayerIDList() (ret []int32) {
	for i := range c.userRankList {
		ret = append(ret, c.userRankList[i].UserID)
	}
	return
}

func (c *chanMatch) PopUser(needNum int32) []userRank {
	ret := make([]userRank, 0)
	i := 0
	for ; i < int(needNum); i++ {
		if i < len(c.userRankList) {
			// tmp := userRank{}
			// tmp = c.userRankList[i]
			ret = append(ret, c.userRankList[i])
		} else {
			break
		}
	}
	if i == len(c.userRankList)-1 {
		c.userRankList = make([]userRank, 0)
	} else {
		c.userRankList = c.userRankList[i:]
	}
	c.mergeNum = 0
	return ret
}
