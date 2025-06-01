package db

import (
	"go_game_server/server/include"
)

var StoreMap = make(map[int32]*include.Store) // 玩家商店信息缓存 key: userId

func SaveStoreData(userId int32) interface{} {
	store, ok := StoreMap[userId]
	if !ok {
		return nil
	}
	if store.Update == include.Update {
		insertStore(store)
		store.Update = include.UnUpdate
	}
	return nil
}

func insertStore(store *include.Store) {
	ExecDB(UserDBType, store.UserId, "replace into t_store(user_id, box_use_ad_num, box_last_use_ad_time, mys_skin_buy_num, mys_skin_last_refresh_time, mys_skin_chip_id, "+
		"skin_use_ad_num, skin_last_refresh_time, gold_view_ad_num, gold_last_view_ad_time)"+
		"values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", store.UserId, store.BoxUseAdNum, store.BoxLastUseAdTime, store.MysSkinBuyNum, store.MysSkinLastRefreshTime, store.MysSkinChipId,
		store.SkinUseAdNum, store.SkinLastRefreshTime, store.GoldViewAdNum, store.GoldLastViewAdTime)
}

func SelectStore(userId int32) *include.Store {
	sql := `
	select
		user_id, box_use_ad_num, box_last_use_ad_time, mys_skin_buy_num, mys_skin_last_refresh_time, mys_skin_chip_id, skin_use_ad_num, skin_last_refresh_time, gold_view_ad_num, gold_last_view_ad_time
	from
		t_store
	where
		user_id = ?
	`
	rows, _ := DB.Query(sql, userId)
	defer rows.Close()
	var s *include.Store
	for rows.Next() {
		s = &include.Store{}
		_ = rows.Scan(&s.UserId, &s.BoxUseAdNum, &s.BoxLastUseAdTime, &s.MysSkinBuyNum, &s.MysSkinLastRefreshTime, &s.MysSkinChipId, &s.SkinUseAdNum,
			&s.SkinLastRefreshTime, &s.GoldViewAdNum, &s.GoldLastViewAdTime)
	}
	return s
}
