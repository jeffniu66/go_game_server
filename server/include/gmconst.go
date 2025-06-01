package include

const (
	GMAllRoom        = "all_room"         // 查找所有正在运行的房间 eg. all_room 返回 数量，房间ids
	GMGetRoom        = "get_room"         // 查找某个房间的信息 eg.get_room^1983
	GMAddItem        = "gmAddItem"        // 加道具 格式: gmAddItem^itemId,itemNum|itemId,itemNum
	GMAddSkill       = "gmAddSkill"       // 加技能 格式： gmAddSkill^skillId
	GMGetOnlineNum   = "gmGetOnlineNum"   // 获取在线用户数
	GMGetRegisterNum = "gmGetRegisterNum" // 获取在线用户数
	GMUptRank        = "uptRank"          // 更新排行榜
)
