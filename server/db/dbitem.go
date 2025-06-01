package db

import (
	"database/sql"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/util"
)

var ItemMap = make(map[int32]map[int32]*include.Item) // 全局道具 key1: userId  key2: itemId

func SaveItemsData(userId int32) interface{} {
	itemMap, ok := ItemMap[userId]
	if !ok {
		return nil
	}
	for _, item := range itemMap {
		if item.Update == include.Update {
			InsertItem(item)
			item.Update = include.UnUpdate
		}
	}
	return nil
}

func GetItemMaxID() int32 {
	rows := DB.QueryRow("select max(uid)+1 as id from t_item")
	var num sql.NullInt64
	err := rows.Scan(&num)
	util.CheckErr(err)

	var uid int32
	// 因为这里有可能为空值
	if !num.Valid {
		uid = 1
	} else {
		uid = int32(num.Int64)
	}

	// 32bit = 14bit serverNo + 18bit incrUid
	sererNo := global.MyConfig.ReadInt32("server", "serverno")
	uid |= sererNo << 18
	return uid
}

// 添加物品
func InsertItem(item *include.Item) {
	ExecDB(UserDBType, item.UserId, "replace into t_item(uid, item_id, num, user_id)"+
		"values(?, ?, ?, ?)", item.Uid, item.ItemId, item.Num, item.UserId)
	//query := "replace into t_item(item_id, num, user_id)" +
	//	"values(?, ?, ?)"
	//result, err := DB.Exec(query, item.ItemId, item.Num, item.UserId)
	//if err != nil {
	//	logger.Log.Errorf("exec db err:%v, query = %s", err, query)
	//}
	//lastId, _ := result.LastInsertId()
	//item.Uid = int32(lastId)
}

//func SaveData(i *include.Item) {
//	ExecDBAsync(UserDBType, i.UserId, "replace into t_item(item_id, num, user_id)"+
//		"values(?, ?, ?)", i.ItemId, i.Num, i.UserId)
//}

func GetPlayerItems(userId int32) []*include.Item {
	sql := `
	select
		uid, item_id, num, user_id
	from
		t_item
	where
		user_id = ?
	`
	rows, _ := DB.Query(sql, userId)
	defer rows.Close()
	var items []*include.Item
	for rows.Next() {
		item := &include.Item{}
		_ = rows.Scan(&item.Uid, &item.ItemId, &item.Num, &item.UserId)
		items = append(items, item)
	}
	return items
}
