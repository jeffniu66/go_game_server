package db

import (
	"fmt"
	"go_game_server/server/include"
	"go_game_server/server/util"
	"strings"
)

func InsertItemData(o *include.ItemData) {
	ExecDB(UserDBType, o.UserId, "insert into t_item_data(user_id, item_id, num, add_time)"+
		"values(?, ?, ?, ?)", o.UserId, o.ItemId, o.Num, util.UnixTime())
}

func BatchInsertItem(items []*include.ItemData) {
	// 存放 (?, ?, ?, ?) 的slice
	valueStrings := make([]string, 0, len(items))
	// 存放values的slice
	valueArgs := make([]interface{}, 0, len(items)*2)
	time := util.UnixTime()
	for _, u := range items {
		valueStrings = append(valueStrings, "(?, ?, ?, ?)")
		valueArgs = append(valueArgs, u.UserId)
		valueArgs = append(valueArgs, u.ItemId)
		valueArgs = append(valueArgs, u.Num)
		valueArgs = append(valueArgs, time)
	}
	sql := fmt.Sprintf("INSERT INTO t_item_data(user_id, item_id, num, add_time) VALUES %s", strings.Join(valueStrings, ","))
	ExecDB(UserDBType, 0, sql, valueArgs...)
}
