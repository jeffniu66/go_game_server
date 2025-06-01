package constant

const (
	ActionSingle = 1
	ActionMulti  = 0
)

const (
	ClickMatch       = 1  // 点击匹配进入游戏
	ClickAchievement = 2  // 点击成就进入分页内
	ClickRole        = 3  // 点击角色进入分页内
	ClickStore       = 4  // 点击商店进入分页内
	ClickCancelMatch = 5  // 匹配界面--点击取消匹配
	ClickAuction     = 6  // 候场--点击抢一下
	ClickSkin        = 7  // 候场--点击皮肤页
	ClickLuckyCard   = 8  // 候场--点击幸运卡页
	IDClose          = 9  // 玩法--身份分配界面关闭
	KillPeople       = 10 // 玩法--杀人
	CallPolice       = 11 // 玩法--报警
	DoTask           = 12 // 玩法--触发任务
	InVote           = 13 // 玩法--进入投票界面
	Vote             = 14 // 玩法--投票任意玩家
	GiveUpVote       = 15 // 玩法--弃票
	VoteKill         = 16 // 玩法--被投死亡
	Killed           = 17 // 玩法--被杀死亡
	Victory          = 18 // 玩法--游戏胜利
	Failure          = 19 // 玩法--游戏失败
	AcqMatchReward   = 20 // 玩法--领取匹配福利
)
