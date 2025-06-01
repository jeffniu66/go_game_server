package db

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
)

var ActionMap = make(map[int32]*include.Action) // key: userId

func SaveActionData(userId int32) interface{} {
	action, ok := ActionMap[userId]
	if !ok {
		return nil
	}
	if action.Update == include.Update {
		insertAction(action)
		action.Update = include.UnUpdate
	}
	return nil
}

func insertAction(a *include.Action) {
	ExecDB(UserDBType, a.UserId, "replace into t_action(user_id, actions, update_time)"+
		"values(?, ?, ?)", a.UserId, a.Actions, util.UnixTime())
}

func SelectAction(userId int32) *include.Action {
	sql := `
	select
		user_id, actions
	from
		t_action
	where
		user_id = ?
	`
	rows, _ := DB.Query(sql, userId)
	defer rows.Close()
	var g *include.Action
	for rows.Next() {
		g = &include.Action{}
		_ = rows.Scan(&g.UserId, &g.Actions)
	}
	return g
}
