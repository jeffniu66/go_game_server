package include

import "go_game_server/proto3"

type Task struct {
	UserId           int32
	PointInfo        string                   // 完成的任务: id,id,id
	AssignTaskPoint  string                   // 分配的任务位置: 1,2,3,4
	Skill            *proto3.Item             // 选择获取到的技能
	TempDropItemsMap map[int32][]*proto3.Item // 临时掉落道具 key: 道具类型
	LuckCardMap      map[int32]int32          // 幸运卡使用数据 key: 幸运卡道具id  value: 幸运卡使用数量
	TotalGetSkills   []int32                  // 所有获得的技能id
	SkillTabNumMap   map[int32]int32          // 技能选项卡数量 key: UserId value: 数量
}
