package game

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"runtime/debug"
	"strconv"
)

// PlayerRankType 玩家段位类型
type PlayerRankType = proto3.PlayerRankEnum // 段位枚举

type mergeRankMatchReq struct {
	RankMatch     *rankMatch
	NeedNum       int32 // 需求数量
	RankStarLevel int32 // 需求对于星级
	IsNext        bool
	Channel       string
}

type mergeRankMatchReply struct {
	PlayerIdList []int32
	Channel      string
}

/*
 * rankMatch
 */
type rankMatch struct {
	matchNum      int32          // 策划配置
	rank          PlayerRankType // 段位
	preRankMatch  *rankMatch     // 上一段位
	nextRankMatch *rankMatch     // 下一段位
	// userRankList   *userRankList
	rankMatchPid *global.PidObj
	// maxStarLevel   int32  // 本段位最大星数
	// minStarLevel   int32  // 本段位最小星数
	// mergeNum       int32  // 合并次数
	mergeTimerKey  string // 合并请求定时任务KEY
	mergeTimeout   int32  // 毫秒
	mergeAiTimeout int32  // 毫秒
	rankMatchUser  *rankMatchUser
}

func (r *rankMatch) ClearRankMatch(ch string) {
	// r.userRankList.ClearRankList()
	r.rankMatchUser.ClearRankList(ch)
	// r.mergeNum = 0
	r.rankMatchPid.StopPidTimer(r.mergeTimerKey)
}

func initRankMatch(rank PlayerRankType, minStarLevel, maxStarLevel int32) *rankMatch {
	time1, time2 := tableconfig.QuaLevelConfs.GetMergeTimeout(int32(rank))
	time1 *= 1000
	time2 *= 1000
	rankMatchObj := &rankMatch{
		// userRankList:   &userRankList{},
		// minStarLevel:   minStarLevel,
		// maxStarLevel:   maxStarLevel,
		mergeTimeout:   time1,
		mergeAiTimeout: time2,
	}
	rankMatchObj.rank = rank
	rankMatchObj.rankMatchUser = newRankMatchUser(minStarLevel, maxStarLevel)
	pidName := "rank_match_" + strconv.Itoa(int(rank))
	rankMatchObj.rankMatchPid = global.RegisterPid(pidName, 2048, rankMatchObj)
	return rankMatchObj
}

func (r *rankMatch) Start() {
	logger.Log.Infof("%v rank match mgr start ... ", r.rank)
}

func (r *rankMatch) AddRobot(num int32) []int32 {
	ret := make([]int32, 0)
	for i := int32(0); i < num; i++ {
		ret = append(ret, i+1)
	}
	return ret
}
func (r *rankMatch) HandleCall(req global.GenReq) global.Reply {
	var ret global.Reply
	switch req.Method {

	default:
	}
	return ret
}

