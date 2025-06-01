package include

const (
	UnUpdate  = iota // 不需要回写数据库
	Update           // 需要回写数据库
	DelUpdate        // 需要删除数据库
)

const (
	Offline = iota // 不在线
	Online         // 在线
)

const (
	Normal  = iota // 普通战役关卡
	Special        // 精英战役关卡
)

// 发送给玩家协程的消息名
const (
	PlayerRefreshIncident            = "refreshIncident"            // 刷新玩家事件
	PlayerAllyInvite                 = "allyInvite"                 // 同盟邀请
	PlayerZeroRefresh                = "zeroRefresh"                // 0点刷新时间
	PlayerFourRefresh                = "fourRefresh"                // 4点刷新时间
	PlayerAddGoods                   = "addGoods"                   // 添加道具
	PlayerCostGoods                  = "costGoods"                  // 玩家扣除道具
	PlayerWorldBoxes                 = "worldBoxes"                 // 地图宝箱
	PlayerRefreshLandEffect          = "refreshLandEffect"          // 刷新玩家地块效果
	PlayerAddLand                    = "addLand"                    // 添加地块
	PlayerRemoveLand                 = "removeLand"                 // 删除地块
	PlayerAddHeroStrength            = "addHeroStrength"            // 玩家增加体力
	PlayerSetScienceLevel            = "setScienceLevel"            // 升级玩家科技等级
	PlayerAddBasicAttr               = "addBasicAttr"               // 玩家增加基础资源
	PlayerSaveDbData                 = "saveUserDbData"             // 定时保存玩家数据
	PlayerHeroGradeUp                = "heroGradeUp"                // 玩家英雄进阶
	PlayerAddWarMail                 = "addWarMail"                 // 添加战报邮件
	PlayerAddMail                    = "addMail"                    // 添加邮件
	PlayerUpdatePower                = "player_update_power"        // 玩家更新势力值
	PlayerSetFallEmoji               = "player_set_fall_emoji"      // 设置玩家沦陷表情
	PlayerUpdateTaskProgress         = "updateTaskProgress"         // 从同盟同步任务数据
	PlayerIncreaseResource           = "increaseResource"           // 同步资源增长
	PlayerUpdateTaskAttackCity       = "updateTaskAttackCity"       // 更新玩家同盟攻城任务
	PlayerRegisterActivityRewardMail = "registerActivityRewardMail" // 以玩家注册时间算结束的活动发邮件奖励
	PlayerUpdateTaskFinish           = "UpdateTaskFinish"           // 检查玩家的任务是否完成并更新状态
)

const (
	LevelAttr                = 20 // 等级
	ExpAttr                  = 21 // 经验值
	WoodAttr                 = 22 // 木材
	IronAttr                 = 23 // 铁矿
	StoneAttr                = 24 // 石料
	ForageAttr               = 25 // 粮草
	GoldAttr                 = 26 // 金币
	DiamondAttr              = 27 // 钻石
	BindDiamondAttr          = 28 // 绑定钻石(保留，但不使用)
	DecreeAttr               = 29 // 政令
	ArmyOrderAttr            = 30 // 军令
	PowerAttr                = 31 // 势力值
	DomainAttr               = 32 // 领地个数
	RenownAttr               = 33 // 名望
	HeroicStrength           = 34 // 英雄体力
	RenownLimitAttr          = 35 // 名望上限
	LordLevelAttr            = 36 // 领主的爵位
	FortressAttr             = 37 // 要塞
	FeatsAttr                = 38 // 武勋
	AllyAttr                 = 39 // 同盟id
	AllyPosAttr              = 40 // 同盟职位
	BirthAttr                = 41 // 出生点
	EnergyAttr               = 42 // 精力
	NormalSectionAttr        = 43 // 普通战役关卡
	SpecialSectionAttr       = 44 // 精英战役关卡
	MasterAllyIDAttr         = 45 //	上司同盟ID
	ConscriptSpeedAttr       = 46 // 征兵加速
	ConscriptLimitSpeedAttr  = 47 // 征兵加速上限
	CityMasterIndexAttr      = 51 // 城主地块ID
	MapExploreDailyTimesAttr = 61 // 大地图每日探索次数
	NextMapExploreStampAttr  = 62 // 下次大地图探索时间戳
)

