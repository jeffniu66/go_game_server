package include

type Auction struct {
	CurPrice       int32           // 当前价格
	UserId         int32           // 玩家id
	UserName       string          // 玩家角色名
	EndTime        int32           // 结束时间
	ClickSum       int32           // 总竞拍次数
	AuctionTimeMap map[int32]int32 // 竞拍时间
	ClickNumMap    map[int32]int32 // 玩家点击次数
}
