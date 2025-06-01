package db

import (
	"go_game_server/server/include"
)

func InsertLoginData(o *include.LoginData) {
	ExecDB(UserDBType, o.UserId, "insert into t_login_data(user_id, username, login_date, register_time)"+
		"values(?, ?, ?, ?)", o.UserId, o.Username, o.LoginDate, o.RegisterTime)
}

func SelectLoginData(userId int32, loginDate string) *include.LoginData {
	sql := `
	select
		*
	from
		t_login_data
	where
		user_id = ? and login_date = ?
	`
	rows, _ := DB.Query(sql, userId, loginDate)
	defer rows.Close()
	var o *include.LoginData
	for rows.Next() {
		o = &include.LoginData{}
		_ = rows.Scan(&o.Uid, &o.UserId, &o.LoginDate)
	}
	return o
}
