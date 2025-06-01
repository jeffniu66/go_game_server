package db

import (
	"go_game_server/server/include"
)

func InsertGroupWinData(o *include.GroupWinData) {
	ExecDB(UserDBType, 0, "insert into t_win_data(wolf_man, normal_man, stat_date)"+
		"values(?, ?, ?)", o.WolfMan, o.NormalMan, o.StatData)
}

func SelectGroupWinData(statDate string) *include.GroupWinData {
	sql := `
	select
		*
	from
		t_win_data
	where
		stat_date = ? 
	`
	rows, _ := DB.Query(sql, statDate)
	defer rows.Close()
	var o *include.GroupWinData
	for rows.Next() {
		o = &include.GroupWinData{}
		_ = rows.Scan(&o.Uid, &o.WolfMan, &o.NormalMan, &o.StatData)
	}
	return o
}

func UpdateGroupWinData(o *include.GroupWinData) {
	sql := `
	update
		t_win_data
	set
		wolf_man = ?, normal_man = ? 
	where
		uid = ?
	`
	ExecDB(UserDBType, 0, sql, o.WolfMan, o.NormalMan, o.Uid)
}
