package db

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
)

func InsertGameData(o *include.GameData) {
	ExecDB(UserDBType, o.UserId, "insert into t_game_data(user_id, username, game_num, login_date, register_time, update_time)"+
		"values(?, ?, ?, ?, ?, ?)", o.UserId, o.Username, o.GameNum, o.LoginDate, o.RegisterTime, util.UnixTime())
}

func SelectGameData(userId int32, loginDate string) *include.GameData {
	sql := `
	select
		uid, user_id, game_num
	from
		t_game_data
	where
		user_id = ? and login_date = ?
	`
	rows, _ := DB.Query(sql, userId, loginDate)
	defer rows.Close()
	var g *include.GameData
	for rows.Next() {
		g = &include.GameData{}
		_ = rows.Scan(&g.Uid, &g.UserId, &g.GameNum)
	}
	return g
}

func UpdateGameData(o *include.GameData) {
	sql := `
	update
		t_game_data
	set
		game_num = ?
	where
		uid = ?
	`
	ExecDB(UserDBType, 0, sql, o.GameNum, o.Uid)
}
