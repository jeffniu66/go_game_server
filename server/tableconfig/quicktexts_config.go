package tableconfig

import (
	"go_game_server/server/constant"
	"go_game_server/server/util"
	"strings"
)

type QuickTextConfig struct {
	Id          int32  `json:"id"`
	Gender      int32  `json:"gender"`
	QuickType   int32  `json:"type"`
	AiDataType  int32  `json:"aidata_type"`
	DataNumber  int32  `json:"data_number"`
	AiChatLimit string `json:"robottrigger"` // [low, high]
	QuickText   string `json:"quicktext"`
	VoiceKey    string `json:"vo"`
	LimitLow    int32  // limit 低位
	LimitHigh   int32  // limit 高位
	CDEndTime   int32  // cd结束时间
}

type QuickTextCol struct {
	QuickTextConfigList []QuickTextConfig
	QuickIDMap          map[int32]*QuickTextConfig
	QuickTypeMap        map[int32][]*QuickTextConfig // key-type
}

func (q *QuickTextCol) InitMap() {
	q.QuickIDMap = make(map[int32]*QuickTextConfig, 0)
	q.QuickTypeMap = make(map[int32][]*QuickTextConfig, 0)
	for i := range q.QuickTextConfigList {
		v := &q.QuickTextConfigList[i]
		if v.AiChatLimit != "" {
			l := strings.Split(v.AiChatLimit, ",")
			if len(l) == 2 {
				v.LimitLow = util.ToInt(l[0])
				v.LimitHigh = util.ToInt(l[1])
			}
		}
		q.QuickIDMap[v.Id] = v
		if mv, ok := q.QuickTypeMap[v.QuickType]; ok {
			mv = append(mv, v)
			q.QuickTypeMap[v.QuickType] = mv
		} else {
			tmp := []*QuickTextConfig{v}
			q.QuickTypeMap[v.QuickType] = tmp
		}
	}
}

func (q *QuickTextCol) GetRandText(textType, liveNum int32) *QuickTextConfig {
	var ret *QuickTextConfig
	mp, ok := q.QuickTypeMap[textType]
	if !ok {
		return nil
	}
	n := 0
	for {
		n++
		if n > 100 { // 防止死循环，找第一个
			for _, v := range mp {
				if liveNum >= v.LimitLow && liveNum <= v.LimitHigh {
					ret = v
					break
				}
			}
			break
		}
		r := util.RandInt(0, int32(len(mp))-1)
		ret = mp[r]
		if textType == constant.TextTypeWait {
			break // 候场聊天直接返回
		}
		if liveNum >= ret.LimitLow && liveNum <= ret.LimitHigh {
			break
		}
	}
	return ret
}

/*
函数合并至 GetRandText order by textType
func (q *QuickTextCol) GetRandWaitText() int32 {
	mp, ok := q.QuickTypeMap[constant.TextTypeWait]
	if !ok {
		return -1
	}
	r := util.RandInt(0, int32(len(mp))-1)
	return mp[r].Id
}

func (q *QuickTextCol) GetRandVoteText() int32 {
	mp, ok := q.QuickTypeMap[constant.TextTypeVote]
	if !ok {
		return -1
	}
	r := util.RandInt(0, int32(len(mp))-1)
	return mp[r].Id
}
*/

var QuickTextCols *QuickTextCol
