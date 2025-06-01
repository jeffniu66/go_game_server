package global

import (
	"go_game_server/server/logger"
	"runtime/debug"
	"sync"
	"time"
)

var PidObjMap *sync.Map

type PidObj struct {
	Callback     GenServer
	ReqCh        chan GenReq
	ReplyCh      chan Reply
	StopTimer    bool
	TimerMap     *sync.Map
	PidName      string
	MessageBox   []GenReq
	ReqChClose   bool
	ReplyChClose bool
}

type GenReq struct {
	Method  string
	MsgData interface{}
	t       int
	time    int32
}

type Reply interface{}

const (
	call = iota
	cast
	timer
	shutdown
)

// Start server
func RegisterPid(pidName string, chanLen int, callback GenServer) *PidObj {
	pidObj := &PidObj{Callback: callback}
	pidObj.ReqCh = make(chan GenReq, chanLen)
	pidObj.ReplyCh = make(chan Reply, 5)
	pidObj.TimerMap = &sync.Map{}
	pidObj.PidName = pidName
	callback.Start()
	if PidObjMap == nil {
		PidObjMap = &sync.Map{}
	}
	PidObjMap.Store(pidName, pidObj) // 用于存储协程对应的通道里面的消息,方便查找阻塞内容
	go pidObj.loop()
	return pidObj
}

func (p *PidObj) loop() {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("genserver loop err: ", err, string(debug.Stack()))
		}
	}()
	for {
		req := <-p.ReqCh
		p.MessageBox = append(p.MessageBox, req)
		PidObjMap.Store(p.PidName, p)
		switch req.t {
		case call:
			reply := p.Callback.HandleCall(req)
			p.ReplyCh <- reply
		case cast:
			p.Callback.HandleCast(req)
		case timer:
			p.Callback.HandleInfo(req)
		case shutdown:
			p.StopAllTimer()
			p.StopTimer = true
			p.Callback.Terminate()
			close(p.ReqCh)
			p.ReqChClose = true
			close(p.ReplyCh) // 相当于往通道发空信息p.ReplyCh <- nil
			p.ReplyChClose = true
			return
		default:
			logger.Log.Errorln("------------- pid loop error : ", req.t, len(p.ReqCh))
		}
		p.MessageBox = p.MessageBox[1:]
		PidObjMap.Store(p.PidName, p)
	}
}

func (p *PidObj) Call(method string, msg interface{}) (reply Reply) {
	if p.ReqChClose || p.ReplyChClose {
		logger.Log.Warnln("replych channel is closed: ", p.PidName, method, msg)
		return
	}
	if len(p.ReqCh) >= 200 {
		logger.Log.Warnln("------------------------- genserver call: ", p.PidName, len(p.ReqCh))
	}
	p.ReqCh <- GenReq{Method: method, MsgData: msg, t: call}
	reply = <-p.ReplyCh
	return
}

func (p *PidObj) Cast(method string, msg interface{}) {
	if p.ReqChClose || p.ReplyChClose {
		logger.Log.Warnln("req channel is closed: ", p.PidName, method, msg)
		return
	}
	if len(p.ReqCh) >= 200 {
		logger.Log.Warnf("------------------------- genserver pid:%v cast: %v len:%v", p.PidName, method, len(p.ReqCh))
	}
	p.ReqCh <- GenReq{Method: method, MsgData: msg, t: cast}
}

func (p *PidObj) SendAfter(method, timerKey string, millisecond int32, msg interface{}) {
	if _, ok := p.TimerMap.Load(timerKey); !ok { // 如果map里面没有该定时器，即ok为false，就要起定时器
		timer := time.AfterFunc(time.Duration(millisecond)*time.Millisecond, func() {
			if !p.StopTimer {
				p.TimerMap.Delete(timerKey)
				if len(p.ReqCh) >= 200 {
					logger.Log.Warnln("------------------------- genserver sendafter: ", p.PidName, len(p.ReqCh))
				}
				p.ReqCh <- GenReq{Method: method, MsgData: msg, t: timer, time: millisecond} // 通过读chan时的第二个返回值可以知道，但仅适用于接收端； 通过recover的方式可以勉强实现，但是不推荐；
			}
		})
		p.TimerMap.Store(timerKey, timer)
	}
}

// 禁止自己协程里面操作Stop, 即自己Stop自己，会造成阻塞
func (p *PidObj) Stop() (reply Reply) {
	if p.ReqChClose || p.ReplyChClose {
		logger.Log.Warnln("channel is closed: ", p.PidName)
		return
	}
	if len(p.ReqCh) >= 200 {
		logger.Log.Warnln("------------------------- genserver Stop: ", p.PidName, len(p.ReqCh))
	}
	p.ReqCh <- GenReq{t: shutdown}
	reply = <-p.ReplyCh // 同步返回，是否有可能阻塞
	return
}

// 可以在自己协程里面操作Stop, 不必同步等待返回
func (p *PidObj) CastStop() {
	if len(p.ReqCh) >= 200 {
		logger.Log.Warnln("------------------------- genserver CastStop: ", p.PidName, len(p.ReqCh))
	}
	p.ReqCh <- GenReq{t: shutdown}
	return
}

func (p *PidObj) GetPidTimer(method string) *time.Timer {
	if timer, ok := p.TimerMap.Load(method); ok {
		return timer.(*time.Timer)
	}
	return nil
}

func (p *PidObj) StopPidTimer(timerKey string) {
	if timer, ok := p.TimerMap.Load(timerKey); ok {
		timer.(*time.Timer).Stop()
		p.TimerMap.Delete(timerKey)
	}
}

func (p *PidObj) StopAllTimer() {
	p.TimerMap.Range(func(key, timer interface{}) bool {
		timer.(*time.Timer).Stop()
		return true
	})
	p.TimerMap = &sync.Map{}
}
