package db

import (
	"database/sql"
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
)

type Attr struct {
	include.PlayerAttr
	UserTitle          *UserTitle
	AchievementMap     map[int32]int32          // key-archivetype value-archive_num
	UserArchiveTypeMap map[int32][]*UserArchive // key-archiveType value UserArchive
}

var NameMap = make(map[string]string) // key: name value: name

func (a *Attr) InitData(name, acctName, nameIndex string, sex int32, openid, channel string) *Attr {
	selSQL := `
	select 
		user_id, acct_name, server_no, username, sex, name_index, user_photo, user_border, gold, gemstone, level, exp, max_exp, rank_id, star, star_count, 
		his_rank_id, ninja_id, ninja_id_gift, archive_point, max_archive_point, use_skin, got_skins, 
		game_duration, match_game_num, match_win_num, match_wolf_num, wolf_win_num, poor_win_num, offline_num, vote_total, vote_correct_total, vote_failed_total, 
		kill_total, wolf_kill_total, pool_kill_total, bekilled_total, bevoteed_total, update_time, room_id, openid, channel, fresh_gift_step, fresh_end_time, 
		register_time
	FROM 
		t_user 
	WHERE 
		username = ?
	`
	rows, err := DB.Query(selSQL, name)
	defer rows.Close()
	util.CheckErr(err)
	emptyData := true
	var serverNo int32
	for rows.Next() {
		err = rows.Scan(&a.UserID, &a.AcctName, &serverNo, &a.Username, &a.Sex, &a.NameIndex, &a.UserPhoto, &a.UserBorder, &a.Gold, &a.GemStone, &a.Level, &a.Exp, &a.MaxExp, &a.RankID, &a.Star, &a.StarCount,
			&a.HisRankID, &a.NinjaID, &a.NinjaIDGift, &a.ArchivePoint, &a.MaxArchivePoint, &a.UseSkin, &a.GotSkins,
			&a.GameDuration, &a.MatchGameNum, &a.MatchWinNum, &a.MatchWolfNum, &a.WolfWinNum, &a.PoorWinNum, &a.OfflineNum, &a.VoteTotal, &a.VoteCorrectTotal, &a.VoteFailedTotal,
			&a.KillTotal, &a.WolfKillTotal, &a.PoorKillTotal, &a.BekilledTotal, &a.BevoteedTotal, &a.UpdateTime, &a.RoomID, &a.OpenId, &a.Channel, &a.FreshGiftStep, &a.FreshEndTime,
			&a.RegisterTime,
		)
		util.CheckErr(err)
		emptyData = false
	}
	logger.Log.Info("============", a.UserID, emptyData)
	curTime := util.UnixTime()
	a.LoginTime = curTime // 登录时间
	if emptyData {        // 创建新玩家
		userID := IDInstance.GetNewUserID()
		a.UserID = userID
		a.AcctName = acctName
		a.Username = name
		a.NameIndex = nameIndex
		a.Sex = sex
		defSkin := tableconfig.ConstsConfigs.GetIdValue(constant.DefaultSkin)
		if defSkin <= 0 {
			a.GotSkins = "1" // "1"是默认皮肤，
		} else {
			a.GotSkins = "1," + util.ToStr(defSkin) // "1"是默认皮肤，defskin是赠送皮肤
		}
		a.UseSkin = defSkin
		a.MaxExp = tableconfig.LevelConfigs.GetLevelConfig(1).NeedNum
		a.RankID = tableconfig.QuaLevelConfs.StartID
		a.NinjaID = tableconfig.NinjaConfigs.StartID
		a.Channel = channel
		ninja := tableconfig.NinjaConfigs.GetNinjaConfig(a.NinjaID)
		if ninja != nil {
			a.MaxArchivePoint = ninja.NeedNum
		}
		// 第一次登陆标志 -注册
		a.FirstLogin = true
		a.OpenId = openid
		a.SexModify = 1
		a.RegisterTime = util.UnixTime()

		logger.Log.Info("create new userID = ", userID)
		//a.RegistTime = curTime // 注册时间
		a.SaveUser()
	}
	// get title
	a.UserTitle = &UserTitle{}
	a.UserTitle.UserID = a.UserID
	a.UserTitle.GetUserTitle(a.UserID)

	// get archive
	a.UserArchiveTypeMap = GetUserArchivement(a.UserID)
	a.AchievementMap = make(map[int32]int32, 0)

	logger.Log.Infof("player：%v, user archive:%v", a, a.UserArchiveTypeMap)
	return a
}