func (r *rankMatch) HandleCast(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("matchMgr SOCKET_EVENT HandleCast error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			debug.PrintStack()
			logger.Log.Errorf("rankMatch HandleCast panic-------------rankMatch:%v", r)
		}
	}()
	switch req.Method {
	case "onlyGame":
		player := req.MsgData.(*Player)
		userID := player.Attr.UserID
		playerIDList := []int32{userID}
		// 临时增加机器人
		playerIDList = append(playerIDList, r.AddRobot(3)...)
		roomID := global.GloInstance.GetUseableRoomId()
		if roomID < 0 {
			logger.Log.Error("房间已满进入不了")
			return
		}
		room := InitRoomMgr(roomID, playerIDList, true)
		global.GloInstance.ChangeRoomIdUsed(roomID, room)

		logger.Log.Infof("roomID:%d start, players: %v", roomID, playerIDList)
		room.roomPid.Cast("enterRoom", playerIDList)
		room.roomPid.Cast("gameWait", '0')
	case "enterRankMatch": // 进入匹配
		player := req.MsgData.(*Player)
		userID := player.Attr.UserID

		r.rankMatchUser.addPlayerID(userID, player.Attr.StarCount, player.IsRobot, player.Attr.Channel)

		playerIDList := r.rankMatchUser.GetPlayerIDList(player.Attr.Channel)
		langNum, pingNum := tableconfig.ConstsConfigs.GetMatchGameParam()
		r.matchNum = int32(langNum + pingNum)

		logger.Log.Info("enter match >>>>>>>>> ", playerIDList, player.Attr.UserID, player.Attr.Channel, r.matchNum, r.rank)
		if int32(len(playerIDList)) >= r.matchNum {
			roomID := global.GloInstance.GetUseableRoomId()
			if roomID < 0 {
				logger.Log.Error("房间已满进入不了")
				return
			}
			room := InitRoomMgr(roomID, playerIDList, false)
			global.GloInstance.ChangeRoomIdUsed(roomID, room)

			logger.Log.Infof("roomID:%d start, players: %v", roomID, playerIDList)
			room.roomPid.Cast("enterRoom", playerIDList)
			room.roomPid.Cast("gameWait", '0')
			r.ClearRankMatch(player.Attr.Channel) // 清空匹配器
		} else {
			if int32(len(playerIDList)) > 0 {
				if player.IsRobot != 1 {
					r.rankMatchUser.SetMergeZero(player.Attr.Channel)
				}
				r.sendMergeRankTimer(player.Attr.Channel)
			}
			logger.Log.Infof("add playerID rankUser:%v playerID:%d", r.rankMatchUser.userMap[player.Attr.Channel], userID)
		}
		if player.IsRobot != 1 {
			rankInfo := tableconfig.QuaLevelConfs.GetQuaConfig(player.Attr.RankID)
			t := rankInfo.AiAllTime * 1000
			r.rankMatchPid.SendAfter("rankAddAIALL", "rankAddAIALL_"+player.Attr.Channel, t, player.Attr.Channel)
		}
	case "quitRankMatch":
		player := req.MsgData.(*Player)
		userID := player.Attr.UserID

		r.rankMatchUser.deletePlayerID(userID, player.Attr.Channel)
		if len(r.rankMatchUser.GetPlayerIDList(player.Attr.Channel)) == 0 {
			r.ClearRankMatch(player.Attr.Channel)
		}
		logger.Log.Info("quitRankMatch:userID:", userID)

	case "mergeRankReq":
		merData := req.MsgData.(*mergeRankMatchReq)
		// sort.Sort(r.userRankList)
		list := r.rankMatchUser.GetPlayerIDListByStarLevel(merData.RankStarLevel, merData.NeedNum, merData.IsNext, merData.Channel)
		for i := range list {
			// 清理当前匹配器
			r.rankMatchUser.deletePlayerID(list[i], merData.Channel)
			if len(r.rankMatchUser.GetPlayerIDList(merData.Channel)) == 0 {
				r.ClearRankMatch(merData.Channel)
			}
		}
		reply := &mergeRankMatchReply{PlayerIdList: list, Channel: merData.Channel}
		merData.RankMatch.rankMatchPid.Cast("mergeRankResp", reply)
		logger.Log.Infof("mergeRankReq:%v reply:%v", merData.RankMatch.rank, reply)

	case "mergeRankResp":
		reply := req.MsgData.(*mergeRankMatchReply)
		for i := range reply.PlayerIdList {
			if reply.PlayerIdList[i] <= 0 {
				continue
			}
			player := global.GloInstance.GetPlayer(reply.PlayerIdList[i])
			logger.Log.Infof("mergeRankResp rank:%d userid:%d", r.rank, player.(*Player).Attr.UserID)
			r.rankMatchPid.Cast("enterRankMatch", player)
		}
		// 如果没找到需要合并的人员，再次发送
		if reply.Channel != "" {
			l := len(r.rankMatchUser.GetPlayerIDList(reply.Channel))
			if l < int(r.matchNum) && l > 0 {
				r.sendMergeRankTimer(reply.Channel)
			}
		}
	default:
		fmt.Println("err matchMgr handle call method")
	}
}

func (r *rankMatch) HandleInfo(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("rank matchMgr SOCKET_EVENT HandleInfo error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			debug.PrintStack()
			logger.Log.Errorf("rankMatch HandleInfo panic-------------rankMatch:%v", r)
		}
	}()
	switch req.Method {
	case "mergeRankTimer":
		ch := req.MsgData.(string)
		logger.Log.Infof("start mergeRank rankMatch:%v", r.rank)
		l := len(r.rankMatchUser.GetPlayerIDList(ch))
		if l < int(r.matchNum) {
			chanMat := r.rankMatchUser.userMap[ch]
			if chanMat.mergeNum > 6 { // 30秒后增加机器人，5s进行一次合并
				// add robot
				if chanMat.mergeNum < 9 {
					idList, b := r.rankMatchUser.GetPlayerIDListByChannel(ch, r.matchNum)
					logger.Log.Infof("GetPlayerIDListByChannel players:%v matchNum:%d", chanMat, chanMat.mergeNum)
					if b {
						reply := &mergeRankMatchReply{PlayerIdList: idList, Channel: ch}
						r.rankMatchPid.Cast("mergeRankResp", reply)
					} else {
						r.sendMergeRankTimer(ch)
					}
					return
				}
				if chanMat.mergeNum > 12 {
					// 当前机器人加玩，检查其他channel
					c := r.rankMatchUser.GetOtherChannel(ch)
					r.sendMergeRankTimer(c)
					return
				}
				aiId := r.rankMatchUser.GetRobotNum(ch) + 1
				p := NewAIPlayer(aiId, ch)
				r.rankMatchPid.Cast("enterRankMatch", p)
				logger.Log.Infof("mergeRankTimer robot current players:%v, and robot matchNum:%d", chanMat, chanMat.mergeNum)
				return
			}
			// 以自身为定点先向下一段位请求匹配，如果没有则向上一段位请求匹配
			// 上下都请求完一次，则在夸大上下匹配星级，请求匹配
			if chanMat.mergeNum%2 == 0 {
				// 向高级段位的请求合并
				if r.nextRankMatch != nil {
					r.nextMergeRank(ch)
				} else {
					r.preMergeRank(ch)
				}
			} else {
				if r.preRankMatch != nil {
					r.preMergeRank(ch)
				} else {
					r.nextMergeRank(ch)
				}
			}
		}
	case "rankAddAIALL":
		ch := req.MsgData.(string)
		num := r.matchNum - int32(len(r.rankMatchUser.GetPlayerIDList(ch)))
		if num == r.matchNum {
			r.rankMatchUser.ClearRankList(ch)
			return
		}
		startNum := r.rankMatchUser.GetRobotNum(ch) + 1
		playerList := GetPlayerAiList(startNum, startNum+num, ch)
		for i := 0; i < len(playerList); i++ {
			logger.Log.Infof("rankAddAIALL rank:%d userid:%d", r.rank, playerList[i].Attr.UserID)
			r.rankMatchPid.Cast("enterRankMatch", playerList[i])
		}
		logger.Log.Info("rank add ai num:", num)
	default:
		fmt.Println("err rank matchMgr handle call method")
	}
}

