package db

import "go_game_server/server/include"

func InsertConferData(o *include.ConferData) {
	ExecDB(UserDBType, 0, "insert into t_confer_data(user_id, rank_id, confer_num, stat_date, register_time)"+
		"values(?, ?, ?, ?, ?)", o.UserId, o.RankId, o.ConferNum, o.StatDate, o.RegisterDate)
}
