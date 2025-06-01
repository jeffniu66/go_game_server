package db

import (
	"go_game_server/server/include"
)

func InsertAdData(o *include.AdData) {
	ExecDB(UserDBType, o.UserId, "insert into t_ad_data(user_id, ad_type, ad_num, register_time, stat_date)"+
		"values(?, ?, ?, ?, ?)", o.UserId, o.AdType, o.AdNum, o.RegisterTime, o.StatDate)
}

func SelectAdData(userId, adType int32, statDate string) *include.AdData {
	sql := `
	select
		uid, user_id, ad_type, ad_num
	from
		t_ad_data
	where
		user_id = ? and ad_type = ? and stat_date = ?
	`
	rows, _ := DB.Query(sql, userId, adType, statDate)
	defer rows.Close()
	var o *include.AdData
	for rows.Next() {
		o = &include.AdData{}
		_ = rows.Scan(&o.Uid, &o.UserId, &o.AdType, &o.AdNum)
	}
	return o
}

func UpdateAdData(o *include.AdData) {
	sql := `
	update
		t_ad_data
	set
		ad_num = ?
	where
		uid = ?
	`
	ExecDB(UserDBType, 0, sql, o.AdNum, o.Uid)
}
