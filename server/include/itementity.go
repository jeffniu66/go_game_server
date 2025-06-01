package include

type Item struct { // 玩家道具 持久化
	Uid    int32
	ItemId int32
	Num    int32
	Update int32 // 更新标志，1为需要回写数据库
	UserId int32
}
