package game

import "go_game_server/server/db"

type GoodsBag struct {
	GoodsMap db.GoodsMap
}

func InitGoodsData(userId int32) *GoodsBag {
	goods := make(db.GoodsMap)
	goodsMap := goods.InitData(userId)
	return &GoodsBag{GoodsMap: goodsMap}
}
