package include

// 玩家行为
type Action struct {
	UserId     int32  // 玩家id
	Actions    string // 行为
	UpdateTime int32
	Update     int32 // 更新标志，1为需要回写数据库
}
