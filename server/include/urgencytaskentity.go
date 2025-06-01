package include

type UrgencyTask struct {
	TaskPoints      []int32           // 完成的任务点
	IngPoint        int32             // 正在进行的紧急任务点 0表示没有
	IngUserId       int32             // 当前发起的任务玩家id
	LastTriggerTime int64             // 上次触发紧急任务时间
	TriggerNumMap   map[int32][]int32 // 狼人触发任务数量 key: userId value: 触发的紧急任务数量
}
