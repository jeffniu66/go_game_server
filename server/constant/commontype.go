package constant

// 完成状态
const (
	StateNotFinish   = iota // 未完成
	StateFinishNoGet        // 已完成未领取
	StateFinishGet          // 已完成已领取
)

const (
	ResultSuccess = 0
	ResultFail
)

const (
	WeekDay    = 7
	DayHour    = 24
	MinSecond  = 60
	HourSecond = 3600
	DaySecond  = 86400
	WeekSecond = 604800
)

const (
	DropGroupProbability  = 1000 // 掉落组概率，<1000表示有不掉落概率 >1000表示必掉落
	TaskPoints            = 6    // 任务位置点个数
	TaskTypeNormal        = 1    // 普通任务
	TaskTypeUrgency       = 2    // 紧急任务
	UrgencyTaskTypeNoDead = 0    // 非致命
	UrgencyTaskTypeDead   = 1    // 致命

	ItemTabSecond       = 101 // 2号道具选项卡
	ItemTabThird        = 102 // 3号道具选项卡
	KillPeopleNum       = 2   // 狼人杀人数掉落池
	WolfManCdAllowRange = 2   // 狼人杀人CD允许的范围
	MaxNameLen          = 32  // 最大名字长度
	AIMaxNum            = 9   // 最大AI数量
)

const (
	ActTypeSingle = 1 // 玩家行为单次
	ActTypeMulti  = 0 // 玩家行为多次
)
