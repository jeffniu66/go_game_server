package tableconfig

import "go_game_server/server/util"

type QuaConfig struct {
	ID          int32 `json:"id"`
	StartTime   int32
	EndTime     int32
	NeedLevel1  int32  `json:"needlevel1"`
	Rewards1    string `json:"rewards1"`
	NeedLevel2  int32  `json:"needlevel2"`
	Rewards2    string `json:"rewards2"`
	StartimeStr string `json:"startime"`
	EndtimeStr  string `json:"endtime"`
}

type QuaConfigCol struct {
	QuaConfigList []QuaConfig
	QuaConfigMap  map[int32]*QuaConfig
}

// 判断是否在当前赛季 返回param1 旧赛季奖励，param bool-true
func (q *QuaConfigCol) IsNowQua(t int32) (*QuaConfig, bool) {
	now := util.UnixTime()
	var oldQua, nowQua *QuaConfig
	var b bool = false
	for _, v := range q.QuaConfigList {
		if util.IsTwoTimeHasTime(v.StartTime, v.EndTime, t) {
			oldQua = &v
		}
		if util.IsTwoTimeHasTime(v.StartTime, v.EndTime, now) {
			nowQua = &v
		}
		if oldQua != nil && nowQua != nil {
			break
		}
	}
	if oldQua != nil && nowQua != nil && oldQua.ID != nowQua.ID {
		b = true
	}
	return oldQua, b
}

func (q *QuaConfigCol) GetQuaConfig(ID int32) *QuaConfig {
	v, ok := q.QuaConfigMap[ID]
	if !ok {
		return q.QuaConfigMap[ID]
	}
	return v
}

func (q *QuaConfigCol) InitMap() {
	q.QuaConfigMap = make(map[int32]*QuaConfig, 0)
	for i := 0; i < len(q.QuaConfigList); i++ {
		v := &q.QuaConfigList[i]
		v.StartTime = util.ToTimeStamp(v.StartimeStr + " 00:00:00")
		v.EndTime = util.ToTimeStamp(v.EndtimeStr + " 23:59:59")
		q.QuaConfigMap[v.ID] = v
	}
}

var QuaConfigs *QuaConfigCol
