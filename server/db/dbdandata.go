package db

import "go_game_server/server/include"

func InsertDanStat(o *include.DanData) {
	ExecDB(UserDBType, o.UserId, "insert into t_dan_stat(user_id, username, rank_id, star, register_date)"+
		"values(?, ?, ?, ?, ?)", o.UserId, o.Username, o.RankId, o.Star, o.RegisterDate)
}
