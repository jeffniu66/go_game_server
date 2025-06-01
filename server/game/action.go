package game

import (
	"fmt"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"runtime/debug"
	"strings"
)

func saveAction(action *include.Action) {
	if action == nil {
		return
	}
	action.Update = include.Update
	db.ActionMap[action.UserId] = action
}

func getAction(userId int32) *include.Action {
	action, ok := db.ActionMap[userId]
	if !ok {
		action = db.SelectAction(userId)
		if action == nil {
			return nil
		}
		db.ActionMap[userId] = action
	}
	return action
}

func RecordAction(userId int32, actId int32) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("RecordAction err: ", err, string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
			debug.PrintStack()
		}
	}()
	if userId <= 0 || actId <= 0 {
		logger.Log.Errorf("=========RecordAction userId = %v, actId = %v", userId, actId)
		return
	}
	action := getAction(userId)
	if action == nil {
		action = &include.Action{UserId: userId, Actions: util.ToStr(actId) + "_1"}
	} else {
		if action.UserId <= 0 {
			logger.Log.Errorf("=========RecordAction userId = %v", userId)
			return
		}
		actions := action.Actions
		actConf := tableconfig.ActionConfigs.GetActionById(actId)
		if actConf.Type == constant.ActionSingle {
			if existAction(actions, actId) {
				return
			}
			actions = actions + "," + util.ToStr(actId) + "_1"
			action.Actions = actions
		} else {
			if actId != constant.Victory && actId != constant.Failure {
				if isOneGame(actions) {
					action.Actions = addActionNum(actions, actId)
				}
			} else {
				action.Actions = addActionNum(actions, actId)
			}
		}
	}

	saveAction(action)
}

func isOneGame(actions string) bool {
	if len(actions) == 0 {
		return true
	}
	acts := strings.Split(actions, ",")
	for _, v := range acts {
		acs := strings.Split(v, "_")
		if util.ToInt(acs[0]) == constant.Victory || util.ToInt(acs[0]) == constant.Failure {
			return false
		}
	}
	return true
}

func addActionNum(actions string, actId int32) string {
	var exist bool
	var arr []string
	acts := strings.Split(actions, ",")
	for _, v := range acts {
		acs := strings.Split(v, "_")
		if util.ToInt(acs[0]) == actId {
			num := util.ToInt(acs[1])
			num += 1
			arr = append(arr, acs[0]+"_"+util.ToStr(num))
			exist = true
		} else {
			arr = append(arr, acs[0]+"_"+acs[1])
		}
	}
	if !exist {
		actions = actions + "," + util.ToStr(actId) + "_1"
	} else {
		actions = strings.Join(arr, ",")
	}
	return actions
}

func existAction(actions string, actId int32) bool {
	if len(actions) == 0 {
		return false
	}
	acts := strings.Split(actions, ",")
	for _, v := range acts {
		acs := strings.Split(v, "_")
		if util.ToInt(acs[0]) == actId {
			return true
		}
	}
	return false
}
