package db

import (
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"runtime/debug"
)

// 数据库协程 [类型][]*PidObj对象
var DBLineMap map[int][]*global.PidObj

const (
	UserDBType   = iota // 玩家数据库协程
	AllyDBType          // 同盟数据库协程
	CommonDBType        // 通用数据库协程
)

func InitDBPid() {
	userDbLine := global.MyConfig.ReadInt32("goroutine", "user_db_line")
	allyDbLine := global.MyConfig.ReadInt32("goroutine", "ally_db_line")
	commonDbLine := global.MyConfig.ReadInt32("goroutine", "common_db_line")

	DBLineMap = make(map[int][]*global.PidObj)
	var i int32
	// 玩家数据库协程
	userDbArray := make([]*global.PidObj, 0, userDbLine)
	for i = 0; i < userDbLine; i++ {
		dbUserPid := global.RegisterPid("DBUserPid", 2048, &DBPidObject{Index: i})
		userDbArray = append(userDbArray, dbUserPid)
	}
	DBLineMap[UserDBType] = userDbArray
	// 同盟数据库协程
	allyDbArray := make([]*global.PidObj, 0, allyDbLine)
	for i = 0; i < allyDbLine; i++ {
		dbAllyPid := global.RegisterPid("DBAllyPid", 2048, &DBPidObject{Index: i})
		allyDbArray = append(allyDbArray, dbAllyPid)
	}
	DBLineMap[AllyDBType] = allyDbArray
	// 通用数据库协程
	commonDbArray := make([]*global.PidObj, 0, commonDbLine)
	for i = 0; i < commonDbLine; i++ {
		dbCommonPid := global.RegisterPid("DBCommonPid", 1024, &DBPidObject{Index: i})
		commonDbArray = append(commonDbArray, dbCommonPid)
	}
	DBLineMap[CommonDBType] = commonDbArray
}

// 获取pidOject对象，按id取模
func getDBPid(t int, id int32) *global.PidObj {
	if id < 0 {
		return nil
	}

	if m, ok := DBLineMap[t]; ok {
		if len(m) == 0 {
			return nil
		}

		return m[int(id)%len(m)]
	}

	return nil
}

// 同步插入，防止停服时有些数据没执行完
func ExecDB(t int, id int32, query string, args ...interface{}) {
	_, err := DB.Exec(query, args...)
	if err != nil {
		logger.Log.Errorf("exec db err:%v, query = %s", err, query)
	}
}

// 异步插入，不需要等待数据库写完数据
func ExecDBAsync(t int, id int32, query string, args ...interface{}) {
	pid := getDBPid(t, id)
	if pid == nil {
		logger.Log.Warnf("get db pid is nil! type: %v, id: %v", t, id)
	} else {
		pid.Cast(Exec, &MsgData{query: query, args: args})
	}
}

/*****************************************************
	db pid object
********************************************************/
const (
	Exec = "Exec"
)

type MsgData struct {
	query string
	args  []interface{}
}

type DBPidObject struct {
	Index int32
}

func (d *DBPidObject) Start() {

}

func (d *DBPidObject) HandleCall(req global.GenReq) global.Reply {
	return nil
}

func (d *DBPidObject) HandleCast(req global.GenReq) {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Errorf("db goroutine error:%v\n, req: %v\n, msg: %v\n, stack: %v\n ",
				err, req, req.MsgData, string(debug.Stack()))
			debug.PrintStack()
		}
	}()

	data, ok := req.MsgData.(*MsgData)
	if !ok {
		return
	}

	switch req.Method {
	case Exec:
		_, err := DB.Exec(data.query, data.args...)
		if err != nil {
			logger.Log.Errorln("DBPidObject ----err", err, data.query)
		}
	}
}

func (d *DBPidObject) HandleInfo(global.GenReq) {

}

func (d *DBPidObject) Terminate() {

}