func (a *Attr) SaveUser() interface{} {

	a.LogoutTime = util.UnixTime() // 登出时间
	serverNo := global.MyConfig.Read("server", "serverno")
	repSQL := `replace into t_user(user_id, acct_name, server_no, username, sex, name_index, user_photo, user_border, gold, gemstone, level, exp, max_exp, rank_id, star, star_count, 
		his_rank_id, ninja_id, archive_point, max_archive_point, use_skin, got_skins, openid, channel, fresh_gift_step, fresh_end_time, sex_modify, register_time)`
	ExecDB(UserDBType, a.UserID, repSQL+
		"values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		a.UserID, a.AcctName, util.ToInt(serverNo), a.Username, a.Sex, a.NameIndex, a.UserPhoto, a.UserBorder, a.Gold, a.GemStone, a.Level, a.Exp, a.MaxExp, a.RankID, a.Star, a.StarCount,
		a.HisRankID, a.NinjaID, a.ArchivePoint, a.MaxArchivePoint, a.UseSkin, a.GotSkins, a.OpenId, a.Channel, a.FreshGiftStep, a.FreshEndTime, a.SexModify, a.RegisterTime)
	return nil
}

func GetUserMaxID() int32 {
	rows := DB.QueryRow("select max(user_id)+1 as id from t_user")
	var num sql.NullInt64
	err := rows.Scan(&num)
	util.CheckErr(err)

	var userID int32
	// 因为这里有可能为空值
	if !num.Valid {
		userID = 100001
	} else {
		userID = int32(num.Int64)
	}

	// 32bit = 14bit serverNo + 18bit incrUid
	sererNo := global.MyConfig.ReadInt32("server", "serverno")
	userID |= sererNo << 18
	return userID
}

func JudgeUserName(name string) bool {
	exist := false
	err := DB.QueryRow("SELECT EXISTS (SELECT username FROM t_user WHERE username = ?)", name).Scan(&exist)
	//fmt.Println(err)
	util.CheckErr(err)
	return exist
}

func CountUsername(nameIndex string) int32 {
	var count int32 = -1
	err := DB.QueryRow("SELECT COUNT(*) FROM t_user WHERE name_index = ?", nameIndex).Scan(&count)
	util.CheckErr(err)
	return count
}

func GetUsername(userID int32) string {
	var s string
	err := DB.QueryRow("SELECT username FROM t_user WHERE user_id = ?", userID).Scan(&s)
	util.CheckErr(err)
	return s
}

func GetUserIdByUserName(name string) (userID int32) {
	err := DB.QueryRow("SELECT user_id FROM t_user WHERE username = ?", name).Scan(&userID)
	util.CheckErr(err)
	return
}

func GetUserByUserName(username string) (a *Attr) {
	sql := `
	select 
		user_id, acct_name, server_no, username, sex, name_index, user_photo, user_border, gold, gemstone, level, exp, max_exp, rank_id, star, star_count, 
		his_rank_id, ninja_id, ninja_id_gift, archive_point, max_archive_point, use_skin, got_skins, 
		game_duration, match_game_num, match_win_num, match_wolf_num, wolf_win_num, poor_win_num, offline_num, vote_total, vote_correct_total, vote_failed_total, 
		kill_total, wolf_kill_total, pool_kill_total, bekilled_total, bevoteed_total, update_time, room_id, openid, channel, fresh_gift_step, fresh_end_time, sex_modify,
		register_time
	FROM 
		t_user 
	WHERE 
		username = ?
	`
	rows, err := DB.Query(sql, username)
	defer rows.Close()
	util.CheckErr(err)
	var serverNo int32
	for rows.Next() {
		a = &Attr{}
		err = rows.Scan(&a.UserID, &a.AcctName, &serverNo, &a.Username, &a.Sex, &a.NameIndex, &a.UserPhoto, &a.UserBorder, &a.Gold, &a.GemStone, &a.Level, &a.Exp, &a.MaxExp, &a.RankID, &a.Star, &a.StarCount,
			&a.HisRankID, &a.NinjaID, &a.NinjaIDGift, &a.ArchivePoint, &a.MaxArchivePoint, &a.UseSkin, &a.GotSkins,
			&a.GameDuration, &a.MatchGameNum, &a.MatchWinNum, &a.MatchWolfNum, &a.WolfWinNum, &a.PoorWinNum, &a.OfflineNum, &a.VoteTotal, &a.VoteCorrectTotal, &a.VoteFailedTotal,
			&a.KillTotal, &a.WolfKillTotal, &a.PoorKillTotal, &a.BekilledTotal, &a.BevoteedTotal, &a.UpdateTime, &a.RoomID, &a.OpenId, &a.Channel, &a.FreshGiftStep, &a.FreshEndTime, &a.SexModify, &a.RegisterTime,
		)
		util.CheckErr(err)
	}
	return a
}