func (r *rankMatch) Terminate() {
	fmt.Println("rank match pid terminate ")
}

func (r *rankMatch) sendMergeRankTimer(ch string) {
	chanMat, ok := r.rankMatchUser.userMap[ch]
	if !ok {
		logger.Log.Errorf("this ch isn't exist:%s", ch)
		return
	}
	chanMat.mergeNum++
	r.rankMatchUser.userMap[ch] = chanMat
	if chanMat.mergeNum > 6 { // 30秒后增加机器人，5s进行一次合并
		if r.mergeAiTimeout <= 0 { // 星耀以上不增加机器人
			return
		}
		// if chanMat.mergeNum > 12 { // 机器人最多增加四个 + 2次渠道匹配
		// 	return
		// }
		r.rankMatchPid.SendAfter("mergeRankTimer", r.mergeTimerKey, r.mergeAiTimeout, ch)
		return
	}
	// 合并匹配不同段位人员
	r.rankMatchPid.SendAfter("mergeRankTimer", r.mergeTimerKey, r.mergeTimeout, ch)
}

func (r *rankMatch) preMergeRank(ch string) {
	// 向上请求匹配，如果上下匹配器都有，则r.mergeNum%2 == 0
	// 向上请求匹配，如果下匹配器都有，则r.mergeNum%2 == 0 或 1 r.mergeNum = 1、2、3... 计算第几次向上匹配时 需要对奇数+1
	chanMat := r.rankMatchUser.userMap[ch]
	var starLevel int32
	if chanMat.mergeNum%2 == 0 {
		starLevel = chanMat.minStarLevel - chanMat.mergeNum/2*5
	} else {
		starLevel = chanMat.minStarLevel - (chanMat.mergeNum+1)/2*5
	}

	if r.rank == proto3.PlayerRankEnum_king_rank { // 青铜一直往下合并
		starLevel = chanMat.minStarLevel - chanMat.mergeNum*5
	}
	needNum := r.matchNum - int32(len(chanMat.GetPlayerIDList()))
	req := &mergeRankMatchReq{
		RankMatch:     r,
		NeedNum:       needNum,
		RankStarLevel: starLevel, // 偶数不加一即可得到下一星级
		IsNext:        false,
		Channel:       ch,
	}
	r.preRankMatch.rankMatchPid.Cast("mergeRankReq", req)

	logger.Log.Infof("mergeRankTimer pre match current players:%v starCount:%d needNum:%d mergeNum:%v", chanMat, starLevel, needNum, chanMat.mergeNum)
}

func (r *rankMatch) nextMergeRank(ch string) {
	// 向下请求匹配，如果上下匹配器都有，则r.mergeNum%2 == 1
	// 向下请求匹配，如果上匹配器没有，则r.mergeNum%2 == 0 或 1 r.mergeNum = 1、2、3... 计算第几次向下匹配时 需要对偶数+2
	chanMat := r.rankMatchUser.userMap[ch]
	var starLevel int32
	if chanMat.mergeNum%2 == 0 {
		starLevel = chanMat.maxStarLevel + (chanMat.mergeNum+2)/2*5
	} else {
		starLevel = chanMat.maxStarLevel + (chanMat.mergeNum+1)/2*5
	}
	if r.rank == proto3.PlayerRankEnum_bronze_rank { // 青铜一直往下合并
		starLevel = chanMat.maxStarLevel + chanMat.mergeNum*5
	}
	needNum := r.matchNum - int32(len(chanMat.GetPlayerIDList()))
	req := &mergeRankMatchReq{
		RankMatch:     r,
		NeedNum:       needNum,
		RankStarLevel: starLevel, // 奇数加一即可得到上一星级
		IsNext:        true,
		Channel:       ch,
	}
	r.nextRankMatch.rankMatchPid.Cast("mergeRankReq", req)

	logger.Log.Infof("mergeRankTimer next match current players:%v starCount:%d needNum:%d mergeNum:%d", chanMat, starLevel, needNum, chanMat.mergeNum)
}
