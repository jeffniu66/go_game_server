package db

import (
	"go_game_server/proto3"
	"go_game_server/server/include"
	"go_game_server/server/util"
)

type UserMail include.UserMail

func (u *UserMail) ToProto3UserMail() (ret *proto3.UserMail) {
	ret = &proto3.UserMail{}
	if u == nil {
		ret = nil
		return
	}
	ret.Id = u.Uid
	ret.Content = u.Content
	ret.IsRead = proto3.CommonStatusEnum(u.IsRead)
	ret.IsOpen = proto3.CommonStatusEnum(u.AnnexOpen)
	ret.DetailType = proto3.DetailTypeEnum(u.DetailType)
	ret.Title = u.Title
	ret.ItemIds = u.AnnexItems
	ret.CreateTime = u.CreateTime
	return
}

func SaveMail(userMail *UserMail) {
	userMail.CreateTime = util.UnixTime()
	userMail.UpdateTime = util.UnixTime()
	if userMail == nil || userMail.UserID <= 0 {
		return
	}
	saveSQL := "replace into t_user_mail(user_id, detail_type, title, content, annex_items, annex_open, is_read, is_expire, create_time, update_time)"
	saveSQL += "values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	ExecDB(UserDBType, userMail.UserID, saveSQL,
		userMail.UserID, userMail.DetailType, userMail.Title, userMail.Content, userMail.AnnexItems, userMail.AnnexOpen, userMail.IsRead, userMail.IsExpire, userMail.CreateTime, userMail.UpdateTime)
}

func GetMailList(userID int32) []*UserMail {
	ret := make([]*UserMail, 0)
	selSQL := `
SELECT 
	uid, user_id, detail_type, title, content, annex_items, annex_open, is_read, is_expire, create_time, update_time
FROM
	t_user_mail
WHERE
	user_id = ? AND is_expire = 0
	`
	rows, err := DB.Query(selSQL, userID)
	defer rows.Close()
	util.CheckErr(err)

	for rows.Next() {
		tmp := UserMail{}
		err = rows.Scan(&tmp.Uid, &tmp.UserID, &tmp.DetailType, &tmp.Title, &tmp.Content, &tmp.AnnexItems, &tmp.AnnexOpen, &tmp.IsRead, &tmp.IsExpire, &tmp.CreateTime, &tmp.UpdateTime)
		util.CheckErr(err)
		ret = append(ret, &tmp)
	}
	if len(ret) <= 0 {
		return nil
	}
	return ret
}

func GetMail(Uid int32) *UserMail {
	var ret UserMail
	selSQL := `
SELECT 
	uid, user_id, detail_type, title, content, annex_items, annex_open, is_read, is_expire, create_time, update_time
FROM
	t_user_mail
WHERE
	uid = ? AND is_expire = 0
	`
	rows, err := DB.Query(selSQL, Uid)
	defer rows.Close()
	util.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&ret.Uid, &ret.UserID, &ret.DetailType, &ret.Title, &ret.Content, &ret.AnnexItems, &ret.AnnexOpen, &ret.IsRead, &ret.IsExpire, &ret.CreateTime, &ret.UpdateTime)
		util.CheckErr(err)
	}
	if ret.Uid <= 0 {
		return nil
	}
	return &ret
}

func ReadUserMail(uid, userID int32, isRead int32) {
	if uid < 0 {
		return
	}
	uptSQL := `UPDATE t_user_mail SET is_read = ?, update_time = ? WHERE uid = ?`
	ExecDB(UserDBType, userID, uptSQL, isRead, util.UnixTime(), uid)
}

func OpenOneMail(uid, userID int32, isOpen int32) {
	if uid < 0 {
		return
	}

	uptSQL := `UPDATE t_user_mail SET is_read = 1, annex_open = ?,update_time = ? WHERE uid = ?`
	ExecDB(UserDBType, userID, uptSQL, isOpen, util.UnixTime(), uid)
}

// 获取未打开邮件
func GetNoOpenMail(userID int32) []string {
	selSQL := `SELECT annex_items FROM t_user_mail WHERE user_id = ? AND annex_open = 0`
	rows, err := DB.Query(selSQL, userID)
	defer rows.Close()
	util.CheckErr(err)
	ret := make([]string, 0)
	for rows.Next() {
		tmp := ""
		err = rows.Scan(&tmp)
		util.CheckErr(err)
		ret = append(ret, tmp)
	}
	if len(ret) <= 0 {
		return nil
	}
	return ret
}

func OpenAllMail(userID int32) {
	uptSQL := `UPDATE t_user_mail SET is_read = 1, annex_open = ?, update_time = ? WHERE user_id = ?`
	ExecDB(UserDBType, userID, uptSQL, proto3.CommonStatusEnum_true, util.UnixTime(), userID)
}

func DelAllRead(userID int32) {
	uptSQL := `DELETE FROM t_user_mail WHERE user_id = ? and is_read = 1 and annex_open = 1`
	ExecDB(UserDBType, userID, uptSQL, userID)
}
