package db

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
)

type UserTitle include.UserTitle

func (u *UserTitle) SaveUserTitle() {
	if u == nil {
		return
	}
	saveSQL := "replace into t_user_title(user_id, keep_first_out, keep_wolf, keep_poor, keep_noitem, "
	saveSQL += "total_kill_poor, total_wolf_day, wolf_timestamp, total_task, total_soul_task, total_gold, total_ad, total_archive, use_title, got_titles, skin_red_datas, title_red_datas)"
	saveSQL += `values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	ExecDB(UserDBType, u.UserID, saveSQL, u.UserID, u.KeepFirstOut, u.KeepWolf, u.KeepPoor, u.KeepNoItem,
		u.TotalKillPoor, u.TotalWolfDay, u.WolfTimestamp, u.TotalTask, u.TotalSoulTask, u.TotalGold, u.TotalAd, u.TotalArchive, u.UseTitle, u.GotTitles, u.SkinRedData, u.TitleRedData)
}

func (u *UserTitle) GetUserTitle(userID int32) *UserTitle {
	selSQL := `
SELECT 
	user_id, keep_first_out, keep_wolf, keep_poor, keep_noitem,
	total_kill_poor, total_wolf_day, wolf_timestamp, total_task, total_soul_task, total_gold, total_ad, total_archive, use_title, got_titles, skin_red_datas, title_red_datas
FROM
	t_user_title
WHERE
	user_id = ?
	`
	rows, err := DB.Query(selSQL, userID)
	defer rows.Close()
	util.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&u.UserID, &u.KeepFirstOut, &u.KeepWolf, &u.KeepPoor, &u.KeepNoItem,
			&u.TotalKillPoor, &u.TotalWolfDay, &u.WolfTimestamp, &u.TotalTask, &u.TotalSoulTask, &u.TotalGold, &u.TotalAd, &u.TotalArchive, &u.UseTitle, &u.GotTitles, &u.SkinRedData, &u.TitleRedData)
		util.CheckErr(err)
	}
	return u
}

type UserArchive include.UserArchive

func GetUserArchivement(userID int32) map[int32][]*UserArchive {
	selSQL := `
SELECT 
	user_id, archive_type, archive_id, got_status, archive_next
FROM
	t_user_archive
WHERE 
	user_id = ?
	`
	rows, err := DB.Query(selSQL, userID)
	defer rows.Close()
	util.CheckErr(err)
	retMap := make(map[int32][]*UserArchive, 0)
	for rows.Next() {
		tmp := UserArchive{}
		err = rows.Scan(&tmp.UserID, &tmp.ArchiveType, &tmp.ArchiveID, &tmp.GotStatus, &tmp.ArchiveNext)
		util.CheckErr(err)

		retTmp := make([]*UserArchive, 0)
		if _, ok := retMap[tmp.ArchiveType]; ok {
			retTmp = retMap[tmp.ArchiveType]
		}
		retTmp = append(retTmp, &tmp)
		retMap[tmp.ArchiveType] = retTmp

	}
	return retMap
}

func SaveUserArchive(uArchive *UserArchive) {
	if uArchive == nil {
		return
	}
	saveSQL := "insert into t_user_archive(user_id, archive_type, archive_id, got_status, archive_next)VALUES(?,?,?,?,?)"
	saveSQL += " ON DUPLICATE KEY UPDATE got_status = ? "

	ExecDB(UserDBType, uArchive.UserID, saveSQL,
		uArchive.UserID, uArchive.ArchiveType, uArchive.ArchiveID, uArchive.GotStatus, uArchive.ArchiveNext,
		uArchive.GotStatus)
}

func GetDan(startTime, endTime int32) (dans []*include.DanData) {
	sql := `
		SELECT 
			user_id, username, rank_id, star
		FROM
			t_user
		WHERE 
			register_time > ? and register_time < ?
			`
	rows, err := DB.Query(sql, startTime, endTime)
	defer rows.Close()
	util.CheckErr(err)
	for rows.Next() {
		tmp := &include.DanData{}
		err = rows.Scan(&tmp.UserId, &tmp.Username, &tmp.RankId, &tmp.Star)
		dans = append(dans, tmp)
	}
	return dans
}
