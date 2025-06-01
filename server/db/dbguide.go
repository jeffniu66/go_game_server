package db

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
)

func InsertGuide(guide *include.Guide) {
	ExecDB(UserDBType, guide.UserId, "replace into t_guide(user_id, guide_ids, update_time)"+
		"values(?, ?, ?)", guide.UserId, guide.GuideIds, util.UnixTime())
}

func SelectGuide(userId int32) *include.Guide {
	sql := `
	select
		user_id, guide_ids 
	from
		t_guide
	where
		user_id = ?
	`
	rows, _ := DB.Query(sql, userId)
	defer rows.Close()
	var g *include.Guide
	for rows.Next() {
		g = &include.Guide{}
		_ = rows.Scan(&g.UserId, &g.GuideIds)
	}
	return g
}