func GetUserByUserId(userId int32) (a *Attr) {
	sql := `
	select 
		tu.user_id, acct_name, server_no, username, sex, name_index, user_photo, user_border, gold, gemstone, level, exp, max_exp, rank_id, star, star_count, 
		his_rank_id, ninja_id, ninja_id_gift, archive_point, max_archive_point, use_skin, got_skins, 
		game_duration, match_game_num, match_win_num, match_wolf_num, wolf_win_num, poor_win_num, offline_num, vote_total, vote_correct_total, vote_failed_total, 
		kill_total, wolf_kill_total, pool_kill_total, bekilled_total, bevoteed_total, update_time, room_id, openid, channel, fresh_gift_step, fresh_end_time, sex_modify,
		register_time, 
		tut.use_title
	FROM 
		t_user AS tu
	LEFT JOIN
		t_user_title AS tut ON tu.user_id = tut.user_id
	WHERE 
		tu.user_id = ?
	`
	rows, err := DB.Query(sql, userId)
	defer rows.Close()
	util.CheckErr(err)
	var serverNo int32
	for rows.Next() {
		a = &Attr{}
		a.UserTitle = new(UserTitle)
		err = rows.Scan(&a.UserID, &a.AcctName, &serverNo, &a.Username, &a.Sex, &a.NameIndex, &a.UserPhoto, &a.UserBorder, &a.Gold, &a.GemStone, &a.Level, &a.Exp, &a.MaxExp, &a.RankID, &a.Star, &a.StarCount,
			&a.HisRankID, &a.NinjaID, &a.NinjaIDGift, &a.ArchivePoint, &a.MaxArchivePoint, &a.UseSkin, &a.GotSkins,
			&a.GameDuration, &a.MatchGameNum, &a.MatchWinNum, &a.MatchWolfNum, &a.WolfWinNum, &a.PoorWinNum, &a.OfflineNum, &a.VoteTotal, &a.VoteCorrectTotal, &a.VoteFailedTotal,
			&a.KillTotal, &a.WolfKillTotal, &a.PoorKillTotal, &a.BekilledTotal, &a.BevoteedTotal, &a.UpdateTime, &a.RoomID, &a.OpenId, &a.Channel, &a.FreshGiftStep, &a.FreshEndTime, &a.SexModify, &a.RegisterTime,
			&a.UserTitle.UseTitle,
		)
		util.CheckErr(err)
	}
	if a != nil {
		// get title
		a.UserTitle = &UserTitle{}
		a.UserTitle.UserID = a.UserID
		a.UserTitle.GetUserTitle(a.UserID)

		// get archive
		a.UserArchiveTypeMap = GetUserArchivement(a.UserID)
		a.AchievementMap = make(map[int32]int32, 0)
	}
	return a
}

func GetUserByOpenId(openid string) (userIdList []int32) {
	sql := `
	select
		user_id 
	from
		t_user
	where
		openid = ?
	`
	rows, _ := DB.Query(sql, openid)
	defer rows.Close()
	for rows.Next() {
		var userId int32
		_ = rows.Scan(&userId)
		userIdList = append(userIdList, userId)
	}
	return userIdList
}

