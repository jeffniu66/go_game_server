package game

import (
	"go_game_server/server/constant"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
)

var GMSkillId int32 // GM设置的技能id

func GetDropGroupConfigByDropId(dropId int32) *tableconfig.DropGroupConfig {
	dropGroupList := tableconfig.DropGroupConfigs.GetDropGroupMap(int(dropId))
	if len(dropGroupList) == 0 {
		logger.Log.Errorln("GetDropGroupConfigByDropId dropGroupList is nil")
		return nil
	}
	return GetDropGroupConfig(dropGroupList)
}

func GetDropGroupConfig(dropGroupList []tableconfig.DropGroupConfig) *tableconfig.DropGroupConfig {
	// TODO GM 待删除
	if GMSkillId > 0 {
		return &tableconfig.DropGroupConfig{Id: 9999, DropId: 1, ItemType: 2, ItemId: int(GMSkillId), Num: 1, Probability: 1000}
	}
	var sum int
	var pros []uint32
	for _, v := range dropGroupList {
		sum = sum + v.Probability
		pros = append(pros, uint32(v.Probability))
	}
	if sum < constant.DropGroupProbability {
		pros = append(pros, uint32(constant.DropGroupProbability-sum))
		index := util.RandGroup(pros...)
		if index == len(pros)-1 {
			return nil
		}
		dropGroupConfig := dropGroupList[index]
		return &dropGroupConfig
	} else {
		index := util.RandGroup(pros...)
		dropGroupConfig := dropGroupList[index]
		return &dropGroupConfig
	}
}

// 使用过幸运卡道具 增加技能掉落概率
func GetDropGroupConfWithBuf(task *Task, dropGroupList []tableconfig.DropGroupConfig) *tableconfig.DropGroupConfig {
	if GMSkillId > 0 {
		return &tableconfig.DropGroupConfig{Id: 9999, DropId: 1, ItemType: 2, ItemId: int(GMSkillId), Num: 1, Probability: 1000}
	}
	var sum int
	var pros []uint32
	var sumProb int32
	if task != nil && task.LuckCardMap != nil {
		for k, v := range task.LuckCardMap {
			itemConfig := tableconfig.ItemConfigs.GetItemConfigById(k)
			if itemConfig == nil {
				continue
			}
			sumProb += int32(int(util.ToInt(itemConfig.Addnum)) * int(v))
		}
	}
	for _, v := range dropGroupList {
		prob := v.Probability + v.Probability*int(sumProb)
		sum = sum + (prob)
		pros = append(pros, uint32(prob))
	}
	if sum < constant.DropGroupProbability {
		pros = append(pros, uint32(constant.DropGroupProbability-sum))
		index := util.RandGroup(pros...)
		if index == len(pros)-1 {
			return nil
		}
		dropGroupConfig := dropGroupList[index]
		return &dropGroupConfig
	} else {
		index := util.RandGroup(pros...)
		dropGroupConfig := dropGroupList[index]
		return &dropGroupConfig
	}
}
