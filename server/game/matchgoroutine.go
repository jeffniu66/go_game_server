package game

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"runtime/debug"
)

var MatchMgrPid *global.PidObj

// 该匹配器是单线程运行的
type matchMgr struct {
	playerIdMap  map[int32]int                 // 单线程下并发安全 key userId 用于过滤同号玩家
	matchNum     int32                         // 策划配置
	rankMatchMap map[PlayerRankType]*rankMatch // 段位匹配 0-青铜，6-王者
}

func InitMatchMgr() {
	playerIdMap := make(map[int32]int, 0)
	langNum, pingNum := tableconfig.ConstsConfigs.GetMatchGameParam()
	matchNum := int32(langNum + pingNum)
	matchMgrObj := &matchMgr{playerIdMap, matchNum, nil}

	// 初始化各个段位map
	rankMatchMap := make(map[PlayerRankType]*rankMatch)
	rankMatchMap[proto3.PlayerRankEnum_bronze_rank] = initRankMatch(proto3.PlayerRankEnum_bronze_rank, constant.MinBronzeRankStar, constant.MaxBronzeRankStar)
	rankMatchMap[proto3.PlayerRankEnum_silver_rank] = initRankMatch(proto3.PlayerRankEnum_silver_rank, constant.MinSliverRankStar, constant.MaxSliverRankStar)
	rankMatchMap[proto3.PlayerRankEnum_gold_rank] = initRankMatch(proto3.PlayerRankEnum_gold_rank, constant.MinGoldRankStar, constant.MaxGoldRankStar)
	rankMatchMap[proto3.PlayerRankEnum_platinum_rank] = initRankMatch(proto3.PlayerRankEnum_platinum_rank, constant.MinPlatinumRankStar, constant.MaxPlatinumRankStar)
	rankMatchMap[proto3.PlayerRankEnum_diamond_rank] = initRankMatch(proto3.PlayerRankEnum_diamond_rank, constant.MinDiamondRankStar, constant.MaxDiamondRankStar)
	rankMatchMap[proto3.PlayerRankEnum_starshine_rank] = initRankMatch(proto3.PlayerRankEnum_starshine_rank, constant.MinStarShineRankStar, constant.MaxStarShineRankStar)
	rankMatchMap[proto3.PlayerRankEnum_king_rank] = initRankMatch(proto3.PlayerRankEnum_king_rank, constant.MinKingRankStar, constant.InfRankStar)
	// 上下段位初始化，用于合并段位
	rankMatchMap[proto3.PlayerRankEnum_bronze_rank].preRankMatch = nil
	rankMatchMap[proto3.PlayerRankEnum_bronze_rank].nextRankMatch = rankMatchMap[proto3.PlayerRankEnum_silver_rank]
	rankMatchMap[proto3.PlayerRankEnum_silver_rank].preRankMatch = rankMatchMap[proto3.PlayerRankEnum_bronze_rank]
	rankMatchMap[proto3.PlayerRankEnum_silver_rank].nextRankMatch = rankMatchMap[proto3.PlayerRankEnum_gold_rank]
	rankMatchMap[proto3.PlayerRankEnum_gold_rank].preRankMatch = rankMatchMap[proto3.PlayerRankEnum_silver_rank]
	rankMatchMap[proto3.PlayerRankEnum_gold_rank].nextRankMatch = rankMatchMap[proto3.PlayerRankEnum_platinum_rank]
	rankMatchMap[proto3.PlayerRankEnum_platinum_rank].preRankMatch = rankMatchMap[proto3.PlayerRankEnum_gold_rank]
	rankMatchMap[proto3.PlayerRankEnum_platinum_rank].nextRankMatch = rankMatchMap[proto3.PlayerRankEnum_diamond_rank]
	rankMatchMap[proto3.PlayerRankEnum_diamond_rank].preRankMatch = rankMatchMap[proto3.PlayerRankEnum_platinum_rank]
	rankMatchMap[proto3.PlayerRankEnum_diamond_rank].nextRankMatch = rankMatchMap[proto3.PlayerRankEnum_starshine_rank]
	rankMatchMap[proto3.PlayerRankEnum_starshine_rank].preRankMatch = rankMatchMap[proto3.PlayerRankEnum_platinum_rank]
	rankMatchMap[proto3.PlayerRankEnum_starshine_rank].nextRankMatch = rankMatchMap[proto3.PlayerRankEnum_king_rank]
	rankMatchMap[proto3.PlayerRankEnum_king_rank].preRankMatch = rankMatchMap[proto3.PlayerRankEnum_starshine_rank]
	rankMatchMap[proto3.PlayerRankEnum_king_rank].nextRankMatch = nil

	matchMgrObj.rankMatchMap = rankMatchMap

	fmt.Println("init match mgr >>>>>> ", playerIdMap, matchNum)
	MatchMgrPid = global.RegisterPid("matchMgrPid", 2048, matchMgrObj)
}

func (m *matchMgr) Start() {
	fmt.Println("match mgr start ... ")
}

func (m *matchMgr) HandleCall(req global.GenReq) global.Reply {
	return nil
}

func (m *matchMgr) HandleCast(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("matchMgr SOCKET_EVENT HandleCast error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			debug.PrintStack()
			logger.Log.Errorf("match HandleCast panic-------------matchMgr:%v", m)
		}
	}()
	switch req.Method {
	case "onlyGame":
		player := req.MsgData.(*Player)
		rankInfo := tableconfig.QuaLevelConfs.GetQuaConfig(player.Attr.RankID)
		playerRank := PlayerRankType(rankInfo.Rank)
		if rankMatch, ok := m.rankMatchMap[playerRank]; ok {
			rankMatch.rankMatchPid.Cast("onlyGame", player)
		} else {
			logger.Log.Errorf("this rank match is nil rank:%v", playerRank)
		}
	case "enterMatch": // 进入匹配
		player := req.MsgData.(*Player)
		rankInfo := tableconfig.QuaLevelConfs.GetQuaConfig(player.Attr.RankID)
		playerRank := PlayerRankType(rankInfo.Rank)
		if rankMatch, ok := m.rankMatchMap[playerRank]; ok {
			rankMatch.rankMatchPid.Cast("enterRankMatch", player)
		} else {
			logger.Log.Errorf("this rank match is nil rank:%v", playerRank)
		}

	case "quitMatch":
		player := req.MsgData.(*Player)
		rankInfo := tableconfig.QuaLevelConfs.GetQuaConfig(player.Attr.RankID)
		playerRank := PlayerRankType(rankInfo.Rank)
		rankMatch := m.rankMatchMap[playerRank]
		rankMatch.rankMatchPid.Cast("quitRankMatch", player)

	default:
		fmt.Println("err matchMgr handle call method")
	}
}

func (m *matchMgr) HandleInfo(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("matchMgr SOCKET_EVENT HandleInfo error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			debug.PrintStack()
			logger.Log.Errorf("match HandleInfo panic-------------matchMgr:%v", m)
		}
	}()
	switch req.Method {
	default:
		fmt.Println("err matchMgr handle call method")
	}
}

func (m *matchMgr) Terminate() {
	fmt.Println("match pid terminate ")
}
