package db

import (
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

type Good include.Goods

type GoodsMap map[int32]*Good

func (gMap GoodsMap) InitData(userId int32) GoodsMap {
	rows, err := DB.Query("select goods_id, type, subtype, num, location, create_time from t_user_goods where user_id = ?", userId)
	defer rows.Close()
	util.CheckErr(err)
	for rows.Next() {
		g := &Good{}
		err = rows.Scan(&g.GoodsId, &g.Type, &g.SubType, &g.Num, &g.Location, &g.CreateTime)
		util.CheckErr(err)
		gMap[g.GoodsId] = g
	}
	return gMap
}

func (gMap GoodsMap) SaveData(userId int32) interface{} {
	logger.Log.Infof("----------------------- save goods userId : %d, ", userId)
	for _, goods := range gMap {
		if goods.Update == include.Update {
			InsertGoods(goods, userId)
			goods.Update = include.UnUpdate
		}
	}
	return nil
}

// 添加物品
func InsertGoods(goods *Good, userId int32) {
	ExecDB(UserDBType, userId, "replace into t_user_goods(user_id, goods_id, type, subtype, num, location, create_time) "+
		"values(?, ?, ?, ?, ?, ?, ?)",
		userId, goods.GoodsId, goods.Type, goods.SubType, goods.Num, goods.Location, goods.CreateTime)
}

// 删除物品
func DeleteGoods(goods *Good, userId int32) {
	ExecDB(UserDBType, userId, "delete from t_user_goods where user_id = ? and goods_id = ? ", userId, goods.GoodsId)
}
