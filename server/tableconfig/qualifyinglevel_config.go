package tableconfig

import (
	"go_game_server/server/constant"
	"go_game_server/server/util"
	"strings"
)

type QualifyingLevelConfig struct {
	ID            int32  `json:"id"`
	Prev          int32  `json:"prev"`
	Next          int32  `json:"next"`
	Rank          int32  `json:"rank"`
	Stage         int32  `json:"stage"`
	MaxStage      int32  `json:"maxstage"`
	MaxStar       int32  `json:"star"`
	ChangeLevel   int32  `json:"changelevel"`
	ChangeStar    int32  `json:"changestar"`
	Points        int32  `json:"points"`
	SavePoints    int32  `json:"savepoints"`
	AddTime       int32  `json:"addtime"`
	AiAllTime     int32  `json:"aijointime"`
	Aitime        int32  `json:"aitime"`
	Winpoints     string `json:"winpoints"`
	Killpoints    string `json:"killpoints"`
	MissionPoints string `json:"missionpoints"`
	VotePoints    string `json:"votepoints"`
	Rewards       string `json:"rewards"`
	WinRewards    string `json:"winrewards"`
	TaskFinReward int32  `json:"taskfinishrewards"`
	AllFinReward  int32  `json:"alltaskfinishrewards"`
	KillReward    int32  `json:"killrewards"`
	WinGoldReward int32
	WinExpReward  int32
}

type QualityingConfigCol struct {
	StartID         int32
	EndID           int32
	QuaLevelCofList []QualifyingLevelConfig
	QuaLevelCofMap  map[int32]*QualifyingLevelConfig
}

func (q *QualityingConfigCol) GetQuaConfig(rankID int32) *QualifyingLevelConfig {
	v, ok := q.QuaLevelCofMap[rankID]
	if !ok {
		return q.QuaLevelCofMap[q.StartID]
	}
	return v
}

func (q *QualityingConfigCol) InitMap() {
	q.QuaLevelCofMap = make(map[int32]*QualifyingLevelConfig, 0)
	for i := 0; i < len(q.QuaLevelCofList); i++ {
		v := &q.QuaLevelCofList[i]
		if v.ID == v.Prev {
			q.StartID = v.ID
		}
		if v.ID == v.Next {
			q.EndID = v.ID
		}
		winReList := strings.Split(v.WinRewards, "|")
		for _, row := range winReList {
			r := strings.Split(row, ",")
			if len(r) != 2 {
				continue
			}
			if r[0] == util.ToStr(constant.ItemIdGold) {
				q.QuaLevelCofList[i].WinGoldReward = util.ToInt(r[1])
			}
			if r[0] == util.ToStr(constant.ItemIdExp) {
				q.QuaLevelCofList[i].WinExpReward = util.ToInt(r[1])
			}
		}

		q.QuaLevelCofMap[v.ID] = v
	}
}
func (q *QualityingConfigCol) GetMergeTimeout(rank int32) (int32, int32) {
	defTime := int32(3)
	defTime2 := int32(8)
	for i := range q.QuaLevelCofList {
		if q.QuaLevelCofList[i].Rank == rank {
			defTime = q.QuaLevelCofList[i].AddTime
			defTime2 = q.QuaLevelCofList[i].Aitime
			break
		}
	}
	return defTime, defTime2
}

func (q *QualityingConfigCol) GetGoldAndExp(rankID int32) (int32, int32) {
	gold := int32(50) // 金币
	exp := int32(50)  // 经验
	v, ok := q.QuaLevelCofMap[rankID]
	if !ok {
		return gold, exp
	}

	strList := strings.Split(v.Rewards, "|")
	if len(strList) != 2 {
		return gold, exp
	}
	strData := strings.Split(strList[0], ",")
	if len(strData) != 2 {
		return gold, exp
	}
	gold |= util.ToInt(strData[1]) // 默认50

	strData = strings.Split(strList[1], ",")
	if len(strData) != 2 {
		return gold, exp
	}
	exp |= util.ToInt(strData[1])
	return gold, exp
}

func (q *QualityingConfigCol) GetChangeRank(rankID, star int32) (retRankID, retStar int32) {
	retRankID, retStar = rankID, star

	rankInfo, ok := q.QuaLevelCofMap[rankID]
	if !ok {
		return
	}
	if star < 0 {
		// star = -1 降阶
		retStar = 3
		retRankID = rankInfo.Prev
		if rankInfo.ID == rankInfo.Prev {
			retStar = 0
		}
	}
	if rankInfo.MaxStar < star {
		retStar = 0
		retRankID = rankInfo.Next
	}
	return
}

func (q *QualityingConfigCol) GetPreRank(rankID, star int32) (retRankID, retStar int32) {
	retRankID, retStar = rankID, star
	rankInfo, ok := q.QuaLevelCofMap[rankID]
	if !ok {
		return
	}

	if 0 > star {
		preRank := q.QuaLevelCofMap[rankInfo.Prev]
		retStar = preRank.MaxStar
		retRankID = preRank.ID
	}
	return
}

var QuaLevelConfs *QualityingConfigCol
