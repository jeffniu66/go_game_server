package cfg

// 带上值，去掉前端用的字段
const (
	SysProDuration                    = iota + 1 //资源产出周期
	SysLockBuild                      = 2        //未解锁建筑
	SysUnlockBuild                    = 3        //已解锁建筑
	SysResourceBuild                  = 4        //资源田枚举
	SysConscriptionTime               = 5        //基础征兵时间（毫秒）
	SysHeroMaxLevel                   = 6        //英雄最大等级
	SysHeroSkill3UnLock               = 7        //英雄技能3解锁(3技能2星20级)
	SysHeroSkill4UnLock               = 8        //英雄技能4解锁(4技能3星紫色品阶)
	SysGuideBattleReward              = 9        // 引导战斗结束奖励ID及奖励物品ID和数量
	SysGuildResident                  = 10       // 引导用的居民(星级^类型^武力^智力^生产力)
	SysNewbieResident                 = 11       // 新手送的居民(星级^类型^武力^智力^生产力)
	SysMaxViewRadius                  = 15       // 建筑最大视野
	SysStrengthDuration               = 27       //英雄体力恢复周期
	SysMaxAngerValue                             //英雄最大怒气值
	SysBaseMapMoveFactor              = 29       //大地速度系数
	SysBaseMaxSoldierNum              = 30       //英雄基础最大上限带兵数量(英雄等级*这个系数)
	SysBaseSoldierNum                 = 31       //英雄默认带兵数量
	SysInitHeroPosZ                              //战斗中英雄距离中心点站位的距离
	SysFortressDispatchRate                      //要塞调度时间比例
	SysBuildFortressPowerAdd                     //建造要塞增加玩家势力值
	SysOccupyFortressPowerAdd                    //占领要塞增加玩家势力值
	SysMaxFortressNum                            //自建要塞最大数量
	SysQuitFortressTime               = 37       //放弃要塞需要时间(min)
	SysDismantleFortressTime          = 38       //拆除要塞需要时间(min)
	SysQuitLandTime                   = 39       //放弃地块需要时间(min)
	SysHurtSoldierA                              //暂时没用,待删除//伤兵率公式A
	SysHurtSoldierB                              //伤兵率公式（伤兵率=A-B*T/1000）中的参数B(T为时间单位毫秒)
	SysskillTimeInterval                         //技能时间间隔(每次使用技能后的时间间隔)
	SysDefendReadyTime                           //坚守准备持续时间
	SysDefendCityTime                            //坚守持续时间
	SysDefendCDTime                              //坚守冷却持续时间
	SysMarchDistance                  = 46       //行军最大距离
	SysMaxStrength                               //英雄最大体力值
	SysLostLandPercent                = 48       //沦陷丢失土地百分比
	SysPlunderPercent                 = 49       //掠夺当前资源百分比
	SysRebelPercent                   = 50       //反叛缴纳资源当前常量倍数
	SysPlayerWorldExp                 = 51       //玩家双方守军根据配置开关获得经验
	SysSpeedFinish                    = 61       //立即完成消耗钻石（优先绑定钻石）
	SysDarkNoticeTime                            //永夜倒计时提醒时间
	SysInitBarrackSoldierID           = 63       //初始的步兵兵种ID
	SysInitStableSoldierID            = 64       //初始的骑兵兵种ID
	SysInitShootSoldierID             = 65       //初始的弓兵兵种ID
	SysInitCarSoldierID               = 66       //初始的车兵兵种ID
	SysNoWarTime                      = 67       //免战时间(s)
	SysFirstTask                      = 68       //第一个主线任务
	SysFirstTrend                     = 69       //第一个天下大势任务
	SysDefaultWarnTime                = 75       //预警时间
	SysBtleBigSkillShowTime                      //战斗英雄大招定帧(展示)时间
	SysBtleBigSkillEffectID                      //战斗英雄大招特效id
	SysCostStrength                   = 78       //每次出征消耗体力
	SysreturnStrength                 = 79       //每次撤退返还体力
	SysInitArmyHeroes                 = 80       //初始部队中的英雄ID
	SysInitHeroID                     = 81       //初始获取的英雄ID
	SysInvestigateTime                = 82       //侦察时间
	SysBirthTime                      = 83       // 生育时间
	SysBoredomTime                    = 84       // 厌倦时间
	SysGrowTime                       = 85       // 成长时间
	SysPopulation                     = 86       // 人口
	SoldierJump                                  //士兵被骑兵冲击击飞参数(1速度2加速度)
	InitArmyHeroNum                   = 88       //部队初始英雄个数
	SysMaxArmyHeroNum                 = 89       //部队最大英雄个数
	KnightTurnWaitTime                           //骑兵转身后的等待时间
	CityBuildReduce                   = 91       //郡城建筑减少时间比例
	BtleArcherAtkCD                              //弓兵表现攻击CD(下限和上限)
	BtleVehicularAtkCD                           //炮车表现攻击CD(下限和上限)
	BtleKnightFirstRunTime                       //骑兵第一次冲锋前等待时间MS(固定等待时间,士兵间相差时间)
	CreateAllyGold                    = 95       //创建同盟金币
	InitFortressReferID                          //初始玩家要塞referID
	InitTransferReferID               = 97       //初始玩家传送阵referID
	RecruitFreeTime                              //免费招募时间
	IncreaseRecruitCost               = 99       //付费招募递增消耗钻石
	InitCityReferID                   = 100      //初始玩家城区ReferID
	InitCityWallReferID               = 101      //初始玩家城墙ReferID
	FortressMaxDispatchNum            = 102      //中立要塞和军营最大调度数量
	HeroHurtTime                      = 103      //英雄重伤时间
	SoldierReturnCost                 = 104      //征兵返回资源百分比
	MapOffsetX                        = 105      //视野X轴偏移格子数
	MapOffsetY                        = 106      //视野Y轴偏移格子数
	CollectTaxCost                    = 107      //强征消耗钻石
	TransferWallReferID               = 108      //玩家传送阵城墙referID
	MaxMarkLandNum                    = 109      // 标记地块上限
	GuideQuickConscription            = 110      // 引导快速征兵
	MaxEquipCount                     = 111      // 英雄最大装备个数
	MaxHeroGrade                      = 112      // 英雄最大阶级
	MaxHeroStar                       = 113      //英雄最大星级
	HeroInitLevel                     = 114      //英雄初始等级
	HeroInitGrade                     = 115      //	英雄初始阶级
	HeroInitStrength                  = 116      //英雄初始体力值
	HeroInitSkillLevel                = 117      //英雄初始技能等级
	ConscriptSpeed                    = 118      // 征兵加速每个多少秒
	DecreeDuration                    = 119      // 征兵加速产出周期/政令产出周期
	ConscriptSpeedCount               = 120      // 征兵加速每次产出n个
	DecreeAttrCount                   = 121      // 政令每次产出n个
	InitConscriptSpeedLimit           = 122      // 初始征兵加速上限
	InitConscriptSpeed                = 123      // 初始征兵加速
	InitResidentCount                 = 125      // 初始玩家居民个数
	ResidentAdTime                    = 129      // 居民广告招募冷却时间
	ResidentAdCount                   = 134      // 居民免费招募次数
	FastMarch                         = 135      // 新手引导快速行军
	FastMarchSpeed                    = 136      // 新手引导快速行军速度倍数
	SysAllyRequestNum                 = 145      // 同盟申请列表上限
	FallNoWarTime                     = 146      // 沦陷玩家的免战时间
	FallSaveDurableRate               = 148      // 沦陷玩家恢复耐久比例
	FailCompenseNum                   = 149      // 失败补偿：玩家兵力大于一定数量
	FailCompenseList                  = 150      // 失败补偿：中立资源地范围
	FailCompenseMax                   = 151      // 失败补偿：最大补偿次数
	LookoutTowerLandLv                = 152      // 创建瞭望塔土地等级>=4
	CellarLandLv                      = 153      // 创建地窖土地等级
	WorkShopLandLv                    = 154      // 创建工坊土地等级
	TrainCampLandLv                   = 155      // 创建训练营土地等级
	GoldmineLandLv                    = 156      // 创建金矿土地等级
	GhostHallEffect                   = 157      // 英灵殿效果时间（秒）
	FreeTaxTime                       = 159      // 免费税收冷却时间
	GateOfHellTime                    = 160      // 地狱之门的效果时间(秒)
	GateOfHellNum                     = 161      // 地狱之门挑战次数
	GateOfHellRefresh                 = 162      // 地狱之门周六日刷出时间
	ConscriptionGold                  = 165      // 征兵1个兵消耗的金币
	AllyBreakUpTime                   = 166      // 解散同盟CD时间（秒）
	AllyBreakUpMember                 = 167      // 解散同盟人数
	AllyContributeRate                = 169      // 同盟捐献兑换比例
	AssaultPlunderDistance            = 170      // 突击掠夺距离
	AssaultPlunderTime                = 171      // 突击掠夺的保护时间
	AllyCooperationTimes              = 175      // 同盟协助次数
	AllyCooperationRefresh            = 176      // 同盟协助刷出时间
	ThiefIncidentID                   = 177      // 小偷事件id
	ReNameData                        = 178      // 改名数据,36小时内可修改1次，消耗龙晶500
	AllyAutoInviteNum                 = 181      // 同盟自动邀请人数
	MailTitleLen                      = 182      // 邮件标题长度
	MailContentLen                    = 183      // 邮件内容长度
	FarmNeedDecree                    = 184      // 屯田需要的政令
	ArmyNameLen                       = 186      // 军团名字长度
	UnitedBackAtOnce                  = 188      // 联军立即撤退花费
	MaxTransferName                   = 190      // 最长传送阵名字
	MaxFortressName                   = 191      // 最长要塞名字（其他建筑也用这个)
	SysInitSectionArmyHeroes          = 192      // 初始战役部队中的英雄ID
	SysAllyAbdicationTime             = 193      // 禅让盟主CD时间（秒）
	MaxSectionPower                   = 197      // 战役能量最大值
	AllyChestRefreshTime              = 198      // 同盟宝箱刷出的时间
	AllyReleaseTaskCD                 = 199      // 同盟发布任务次数的恢复时间
	JoinAllyCD                        = 200      // 重新加入同盟CD
	AllyCooperationSpeed              = 201      // 同盟小偷协作速度倍数
	AllyInviteExpire                  = 202      // 同盟邀请有效时间
	AllyTreasureHuntTime              = 203      // 同盟寻宝活动开启期间
	AllyTreasureHuntActionRecover     = 204      // 同盟寻宝行动恢复的时间和次数(1小时1次)
	AllyTreasureHuntExploreReduce     = 205      // 同盟寻宝增加同盟成员探索减少X分钟
	AllyTreasureHuntOccupyExtra       = 206      // 同盟寻宝占领额外多的百分比
	AllyTreasureHuntActionMaxTimes    = 207      // 同盟寻宝行动上限最多次数
	SysInitSectionTask                = 208      // 第一个关卡活动任务
	AllyTreasureHuntActionBubbleCount = 209      // 同盟寻宝占地显示泡泡的个数
	AllyTreasureHuntStartPointLand    = 211      // 同盟寻宝地块id
	MapExploreDailyTimes              = 215      // 玩家每天探索的次数
	MapExploreDailyCD                 = 216      // 玩家每天探索的cd(秒)
	MapExploreArmySpeed               = 217      // 探索行军速度
	AllyQADayTimes                    = 218      // 同盟问答每日答题次数
	GuideSummonGoods                  = 219      // 引导五连抽物品列表
	GuideMapBattle                    = 220      // 新手引导主城城皮地块战斗引导ID(引导^step^armyID)
	ExploreEventCityRadius            = 221      // 中立城半径10格
	ExploreEventStageRadius           = 222      // 关卡半径5格
	DispelRadius                      = 223      // 驱散半径
	DispelArmySpeed                   = 224      // 驱散部队速度
	DispelCD                          = 225      // 驱散cd（秒）
	DispelSpend                       = 226      // 驱散一圈消耗的时间（秒），和驱散半径相关
	DispelDailyTimes                  = 227      // 驱散每天次数
	GeneralFragmentID                 = 228      // 通用的英灵碎片
	HeroResetTime                     = 229      // 英雄重置时间冷却时间(秒)
	VillageResidentCount              = 230      // 每次村庄产出居民的个数
	AllyTreasureHuntFirstTimes        = 231      // 同盟寻宝首次体验次数
	ExploreEventArmyID                = 232      // 探索事件的特殊部队id
	VillageReferID                    = 233      // 村庄的ReferID
	EnterCityNeedResident             = 234      // 进入内城需要的获得过的居民数量
	SystemAllyCreateNum               = 237      // 以国家为单位，每创建多少真实玩家（倍数），系统就创建一个同盟
	SystemAllyLeaderMemberNum         = 238      // 同盟人数到达n时，选举出势力值最高的作为盟主
	GuideLand                         = 239      // 新手引导地块不掉居民
	DayRecruitFreeCount               = 240      // 每天玩家能免费招募的居民个数(村庄每天的产出)
	ResidentSeekDuration              = 241      // 居民投奔间隔时间(村庄产出间隔时间)
	LandResidentDropRange             = 242      // 地块掉落居民的范围(与玩家出生点的距离)
	VillageRefreshRange               = 243      // 村庄的刷新范围(玩家出生点几格范围内)
)

