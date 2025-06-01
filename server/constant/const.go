package constant

const (
	AIAddScoreInterTime    = 14  // 机器人自动加分数间隔时间
	TaskTotalScoreId       = 108 // 任务总分数的表id
	UrgencyTaskCdId        = 110 // 紧急任务CD的表id
	CloseDoorDurationTime  = 111 // 关门持续时间（秒）
	CloseDoorCD            = 112 // 关门CD
	NormalManSkillPool     = 300 // 平民技能总池
	WolfManSkillPool       = 350 // 狼人技能总池
	WolfManKillOneLvPool   = 351 // 狼人杀1-2人掉落池
	WolfManKillTwoLvPool   = 352 // 狼人杀2人以上掉落池
	DeathDropPool          = 353 // 死亡掉落池
	DeathDropSumPool       = 354 // 死亡掉落总池
	SmokeEgg               = 501 // 烟雾弹
	NormalBoxDropNum       = 901 // 普通宝箱掉落道具种类
	AdvanceBoxDropNum      = 902 // 高级宝箱掉落道具种类
	OpenBoxNeedGold        = 903 // 普通宝箱开箱费用（金币）
	BoxMaxAdNumPerDay      = 905 // 高级宝箱单日最多开启次数
	ViewAdGetGoldNum       = 906 // 免费金币-看一次广告获取金币数
	GoldMaxAdNumPerDay     = 907 // 免费金币-单日最多观看AD次数
	SkinMaxAdNumPerDay     = 908 //  广告皮肤单日观看上限
	OneAuctionPrice        = 909 // 竞拍单次出价金额
	AuctionCD              = 910 // 出价操作CD
	ReturnGoldRatio        = 911 // 返增金额占比（X分之一）
	AuctionDurationTime    = 912 // 竞拍持续时间
	AuctionDropIdNormalMan = 913 // 好人--竞拍奖品--掉落池ID
	AuctionDropIdWolfMan   = 914 // 狼人--竞拍奖品--掉落池ID
	AdTimeCheckInterval    = 951 // 广告时间校验间隔（秒）
	InitItems              = 999 // 初始档赠送内容
)

const (
	MatchRoomType      = 10   // 匹配赛類型
	GameWaitTime       = 8    // 候场时间
	RankAddAI          = 11   // 该段位直接补充剩下所有机器人
	RankAddAITime      = 12   // 匹配时间达到N时补充剩下所有机器人----弃用
	RobotChatTime      = 13   // 机器人聊天时间----------------------弃用
	AttackFrozenTime   = 101  // 狼人攻击CD
	UrgencyVoteNum     = 102  // 紧急会议次数
	UrgVoteFrozenTime  = 950  // 紧急会议CD
	GameVoteTime       = 204  // 投票时长
	ConstFPSNum        = 205  // 每秒多少帧数
	ChatVoteTime       = 206  // 投票前聊天时长
	VoteSumEndTime     = 207  // 投票结束时长
	ProtectRankID      = 208  // 段位保护
	WinAnnxGlodReward  = 400  // 任务奖励加成金币-放弃使用，转到排位赛表
	WinAnnxExpReward   = 401  // 任务奖励加成经验-放弃使用，转到排位赛表
	TaskFinishReward   = 402  // 任务完成数（单个加成）-放弃使用，转到排位赛表
	KillReward         = 403  // 杀人数（单个加成）-放弃使用，转到排位赛表
	TaskAllReward      = 404  // 总任务完成进度（单个加成）-放弃使用，转到排位赛表
	KillScoreRatio     = 915  // 狼人评分-击破系数
	BreakScoreRatio    = 916  // 狼人评分-破坏系数
	HidScoreRatio      = 917  // 狼人评分-逃票系数
	TaskScoreRatio     = 918  // 好人评分--任务系数
	VoteScoreRatio     = 919  // 好人评分--投凶成功系数
	TaskInScoreRatio   = 920  // 好人评分--任务占比系数
	CallVoteScoreRatio = 921  // 好人评分--发现尸体加分
	RoomCloseTime      = 922  // 房间关闭倒计时
	FreshGiftNeedNum   = 1000 // 新手福利-匹配领奖局数
	FreshGiftRewards   = 1001 // 新手福利-匹配领奖皮肤ID
	FreshGiftEndTime   = 1002 // 新手福利-倒计时领奖时长
	FreshGiftEndReward = 1003 // 新手福利-倒计时领奖奖品
	RandSkin           = 1004 // AI随机皮肤
	RandPhoto          = 1005 // AI随机头像
	RandTitle          = 1006 // AI随机称号
	DefaultSkin        = 1007 // 默认皮肤
	WorldChatTextCD    = 1010 // 公屏机器人快捷文本cd
	WorldChatTextRate  = 1011 // 公屏机器人快捷文本发言频率
)

const TimerRoomRickerPrefix = "room_ticker_" // 定时发送帧key前缀
const TimerAddScorePrefix = "add_score_"     // 机器人定时加积分前缀

const (
	MinBronzeRankStar    = 0
	MaxBronzeRankStar    = 11
	MinSliverRankStar    = 12
	MaxSliverRankStar    = 27
	MinGoldRankStar      = 28
	MaxGoldRankStar      = 48
	MinPlatinumRankStar  = 49
	MaxPlatinumRankStar  = 69
	MinDiamondRankStar   = 70
	MaxDiamondRankStar   = 90
	MinStarShineRankStar = 91
	MaxStarShineRankStar = 111
	MinKingRankStar      = 112
	InfRankStar          = -1 // 无限大
)