// 玩家属性
type PlayerAttr struct {
	UserID        int32
	AcctName      string // 账号名
	Username      string // 角色名名
	NameIndex     string // 角色名索引 对应于一个随机名，空的情况下无对于
	Sex           int32
	GemStone      int32  // 宝石
	Country       int32  // 所属国家
	Power         int32  // 势力值
	FirstLogin    bool   // 第一次登陆
	LoginTime     int32  // 最近登录时间
	LogoutTime    int32  // 最近登出时间
	RegistTime    int32  // 注册时间
	CommonFlag    int32  // 通用的位存储标识
	Ip            string // ip地址
	UserBorder    int32  // 相框
	UserPhoto     int32  // 玩家头像
	UseSkin       int32  // 使用的皮肤
	GotSkins      string // 皮肤ids:2,3
	UserData             // 玩家匹配赛数据
	OpenId        string // 渠道用户唯一标识
	Channel       string // 渠道
	SexModify     int32  // 1未修改 2已修改
	FreshGiftStep int32
	FreshEndTime  int32
	RegisterTime  int32
}

// UserMatchGameData 玩家匹配赛数据
type UserData struct {
	Level            int32 // 等级
	Exp              int32 // 经验
	MaxExp           int32 // 最大经验
	Gold             int32 // 金币
	RankID           int32 // 段位
	Star             int32 // 星级
	StarCount        int32 // 总星级
	HisRankID        int32 // 历史段位
	NinjaID          int32 // 忍阶
	NinjaIDGift      int32 // 忍者礼包-对应的ninjaid
	ArchivePoint     int32 // 成就点
	MaxArchivePoint  int32 // 最大成就点
	GameDuration     int32 // 游戏时长
	MatchGameNum     int32 // 排位总局数
	MatchWinNum      int32 // 排位胜场总局数
	MatchWolfNum     int32 // 狼人总局数
	WolfWinNum       int32 // 排位狼人胜场局数
	PoorWinNum       int32 // 排位平民胜场局数
	OfflineNum       int32 // 掉线局数
	VoteTotal        int32 // 总投票次数
	VoteCorrectTotal int32 // 投票正确次数
	VoteFailedTotal  int32 // 投票失败次数
	KillTotal        int32 // 杀人次数
	WolfKillTotal    int32 // 狼人杀人次数
	PoorKillTotal    int32 // 平民杀人次数
	BekilledTotal    int32 // 被杀害次数
	BevoteedTotal    int32 // 被票杀次数
	UpdateTime       int32 // 最后更新时间
	RoomID           int32 // roomID
}

// 离线玩家活跃数据
type RedisUser struct {
	LordName            string // 领主名字
	Domain              int32  // 领地个数
	Country             int32  // 所属国家
	Power               int32  // 势力值
	AllyId              int32  // 同盟id
	AllyName            string // 同盟名
	AllyReqCD           int32  // cd
	AllyPos             int32  // 同盟职位
	BirthPoint          int32  // 出生点
	FirstOccupy         int32  // 玩家的第一次占领等级土地
	RegistTime          int32  // 注册时间
	LordLevel           int32  // 领主的爵位
	PlunderProtectStamp int32  // 突击掠夺保护时间戳
	Language            int32  // 选择文本语言
	FirebaseToken       string // firebase token
	NotifySetting       int32  //
	AccountName         string // 账号名
	CommonFlag          int32  // 通用的位存储标识
	CityMasterIndex     int32  // 城主府地块ID
	HeadId              int32  // head id
}