const (
	OnlineGiftRewardTimes           = iota + 1 // 第几次领取礼包翻倍数
	OnlineGiftCountDownSection                 // 在线礼包倒计时区间（单位分）
	CommunityConcernInstagram                  // Instagram社区关注奖励
	CommunityConcernTwitter                    // Twitter社区关注奖励
	CommunityConcernFacebook                   // Facebook社区关注奖励
	FirstRecharge                              // 首充奖励
	MonthCardDiamond                           // 月卡每日领取
	FirstJoinAlly                              // 首次入盟
	BlackMarketResetCost                       // 黑市重置消耗钻石数
	BlackMarketResetTime                       // 黑市重置时间点
	BindMail                                   // 绑定邮件
	BindFacebook                               // 绑定facebook
	BindGooglePlay                             // 绑定Google Play
	OpenServerMailReward            = 17       // 开服邮件奖励
	OnlineGiftRewardDay             = 19       // 在线礼包领取天数
	OnlineGiftDayTimes              = 20       // 在线礼包一天领取次数
	SpaceRiftAward                  = 21       // 空间裂缝奖励（24小时获得英雄碎片）
	SurveyReward                    = 22       // 问卷调查奖励
	FirstDownLoadSecondPackageAward = 23       // 第一次下载资源包奖励
	HappyWeekFinishAward            = 24       // 七天乐活动终极大奖
	OldPlayerBindingAward           = 25       // 老玩家绑定奖励
)
