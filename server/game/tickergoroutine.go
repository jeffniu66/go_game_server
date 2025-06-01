package game

import (
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"runtime/debug"
)

type TickerTask struct {
	tickerPid         *global.PidObj
	worldChatTextPool *WorldChatTextPool
}

func CreateTickerTask() *TickerTask {
	tickerObj := new(TickerTask)
	tickerObj.tickerPid = global.RegisterPid("tickerTaskManange", 2048, tickerObj)
	tickerObj.worldChatTextPool = InitWorldChat()
	return tickerObj
}

func (t *TickerTask) Start() {
	logger.Log.Info("tickerTaskManage start ...")
}

func (t *TickerTask) HandleCall(req global.GenReq) global.Reply {
	return nil
}

func (t *TickerTask) HandleCast(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("room SOCKET_EVENT HandleCast error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			logger.Log.Errorf("room HandleCast panic-------------TickerTask:%v", t)
			debug.PrintStack()
		}
	}()
	switch req.Method {
	default:
		logger.Log.Info("err TickerTask HandleCast call method")
	}
}

func (t *TickerTask) HandleInfo(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("TickerTask SOCKET_EVENT HandleCast error:%v, req: %v, stack: %v ", err, req, string(debug.Stack()))
			logger.Log.Errorf("TickerTask HandleInfo panic-------------TickerTask:%v", t)
			debug.PrintStack()
		}
	}()

	switch req.Method {
	case constant.TickerZeroTimer:
		t.zeroTimer()
	case constant.TickerWorldChatTime:
		t.worldChatOnce()
	case constant.TickerWorldTextCD:
		index, ok := req.MsgData.(int32)
		if !ok {
			logger.Log.Errorf("TickerTask TickerWorldTextCD error, reqdata isn't int32, req:%v", req)
			return
		}
		t.worldChatTextPool.TextCDEnd(index)
	default:
		logger.Log.Info("err TickerTask HandleInfo call method")
	}
}

func (t *TickerTask) Terminate() {
	logger.Log.Errorf("TickerTask pid terminate roomPidName:%s", t.tickerPid.PidName)
}
