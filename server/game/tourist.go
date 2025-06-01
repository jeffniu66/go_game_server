package game

import (
	"go_game_server/server/db"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

func CreateTourist(writeChan chan interface{}) *global.PidObj {
	tourist := &Player{WriteChan: writeChan}
	tourist.Attr = &db.Attr{}
	tourist.Attr.UserID = 10
	tourist.Attr.Username = "tourist"
	pidName := "tourist_" + "10"
	pid := global.RegisterPid(pidName, 256, tourist)
	// 1分钟后退出
	lengthTime := 60 * 1000
	pid.SendAfter("player_stop", "player_stop", int32(lengthTime), nil)
	tourist.Pid = pid
	return tourist.Pid
}

func (p *Player) Register(acctName, pw string) (*db.Acct, bool) {
	acct := &db.Acct{}
	acct.AcctName = acctName
	pwS, err := util.PasswordEncode(pw, "", 0)
	if err != nil {
		logger.Log.Errorf("register failed,err: %v", err)
		return nil, false
	}
	acct.Password = pwS
	db.SaveAcct(acct)
	return acct, true
}