func (a *Attr) UpdateUserAttr() {
	if a == nil {
		return
	}
	a.UpdateTime = util.UnixTime()
	uptSQL := `
	UPDATE 
		t_user 
	SET 
		username = ?, user_photo = ?, user_border = ?, gold = ?, gemstone = ?, level = ?, exp = ?, max_exp = ?, 
		rank_id = ?, star = ?, star_count = ?, his_rank_id = ?, ninja_id = ?, ninja_id_gift = ?, archive_point = ?, max_archive_point = ?,
		use_skin = ?, got_skins = ?, game_duration = ?, match_game_num = ?,  match_win_num = ?, match_wolf_num = ?, wolf_win_num = ?, poor_win_num = ?, 
		offline_num = ?, vote_total = ?, vote_correct_total = ?,  vote_failed_total = ?, kill_total = ?, wolf_kill_total = ?, 
		pool_kill_total = ?, bekilled_total = ?, bevoteed_total = ?, update_time = ?, room_id = ?, channel = ?, fresh_gift_step = ?, fresh_end_time = ?, sex_modify = ?, sex = ? 
	WHERE
		user_id = ?
	`
	ExecDB(UserDBType, a.UserID, uptSQL,
		a.Username, a.UserPhoto, a.UserBorder, a.Gold, a.GemStone, a.Level, a.Exp, a.MaxExp,
		a.RankID, a.Star, a.StarCount, a.HisRankID, a.NinjaID, a.NinjaIDGift, a.ArchivePoint, a.MaxArchivePoint,
		a.UseSkin, a.GotSkins, a.GameDuration, a.MatchGameNum, a.MatchWinNum, a.MatchWolfNum, a.WolfWinNum, a.PoorWinNum,
		a.OfflineNum, a.VoteTotal, a.VoteCorrectTotal, a.VoteFailedTotal, a.KillTotal, a.WolfKillTotal,
		a.PoorKillTotal, a.BekilledTotal, a.BevoteedTotal, a.UpdateTime, a.RoomID, a.Channel, a.FreshGiftStep, a.FreshEndTime, a.SexModify, a.Sex, a.UserID)
}

func (a *Attr) SaveUserAchiveList() {
	for _, list := range a.UserArchiveTypeMap {
		for _, v := range list {
			SaveUserArchive(v)
		}
	}
}

func (a *Attr) SaveData() {
	a.UpdateUserAttr()
	a.UserTitle.SaveUserTitle()
	a.SaveUserAchiveList()
}

func GetUserIds() (userIds []int32) {
	q := `select user_id from t_user`
	rows, _ := DB.Query(q)
	defer rows.Close()
	for rows.Next() {
		var userId int32
		_ = rows.Scan(&userId)
		userIds = append(userIds, userId)
	}
	return userIds
}

func SelectOnlineNum() (sum int32) {
	err := DB.QueryRow("SELECT count(1) FROM t_user WHERE update_time > ?", util.UnixTime()-5*60).Scan(&sum)
	util.CheckErr(err)
	return
}

func SelectRegisterNum(channel, startTime, endTime string) (sum int32) {
	err := DB.QueryRow("SELECT count(1) FROM t_user WHERE channel = ? and FROM_UNIXTIME(register_time) > ? and FROM_UNIXTIME(register_time) < ?", channel, startTime, endTime).Scan(&sum)
	util.CheckErr(err)
	return
}

func SelectNames() (names []string) {
	rows, _ := DB.Query("SELECT username FROM t_user")
	defer rows.Close()
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		names = append(names, name)
	}
	return
}

func SelectRankBoard(length int32, ret []*include.UserTop, retLen *int32) {
	var topId int32 = 0
	selSQL := `
	SELECT
		user_id, username,user_photo, level, exp, rank_id, star, star_count, update_time
	FROM
		t_user
	WhERE
		star_count >= 0
	ORDER BY
		star_count DESC, level DESC, exp DESC, update_time DESC
	Limit 0, ?;
	`
	rows, err := DB.Query(selSQL, length)
	util.CheckErr(err)
	defer rows.Close()
	for rows.Next() {
		topId++
		tmp := new(include.UserTop)
		tmp.TopId = topId
		err := rows.Scan(&tmp.UserId, &tmp.Username, &tmp.UserPhoto, &tmp.Level, &tmp.Exp, &tmp.RankId, &tmp.Star, &tmp.StarCount, &tmp.UpdateTime)
		util.CheckErr(err)
		ret[int(topId-1)] = tmp
	}
	*retLen = topId
	return
}
