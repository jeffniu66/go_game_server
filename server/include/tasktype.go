package include

//type Task struct {
//	TaskId     int32
//	Type       int32
//	SubType    int32
//	TargetList []*Target
//	State      int32
//	Update     int32 // 更新标志，1为需要回写数据库
//	CreateTime int32 //
//	ExpireTime int32 // 0为无限时间
//}

type TaskProgressMsgData struct {
	SubType      int32
	AllyId       int32
	UserId       int32
	ProgressMap  map[int32]int32 // 目标id：目标进度
	ProgressType int32           // 添加进度类型，0是增加，1是当前值
}

// 任务完成状态
const (
	TaskStateNotFinish   = 0 // 未完成
	TaskStateFinishNoGet = 1 // 已完成未领取
	TaskStateFinishGet   = 2 // 已完成已领取
)

// 任务目标
type Target struct {
	TargetId int32 // 目标ID
	//TargetSum int32 // 目标总数 废弃,读配置表取最大值
	TargetCur int32 // 目标当前数
}

// 任务类型
const (
	TaskTypeMain           = iota + 1 // 主线任务
	TaskTypeDaily                     // 日常任务
	TaskTypeAlly                      // 工会任务
	TaskTypeLord                      // 领主任务
	TaskTypeSign                      // 签到任务
	TaskTypePowerRank                 // 势力活动任务
	TaskTypeWorldTrend                // 天下大势任务
	TaskTypeAllyDaily                 // 同盟个人日常任务
	TaskTypeAllyLimit                 // 同盟个人限定任务
	TaskTypeAllyGroup                 // 同盟集体任务
	TaskTypeAllyPowerRank             // 同盟势力活动任务
	TaskTypeCityBuild                 // 主城建设任务
	TaskTypeSection                   // 关卡活动任务
	TaskTypeAllyAttackCity            // 同盟攻城活动任务
	TaskTypeHappyWeek                 // 七天乐活动任务
	TaskTypeMoreLand                  // 扩张领土活动任务
	TaskTypeMoreAllies                // 壮大同盟活动任务
)

// 任务子类(9,10,14,16,21,22,23,34,35,37,58,59废弃了)
const (
	TaskSubTypePlayer                 = 1    // 角色属性(角色的某项属性，达到某个值)
	TaskSubTypeBuild                  = 2    // 城建进度(某个建筑，达到某个等级)
	TaskSubTypeSoldierT               = 3    // 4种兵种中有N种达到L级
	TaskSubTypeSection                = 4    // 战役通关M关卡
	TaskSubTypeAllyCity               = 5    // 同盟占领L级城
	TaskSubTypeLand                   = 6    // 打N个L级地
	TaskSubTypeDelLand                = 7    // 成功放弃1块土地
	TaskSubTypeEquip                  = 8    // 进行一次穿装备行为
	TaskSubTypeSkill                  = 9    // 进行一次升级技能行为
	TaskSubTypeHeroStar               = 10   // N个英雄升星到S星
	TaskSubTypeHeroGrade              = 11   // N个英雄进阶到G阶
	TaskSubTypeSoldierN               = 12   // 任意一个部队当前兵力达到N个兵
	TaskSubTypeLord                   = 13   // 领主达到L级
	TaskSubTypeFarm                   = 14   // 进行一次屯田行为(加资源后)
	TaskSubTypeAlly                   = 15   // 加入同盟
	TaskSubTypeHold                   = 16   // 占领一次其他人的领地
	TaskSubTypeTax                    = 18   // 征收一次税收
	TaskSubTypeYield                  = 19   // 产量达到N
	TaskSubTypeScience                = 20   // 科技 S科技研究达到L级
	TaskSubTypeReserveConscript       = 21   // 使用一次预备征兵
	TaskSubTypeNormalConscript        = 22   // 进行一次征兵,只需要点击征兵就算一次
	TaskSubTypeChangeName             = 23   // 修改自己的名字
	TaskSubTypeAllyInvite             = 24   // 邀请同盟人数N
	TaskSubTypeRechargeNum            = 25   // 充值金额N
	TaskSubTypeCostResource           = 26   // 资源消耗{M,N}
	TaskSubTypeSign                   = 27   // 完成签到N次
	TaskSubTypeTalkTimes              = 28   // 完成发言N次
	TaskSubTypeFinishCooperationTimes = 29   // 完成同盟帮助N次
	TaskSubTypeAskCooperationTimes    = 30   // 完成同盟求助N次
	TaskSubTypeAllyContributeSum      = 31   // 完成同盟捐献N
	TaskSubTypeUpgradeCityTimes       = 32   // 完成城建升级N次
	TaskSubTypeSectionTimes           = 33   // 完成攻打战役N次
	TaskSubTypeDecreeTimes            = 34   // 完成政令消耗N次
	TaskSubTypeConscriptSpeedTimes    = 35   // 完成征兵加速消耗N次
	TaskSubTypeUpgradeHeroTimes       = 36   // 完成英雄升级N次
	TaskSubTypeSummonHeroTimes        = 37   // 完成召唤N次
	TaskSubTypeWarMailTimes           = 38   // 查看战报N次
	TaskSubTypeWarVideoTimes          = 39   // 回放战报N次
	TaskSubTypeResidentNum            = 40   // 招募N个居民
	TaskSubTypeScienceTimes           = 41   // 研究科技N次
	TaskSubTypeInvestigationTimes     = 42   // 侦查土地N次
	TaskSubTypeAllyMemberNum          = 43   // 同盟成员N个
	TaskSubTypeKillSoldierNum         = 44   // 杀敌个数N
	TaskSubTypeShareTimes             = 45   // 分享N次
	TaskSubTypeDiamondLottery         = 46   // 钻石召唤N次
	TaskSubTypeGoldLottery            = 47   // 金币召唤N次
	TaskSubTypeDestroyGoblin          = 48   // 消灭哥布林N次
	TaskSubTypeFarmSum                = 49   // 屯田获得的资源总量N
	TaskSubTypeCostEveryResource      = 50   // 消耗任意材料总量N
	TaskSubTypeAllyLevel              = 51   // 同盟达到N级
	TaskSubTypeUpLevelSkill           = 52   // 接任务后，英雄升级技能N次
	TaskSubTypeOccupyLand             = 53   // 接任务后，占领N个大于等于L级地
	TaskSubTypeAttackCity             = 54   // 同盟攻占了多少个城池
	TaskSubTypeAttackLand             = 55   // 出征一块土地(点出征就算)
	TaskSubTypeSeeLand                = 56   // 查看我的领土
	TaskSubTypeHeroLevel              = 57   // N个英雄升到L级
	TaskSubTypeCreateBuild            = 58   // 行为-建造N个要塞(金矿,工坊等)
	TaskSubTypeNeutralFortress        = 59   // 占领N个中立要塞(玩家身上有多少个中立要塞)
	TaskSubTypeOpenFAQ                = 60   // 打开FAQ
	TaskSubTypeFog                    = 61   // 解开N块迷雾
	TaskSubTypeBlackMarket            = 62   // 黑市购买N个道具
	TaskSubTypeOwnHero                = 63   // 拥有N个英雄
	TaskSubTypeRecoverLost            = 64   // 收复失地
	TaskSubTypeBind                   = 65   // 绑定帐号
	TaskSubTypeAchievement            = 1000 // 成就类型的任务
)

// 添加进度类型
const (
	TaskProgressTypeAdd    = iota // 增加
	TaskProgressTypeAssign        // 当前值
	TaskProgressTypeMax           // 保留最大值
	TaskProgressTypeGTEAdd        // 大于等于增加
)
