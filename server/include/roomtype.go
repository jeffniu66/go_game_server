package include

import "go_game_server/proto3"

type VoteMsg struct {
	UserId     int32 // 投票者id
	VoteUserId int32 // 被投者id
	VoteNum    int32
	Status     proto3.CommonStatusEnum // 1-正常投票 2-弃票
}

type AttackReq struct {
	SocketUser   int32 // 连接者
	UserId       int32 // 杀人者
	SufferUserId int32 // 被杀人
}

type SettletData struct {
	OffLine          int32  // 是否掉线 0-未掉线， 1掉线
	GameDuration     int32  // 游戏时长
	MatchGameNum     int32  // 游戏盘数 默认是1
	MatchWinNum      int32  // 0或1
	MatchWolfNum     int32  // 0或1
	StarCount        int32  // 1或-1
	Star             int32  // 1或-1
	WolfWinNum       int32  // 0 1
	PoorWinNum       int32  // 0 1
	Exp              int32  // 经验
	Gold             int32  // 金币
	KillTotal        int32  // 杀人数
	WolfKillTotal    int32  // 狼人杀人数
	PoorKillTotal    int32  // 平民杀人数
	BevoteedTotal    int32  // 被票杀次数 0 1
	BekilledTotal    int32  // 被杀次数 0 1
	VoteTotal        int32  // 投票次数
	VoteCorrectTotal int32  // 投票正确次数
	VoteFailedTotal  int32  // 投票失败次数
	KeepFirstOut     int32  // 连续第一个被杀 0 1
	KeepWolf         int32  // 连续狼人 0 1
	KeepPoor         int32  // 连续平民 0 1
	KeepNoItem       int32  // 连续未得到道具
	TotalWolfDay     int32  // 累计狼人天数 0 1
	TotalTask        int32  // 累计完成任务数
	Rewards          string // 1,50|2,50
}