type UserTitle struct {
	UserID        int32  // 玩家ID
	KeepFirstOut  int32  // 连续第一个出局数
	KeepWolf      int32  // 连续狼人数
	KeepPoor      int32  // 连续平民数
	KeepNoItem    int32  // 连续未得到道具
	TotalKillPoor int32  // 累计杀平民
	TotalWolfDay  int32  // 累计几天狼人
	WolfTimestamp int32  // 累计狼人日期
	TotalTask     int32  // 累积完成任务
	TotalSoulTask int32  // 灵魂状态帮助队友次数
	TotalGold     int32  // 累计金币
	TotalArchive  int32  // 累计获得成就点
	TotalAd       int32  // 累计看广告数量
	UseTitle      int32  // 使用的称号
	GotTitles     string // 获得的称号，称号id
	SkinRedData   string // 皮肤红点数据 1,2
	TitleRedData  string // 称号红点数据 1,2
}

type UserArchive struct {
	UserID      int32
	ArchiveType int32
	ArchiveID   int32
	GotStatus   int32 // 0-已领取，1-未领取, 2-未达成
	ArchiveNext int32
}

// 用户排名
type UserTop struct {
	UserId     int32
	Username   string // 用户名
	UserPhoto  int32  // 头像
	TopId      int32  // 名次
	RankId     int32  // 段位
	Star       int32  // 星
	StarCount  int32  // 总星级
	Level      int32  // 等级
	Exp        int32  // 经验
	UpdateTime int32  // 最近更新时间
}

const (
	GetStatusYes   = 0 // 已领取
	GetStatusNo    = 1 // 未领取
	GetStatusFalse = 2 // 未达成
)

const (
	ArchiveTypeKill       = 1  // 杀人数/(狼人身份*8)≥
	ArchiveTypeVoteOut    = 2  // 被投票出局次数/平民玩家局数≥
	ArchiveTypeVoteWolf   = 3  // 成功投出狼人次数/（平民玩家局数*2）≥
	ArchiveTypeFirstOut   = 4  // 首轮就被杀死的次数/平民玩家局数≥
	ArchiveTypeTask       = 5  // 累计完成的任务次数≥
	ArchiveTypeTaskGold   = 6  // 在任务中累计获得金币数量≥
	ArchiveTypeSoulTask   = 7  // 在死亡状态中累计完成的任务数量≥
	ArchiveTypeLevel      = 8  // 等级≥
	ArchiveTypeNinjaLevel = 9  // 忍阶≥
	ArchiveTypeArcpoint   = 10 // 累计获得成就点≥
	ArchiveTypeKillPoor   = 12 // 累计击杀平民数量≥
	ArchiveTypeKillWolf   = 13 // 累计击杀狼人数量≥
	ArchiveTypeVoteTrue   = 14 // 累计投票成功次数≥
	ArchiveTypeGold       = 15 // 累计获得金币≥
	ArchiveTypeSkin       = 16 // 累计获得皮肤数量≥
	ArchiveTypeMatch      = 17 // 累计排位赛局数≥
	ArchiveTypeMatchWin   = 18 // 累计排位赛胜利局数≥
	ArchiveTypeRank       = 19 // 最高段位达成≥
	ArchiveTypeItem       = 20 // 累计获得道具数量≥
	ArchiveTypeDuration   = 22 // 累计在线时长≥
	ArchiveTypeAd         = 23 // 累计观看广告数量≥
	ArchiveTypePhoto      = 24 // 累计获得头像数量≥
	ArchiveTypeBorder     = 25 // 累计获得头像框数量≥
	ArchiveTypeKeepPoor   = 26 // 连续平民次数
	ArchiveTypeKeepWolf   = 27 // 连续狼人次数
	ArchiveTypeKeepNoItem = 28 // 连续未得到道具
	ArchiveTypeWolfDay    = 29 // 累计狼人天数
	ArchiveTypeVotePoor   = 30 // 投出平民数

	archiveTypeLast = 31 // 最后的类型
)

var ArchiveTypeList []int32

const (
	TitleType   = 1 // 称号
	ArchiveType = 2 // 成就
)

func init() {
	for i := int32(1); i < archiveTypeLast; i++ {
		ArchiveTypeList = append(ArchiveTypeList, i)
	}
}
