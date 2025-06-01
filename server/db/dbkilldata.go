package db

import "go_game_server/server/include"

func InsertKillData(o *include.KillData) {
	ExecDB(UserDBType, 0, "insert into t_kill_data(user_id, rank_id, kill_num, stat_date)"+
		"values(?, ?, ?, ?)", o.UserId, o.RankId, o.KillNum, o.StatDate)
}

func SelectKillData(userId int32, statDate string) *include.KillData {
	sql := `
	select
		*
	from
		t_kill_data
	where
		user_id = ? and stat_date = ? 
	`
	rows, _ := DB.Query(sql, userId, statDate)
	defer rows.Close()
	var o *include.KillData
	for rows.Next() {
		o = &include.KillData{}
		_ = rows.Scan(&o.Uid, &o.UserId, &o.RankId, &o.KillNum, &o.StatDate)
	}
	return o
}

func UpdateKillData(o *include.KillData) {
	sql := `
	update
		t_kill_data
	set
		kill_num = ? 
	where
		uid = ?
	`
	ExecDB(UserDBType, 0, sql, o.KillNum, o.Uid)
}
