package game

import (
	"go_game_server/server/constant"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"strings"
)

type WorldChatRate struct {
	LowTime  int32
	HighTime int32
	Rate     int32
}

type WorldChatTextPool struct {
	PreWorldText   *tableconfig.QuickTextConfig
	ChatRate       []*WorldChatRate // 0,1=20|1,9=60|9,12=20|12,14=15|14,19=20|19,23=10|23,24=20 [0,1)
	NormalTextPool []*tableconfig.QuickTextConfig
	CDTextPool     []*tableconfig.QuickTextConfig
	CDTextCache    map[int32]*tableconfig.QuickTextConfig // key-id
}

func InitWorldChat() *WorldChatTextPool {
	w := new(WorldChatTextPool)
	w.NormalTextPool = make([]*tableconfig.QuickTextConfig, 0)
	w.CDTextPool = make([]*tableconfig.QuickTextConfig, 0)
	w.CDTextCache = make(map[int32]*tableconfig.QuickTextConfig, 0)
	if len(tableconfig.QuickTextCols.QuickTypeMap[constant.TextTypeWorld]) <= 0 {
		logger.Log.Error("quickText TextTypeWorld is zero")
	}
	for _, v := range tableconfig.QuickTextCols.QuickTypeMap[constant.TextTypeWorld] {
		if v.AiDataType == constant.TextNormal {
			w.NormalTextPool = append(w.NormalTextPool, v)
		} else if v.AiDataType == constant.TextCD {
			w.CDTextPool = append(w.CDTextPool, v)
		} else {
			continue
		}
	}
	w.initChatRate()
	return w
}

func (w *WorldChatTextPool) initChatRate() {
	chatRate := tableconfig.ConstsConfigs.GetValueById(constant.WorldChatTextRate)
	rateList := strings.Split(chatRate, "|")
	for _, v := range rateList {
		row := new(WorldChatRate)
		vList := strings.Split(v, "=")
		if len(vList) != 2 {
			continue
		}
		row.Rate = util.ToInt(vList[1])
		hList := strings.Split(vList[0], ",")
		if len(hList) != 2 {
			continue
		}
		row.LowTime = util.ToInt(hList[0])
		row.HighTime = util.ToInt(hList[1])
		w.ChatRate = append(w.ChatRate, row)
	}
}

func (w *WorldChatTextPool) getChatRate(h int32) *WorldChatRate {
	for _, v := range w.ChatRate {
		if h >= v.LowTime && h < v.HighTime {
			return v
		}
	}
	return nil
}

// ret timeout second
func (w *WorldChatTextPool) GetNextChatTime() int32 {
	var timeout int32 = 0
	h, m := util.GetHour(), util.GetMinute()
	chatRate := w.getChatRate(h)
	if chatRate == nil {
		logger.Log.Errorf("this:%d hour don't set rate", h)
		return -1
	}
	// [a, n1, 2a) --> [2a, n2, 3a) timeout = a -n1moda + r
	r := util.RandInt(0, chatRate.Rate-1)
	timeout = chatRate.Rate - m%chatRate.Rate + r
	return timeout
}

func (w *WorldChatTextPool) GetRandText(t *TickerTask) *tableconfig.QuickTextConfig {
	var textConf *tableconfig.QuickTextConfig
	var r int32
	nl, cl := int32(len(w.NormalTextPool)), int32(len(w.CDTextPool))
	if w.PreWorldText != nil && w.PreWorldText.AiDataType == constant.TextCD {
		r = util.RandInt(0, nl-1)
	}
	if w.PreWorldText == nil || w.PreWorldText.AiDataType == constant.TextNormal {
		r = util.RandInt(0, nl+cl-1)
	}

	if r < nl {
		textConf = w.NormalTextPool[r]
	} else {
		index := r - nl
		textConf = w.CDTextPool[index]
		cd := tableconfig.ConstsConfigs.GetIdValue(constant.WorldChatTextCD) // s
		textConf.CDEndTime = util.UnixTime() + cd
		w.CDTextPool = append(w.CDTextPool[0:index], w.CDTextPool[index+1:]...)
		w.CDTextCache[textConf.Id] = textConf
		t.tickerPid.SendAfter(constant.TickerWorldTextCD, constant.TickerWorldTextCD+util.ToStr(textConf.Id), cd*1000, textConf.Id)
	}
	w.PreWorldText = textConf
	return textConf
}

func (w *WorldChatTextPool) TextCDEnd(id int32) {
	textConf := w.CDTextCache[id]
	w.CDTextPool = append(w.CDTextPool, textConf)
	delete(w.CDTextCache, id)
}
