package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
	"time"
)

var tickerTaskManage *TickerTask

func FreshTickerManageData() {
	tickerTaskManage.worldChatTextPool = InitWorldChat()
}

func NewGlobalTimer() {
	tickerTaskManage = CreateTickerTask()

	// newTimer(timerFunc)
	// tickerTaskManage.zeroTimer()
	// tickerTaskManage.worldChatOnce()

	tickerTaskManage.StartTask()
}

var count int

func timerFunc() {
	//logger.Log.Infof(">>>>>>>>>>>> 定时器启动 每十分钟一次ticker \n\n")
	count++
	hour, minute := util.GetHour(), util.GetMinute()
	if hour == 23 && minute >= 50 {
		// time.AfterFunc(time.Duration(util.GetNextIntPoint())*time.Second, zeroTimer)
	}
}

func newTimer(timer func()) {
	timer() // 开服先执行一次
	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			timer()
		}
	}
}

func (t *TickerTask) StartTask() {
	refreshTime := (util.Get24TimeStamp(util.UnixTime()) - util.UnixTime() + 1) * 1000
	t.tickerPid.SendAfter(constant.TickerZeroTimer, constant.TickerZeroTimer, refreshTime, 0)

	timeout := t.worldChatTextPool.GetNextChatTime() * 1000
	t.tickerPid.SendAfter(constant.TickerWorldChatTime, constant.TickerWorldChatTime, timeout, nil)
}

func (t *TickerTask) zeroTimer() {
	playerList := global.GloInstance.GetPlayerIDList()
	logger.Log.Infoln("每天晚上凌晨0点左右更新, 玩家id列表为：", playerList)
	statDan()
	refreshTime := (util.Get24TimeStamp(util.UnixTime()) - util.UnixTime() + 1) * 1000
	t.tickerPid.SendAfter(constant.TickerZeroTimer, constant.TickerZeroTimer, refreshTime, 0)
}

func (t *TickerTask) worldChatOnce() {
	text := t.worldChatTextPool.GetRandText(t)
	if text == nil {
		logger.Log.Errorf("worldChatOnce is nil")
		return
	}
	player := NewAIPlayer(10000, "dev")
	SpreadWorldPlayer(player, proto3.ChatTypeEnum_chat_merge, "", constant.QuickTextType, text.Id, nil)
	timeout := t.worldChatTextPool.GetNextChatTime() * 1000
	t.tickerPid.SendAfter(constant.TickerWorldChatTime, constant.TickerWorldChatTime, timeout, nil)
}

// 定时统计
func ScheduleStat() {
	ticker := time.NewTicker(10 * time.Minute)
	for {
		select {
		case <-ticker.C:
			scheduleFunc()
		}
	}
}

func scheduleFunc() {
	hour, minute := util.GetHour(), util.GetMinute()
	if hour == 23 && minute >= 50 {
		time.AfterFunc(time.Duration(util.GetNextIntPoint())*time.Second, statDan)
	}
}

// 统计新注册玩家段位分布
func statDan() {
	logger.Log.Infoln("==================统计新注册玩家段位分布开始============================")
	curTime := util.UnixTime()
	startTime := curTime - constant.DaySecond
	dans := db.GetDan(startTime, curTime)
	if len(dans) == 0 {
		return
	}
	date := util.GetDateStr(time.Now().Unix() - constant.DaySecond)
	for _, dan := range dans {
		dan.RegisterDate = date
		db.InsertDanStat(dan)
	}
	logger.Log.Infoln("==================统计新注册玩家段位分布结束============================")
}
