package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"strconv"
	"strings"
)

var FinishTaskDropItemNum int32 = 1 // 掉落道具个数 默认为1

func saveItem(item *include.Item) {
	if item == nil {
		return
	}
	item.Update = include.Update
	db.ItemMap[item.UserId][item.ItemId] = item
}

// 取道具
func GetItem(userId int32, itemId int32) *include.Item {
	items, ok := db.ItemMap[userId]
	if ok { // 存在玩家的道具集合
		if item, ok2 := items[itemId]; ok2 {
			return item
		}
		return nil // 道具集合存在 但是道具不存在 则数据库里也不存在
	}
	// 不存在道具集合 从数据库加载
	itemsDB := db.GetPlayerItems(userId)
	if len(itemsDB) == 0 {
		return nil
	}
	items = make(map[int32]*include.Item)
	for _, item := range itemsDB {
		items[item.ItemId] = item
	}
	db.ItemMap[userId] = items
	if item, ok := items[itemId]; ok {
		return item
	}
	return nil
}

// 获取玩家皮肤碎片
func GetSkinPiece(userId int32) (items []*proto3.SkinPiece) {
	ret := make([]*proto3.SkinPiece, 0)
	itemList := GetItems(userId)

	for _, v := range itemList {
		itconfigs := tableconfig.ItemConfigs.GetItemConfigById(v.Id)
		if itconfigs == nil {
			logger.Log.Errorf("this skin:%v ItemConfigs isn't setting please fix item.excel", v.Id)
			return
		}
		if itconfigs.Type == int32(proto3.ItemTypeEnum_skin_item) {
			tmp := &proto3.SkinPiece{Id: v.Id, PieceNum: v.Num}
			ret = append(ret, tmp)
		}
	}
	if len(ret) <= 0 {
		return nil
	}
	return ret
}

// 获取玩家所有道具-登录返回
func GetItems(userId int32) (items []*proto3.Item) {
	itemsMap, ok := db.ItemMap[userId]
	if !ok {
		// 从数据库读取玩家道具
		itemsDB := db.GetPlayerItems(userId)
		if len(itemsDB) == 0 {
			return items
		}
		// 数据库存在 缓存一下
		itemsMap = make(map[int32]*include.Item)
		for _, v := range itemsDB {
			item := &proto3.Item{Id: v.ItemId, Num: v.Num}
			items = append(items, item)

			itemsMap[v.ItemId] = v
		}
		db.ItemMap[userId] = itemsMap
	} else {
		for _, v := range itemsMap {
			item := &proto3.Item{Id: v.ItemId, Num: v.Num}
			items = append(items, item)
		}
	}
	return items
}

func (t *Task) addSkills(dropGroupConfigs []*tableconfig.DropGroupConfig) {
	if len(dropGroupConfigs) == 0 {
		return
	}
	tempDropItemsMap := t.TempDropItemsMap
	if tempDropItemsMap == nil {
		tempDropItemsMap = make(map[int32][]*proto3.Item)
	}
	for _, v := range dropGroupConfigs {
		item := &proto3.Item{Id: int32(v.ItemId), Num: int32(v.Num)}
		if _, ok := tempDropItemsMap[int32(v.ItemType)]; !ok {
			tempDropItemsMap[int32(v.ItemType)] = []*proto3.Item{item}
		} else {
			tempDropItemsMap[int32(v.ItemType)] = append(tempDropItemsMap[int32(v.ItemType)], item)
		}
	}
	t.TempDropItemsMap = tempDropItemsMap
}

// 获取可以掉落道具的数量（任务或狼人杀人）
func (r *Room) getCanDropItemNum(userId int32) int32 {
	taskMap := r.taskMap
	if taskMap == nil {
		return FinishTaskDropItemNum
	}
	task := taskMap[userId]
	if task == nil {
		return FinishTaskDropItemNum
	}
	if task.SkillTabNumMap == nil {
		return FinishTaskDropItemNum
	}
	return FinishTaskDropItemNum + task.SkillTabNumMap[userId]
}

func (r *Room) dropItem(skillPoolType int32, userId int32) map[int32][]*proto3.Item {
	dropItemMap := make(map[int32][]*proto3.Item)
	for i := 0; i < int(r.getCanDropItemNum(userId)); i++ {
		logger.Log.Infof("可以掉落的道具数量: %d", r.getCanDropItemNum(userId))
		var dropGroupConf *tableconfig.DropGroupConfig
		if skillPoolType == constant.NormalManSkillPool {
			dropGroupConf = r.doTaskDropSkill(userId)
		} else {
			dropGroupConf = r.killPeopleDropSkill(userId)
		}
		if dropGroupConf == nil {
			logger.Log.Info("FinishTask dropItem dropGroupConfig is nil")
			continue
		}

		itemType := int32(dropGroupConf.ItemType)
		itemId := int32(dropGroupConf.ItemId)
		itemNum := int32(dropGroupConf.Num)

		if _, ok := dropItemMap[itemType]; !ok {
			items := []*proto3.Item{{Id: itemId, Num: itemNum}}
			dropItemMap[itemType] = items
		} else {
			item := &proto3.Item{Id: itemId, Num: itemNum}
			dropItemMap[itemType] = append(dropItemMap[itemType], item)
		}
		r.addDropItem(userId, dropGroupConf)
	}
	return dropItemMap
}

// 狼人杀人掉落道具
func (r *Room) KillPeopleDropItem(userId int32) {
	if !r.roomInfo.IsWolfMan(userId) {
		return
	}
	dropItemMap := r.dropItem(constant.WolfManSkillPool, userId)
	if len(dropItemMap) != 0 {
		player := global.GloInstance.GetPlayer(userId)
		if player == nil {
			return
		}
		p := player.(*Player)
		p.Pid.Cast("dropItemResp", dropItemMap)
	}
}

// 掉出的道具如果存在技能池内，则减掉池内的数量
func (r *Room) reduceSkillPoolSkillNum(skillPoolType int32, dropConf *tableconfig.DropGroupConfig) {
	if dropConf == nil {
		return
	}
	skillPoolMap := make(map[int32]int32)
	if skillPoolType == constant.NormalManSkillPool {
		skillPoolMap = r.normalManSkillPoolMap
	} else if skillPoolType == constant.WolfManSkillPool {
		skillPoolMap = r.wolfManSkillPoolMap
	}
	skillNum, ok := skillPoolMap[int32(dropConf.ItemId)]
	if !ok {
		return
	}
	if skillNum <= 0 {
		return
	}
	leftNum := skillNum - 1
	if leftNum < 0 {
		leftNum = 0
	}
	skillPoolMap[int32(dropConf.ItemId)] = leftNum
	if skillPoolType == constant.NormalManSkillPool {
		r.normalManSkillPoolMap = skillPoolMap
	} else if skillPoolType == constant.WolfManSkillPool {
		r.wolfManSkillPoolMap = skillPoolMap
	}
}

// 做任务掉技能
func (r *Room) doTaskDropSkill(userId int32) *tableconfig.DropGroupConfig {
	// 根据完成任务数 选择不同掉落池
	task, ok := r.taskMap[userId]
	if !ok {
		return nil
	}
	var dropGroupList []tableconfig.DropGroupConfig
	if r.roomInfo.CheckPlayerGameStatus(userId, proto3.PlayerGameStatus_killed) { // 死亡
		deathDropPool := tableconfig.ConstsConfigs.GetValueById(constant.DeathDropPool)
		dropId, _ := strconv.Atoi(deathDropPool)
		dropGroupList = tableconfig.DropGroupConfigs.GetDropGroupMap(dropId)
		if len(dropGroupList) == 0 {
			logger.Log.Errorln("dropGroupList is null")
			return nil
		}
		dropGroupList = r.excludeComPoolDropConf(constant.DeathDropSumPool, dropGroupList)
	} else {
		finishTaskNum := GetFinishTaskNum(task.PointInfo)
		dropGroupId := finishTaskNum + 1
		dropGroupList = tableconfig.DropGroupConfigs.GetDropGroupMap(int(dropGroupId))
		if len(dropGroupList) == 0 {
			logger.Log.Errorln("dropGroupList is null")
			return nil
		}
		dropGroupList = r.excludeComPoolDropConf(constant.NormalManSkillPool, dropGroupList)
	}
	return GetDropGroupConfWithBuf(task, dropGroupList)
}

// 狼人杀人掉技能
func (r *Room) killPeopleDropSkill(userId int32) *tableconfig.DropGroupConfig {
	if r.roomInfo == nil {
		return nil
	}
	var dropGroupId int32
	killNum := r.roomInfo.GetKillNum(userId)
	if killNum <= constant.KillPeopleNum {
		dropGroupId = util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.WolfManKillOneLvPool))
	} else {
		dropGroupId = util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.WolfManKillTwoLvPool))
	}
	dropGroups := tableconfig.DropGroupConfigs.GetDropGroupMap(int(dropGroupId))
	if len(dropGroups) == 0 {
		logger.Log.Errorln("KillPeopleRandDropItem dropGroupList is null")
		return nil
	}

	dropGroups = r.excludeComPoolDropConf(constant.WolfManSkillPool, dropGroups)

	// 删除公共池内技能
	dropConf := GetDropGroupConfWithBuf(r.taskMap[userId], dropGroups)
	if r.roomInfo.IsWolfMan(userId) {
		r.reduceSkillPoolSkillNum(constant.WolfManSkillPool, dropConf)
	} else {
		r.reduceSkillPoolSkillNum(constant.NormalManSkillPool, dropConf)
	}
	return dropConf
}

// 排除公共技能池内数量为0的掉落项
func (r *Room) excludeComPoolDropConf(skillPoolType int32, dropGroupList []tableconfig.DropGroupConfig) []tableconfig.DropGroupConfig {
	skillPoolMap := make(map[int32]int32)
	if skillPoolType == constant.NormalManSkillPool { // 正常玩家
		skillPoolMap = r.normalManSkillPoolMap
	} else if skillPoolType == constant.WolfManSkillPool { // 狼人
		skillPoolMap = r.wolfManSkillPoolMap
	} else if skillPoolType == constant.DeathDropSumPool { // 死亡
		skillPoolMap = r.deathSkillPoolMap
	}
	dropGroupListTemp := make([]tableconfig.DropGroupConfig, 0)
	for _, v := range dropGroupList {
		skillNum, ok := skillPoolMap[int32(v.ItemId)]
		if !ok {
			continue
		}
		if skillNum > 0 {
			dropGroupListTemp = append(dropGroupListTemp, v)
		}
	}
	return dropGroupListTemp
}

func (r *Room) addDropItem(userId int32, dropGroupConfig *tableconfig.DropGroupConfig) {
	// 存储掉落物品 如果是道具，需要持久化 如果是技能，则不用持久化
	if task, ok := r.taskMap[userId]; ok {
		if dropGroupConfig.ItemType == int(proto3.DropTypeEnum_item) {
			player := global.GloInstance.GetPlayer(userId)
			if player != nil {
				p := player.(*Player)
				// 根据道具id取道具表
				itemConfig := tableconfig.ItemConfigs.GetItemConfigById(int32(dropGroupConfig.ItemId))
				if itemConfig == nil {
					p.ErrorResponse(proto3.ErrEnum_Error_Goods_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Goods_NotExists)])
					return
				}
				p.AddItems(util.ToStr(int32(dropGroupConfig.ItemId)) + "," + util.ToStr(int32(dropGroupConfig.Num)))
			}
		} else if dropGroupConfig.ItemType == int(proto3.DropTypeEnum_skill) {
			task.addSkills([]*tableconfig.DropGroupConfig{dropGroupConfig})
			r.taskMap[userId] = task
		}
	} else { // 狼人不存在
		task := &Task{UserId: userId}
		task.addSkills([]*tableconfig.DropGroupConfig{dropGroupConfig})
		r.taskMap[userId] = task
	}
}

// 加道具-公共方法
func (p *Player) AddItems(itemStr string) {
	if len(itemStr) == 0 {
		return
	}

	itemsResp := p.AddRewards(itemStr)

	if len(itemsResp) != 0 {
		cmd := proto3.ProtoCmd_CMD_ItemsResp
		pbData := &proto3.ItemResp{Items: itemsResp}
		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	}
}

func (p *Player) AddRewards(itemStr string) []*proto3.Item {
	if len(itemStr) == 0 {
		return nil
	}

	var skinRed bool // 皮肤推送红点
	var itemRed bool // 背包推送红点

	itemsResp := make([]*proto3.Item, 0)
	items := strings.Split(itemStr, "|")
	var itemData []*include.ItemData
	for itId, itNum := range mergeSameItem(items) {
		if itId == constant.ItemIdGold { // 货币
			p.Attr.Gold = p.Attr.Gold + itNum
			p.Attr.UserTitle.TotalGold += itNum
		} else if itId == constant.ItemIdArchive {
			p.Attr.ArchivePoint += itNum
			p.Attr.UserTitle.TotalArchive += itNum
			p.UpdateArchive()
		} else if itId == constant.ItemIdExp {
			p.Attr.Exp += itNum
			p.UpdateLevel()
		} else if itId == constant.ItemIdGemStone {
			p.Attr.GemStone += itNum
		} else { // 道具
			userId := p.Attr.UserID
			userItemMap := db.ItemMap[userId]
			var item *include.Item
			if userItemMap == nil {
				userItemMap = make(map[int32]*include.Item)
				item = &include.Item{Uid: db.ItemUidInstance.GetNewItemUid(), ItemId: itId, Num: itNum, UserId: userId}
			} else {
				item = GetItem(userId, itId)
				if item == nil {
					item = &include.Item{Uid: db.ItemUidInstance.GetNewItemUid(), ItemId: itId, Num: itNum, UserId: userId}
				} else {
					item.Num += itNum
				}
			}
			item.Update = include.Update
			userItemMap[itId] = item
			db.ItemMap[userId] = userItemMap

			// 红点判断
			itemConfig := tableconfig.ItemConfigs.GetItemConfigById(itId)
			if itemConfig.Type == constant.ItemTypeSkin {
				skinRed = true
			}
			if itemConfig.Type == constant.ItemTypeFixedRewardBox || itemConfig.Type == constant.ItemTypeRandomRewardBox {
				itemRed = true
			}
		}
		itemsResp = append(itemsResp, &proto3.Item{Id: itId, Num: itNum})

		itemData = append(itemData, &include.ItemData{UserId: p.Attr.UserID, ItemId: itId, Num: itNum})
	}

	if skinRed {
		skinRedData := p.GetSkinRedPoint()
		pbData := &proto3.RedPointResp{RedType: proto3.RedPointEnum_skin_red, RedData: skinRedData}
		p.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_RedPointResp, PbData: pbData})
	} else if itemRed {
		pbData := &proto3.RedPointResp{RedType: proto3.RedPointEnum_items_red, RedData: GetItemRedDot(p.Attr.UserID)}
		p.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_RedPointResp, PbData: pbData})
	}

	RecordItem(itemData)
	return itemsResp
}

// 合并相同itemId道具数量
func mergeSameItem(items []string) map[int32]int32 {
	itemMap := make(map[int32]int32)
	for _, v := range items {
		its := strings.Split(v, ",")
		itemId, _ := strconv.Atoi(its[0])
		itemNum, _ := strconv.Atoi(its[1])
		itemMap[int32(itemId)] += int32(itemNum)
	}
	return itemMap
}

//func (p *Player) AddItem(itemType int32, itemId int32, itemNum int32) {
//	if itemNum <= 0 {
//		return
//	}
//	if itemType == include.ItemTypeItem { // 道具
//		userId := p.Attr.UserID
//		items, ok := db.ItemMap[userId]
//
//		item := GetItem(userId, itemId)
//		if item == nil {
//			item = &include.Item{ItemId: itemId, Num: itemNum, UserId: userId}
//			if !ok { // 玩家道具集合不存在
//				items = make(map[int32]*include.Item)
//			}
//		} else {
//			item.Num += itemNum
//		}
//		items[itemId] = item
//		db.ItemMap[userId] = items
//		//(*db.Item)(item).SaveData()
//		//db.SaveItemsData(item)
//	} else if itemType == include.ItemTypeCurrency { // 货币
//		switch itemId {
//		case include.ItemIdGold:
//			p.Attr.Gold = p.Attr.Gold + itemNum
//		default:
//			break
//		}
//		p.SaveDBData()
//	}
//}

// 减道具-公共方法
func (p *Player) ReduceItem(itemType int32, itemId int32, itemNum int32) bool {
	if itemNum <= 0 {
		return false
	}
	if itemType == constant.ItemTypeCurrency { // 货币
		switch itemId {
		case constant.ItemIdGold:
			p.Attr.Gold -= itemNum
		default:
			break
		}
	} else {
		userId := p.Attr.UserID
		item := GetItem(userId, itemId)
		if item == nil { // 道具不存在
			if itemType == constant.ItemTypeSkin { // 全扣万能碎片
				allPowerItem := GetItem(userId, constant.ItemIdAllPowerChip)
				if allPowerItem == nil {
					return false
				}
				if allPowerItem.Num < itemNum {
					return false
				}
				// 扣除万能碎片
				allPowerItem.Num -= itemNum
				saveItem(allPowerItem)
				return true
			}
			return false
		}
		if item.Num < itemNum {
			// 数量不足
			if itemType != constant.ItemTypeSkin { // 皮肤碎皮不足 是否有万能碎片
				return false
			}
			allPowerItem := GetItem(userId, constant.ItemIdAllPowerChip)
			if allPowerItem == nil {
				return false
			}
			leftNeedNum := itemNum - item.Num
			if allPowerItem.Num < leftNeedNum { // 剩余碎片不足
				return false
			}
			// 扣除万能碎片
			allPowerItem.Num -= leftNeedNum
			saveItem(allPowerItem)
			itemNum -= leftNeedNum
		}
		item.Num -= itemNum
		saveItem(item)
	}
	// 记录道具
	var itemDatas []*include.ItemData
	itemDatas = append(itemDatas, &include.ItemData{UserId: p.Attr.UserID, ItemId: itemId, Num: -itemNum})
	RecordItem(itemDatas)
	return true
}

// 选择或道具
func (r *Room) ChoiceOrGiveUpItem(userId int32, choiceItemReq *proto3.ChoiceItemReq) (resp *proto3.ChoiceItemResp) {
	task, ok := r.taskMap[userId]
	if !ok {
		return &proto3.ChoiceItemResp{Ret: proto3.ErrEnum_Error_Operation_Fail}
	}
	tempDropItemsMap := task.TempDropItemsMap
	if tempDropItemsMap == nil {
		return &proto3.ChoiceItemResp{Ret: proto3.ErrEnum_Error_Operation_Fail}
	}

	if choiceItemReq.IsChoice == 1 { // 选择道具
		items := tempDropItemsMap[int32(choiceItemReq.Type)]
		if len(items) == 0 {
			return &proto3.ChoiceItemResp{Ret: proto3.ErrEnum_Error_Operation_Fail}
		}
		var exists = false
		var item *proto3.Item
		for _, v := range items {
			if v.Id == choiceItemReq.ItemId {
				item = v
				exists = true
			}
		}
		if !exists {
			return &proto3.ChoiceItemResp{Ret: proto3.ErrEnum_Error_Operation_Fail}
		}

		task.Skill = item
		if item != nil {
			task.TotalGetSkills = append(task.TotalGetSkills, item.Id)
		}
	} else { // 放弃道具
		if r.roomInfo.IsWolfMan(userId) {
			r.addSkillToSkillPool(constant.WolfManSkillPool, tempDropItemsMap)
		} else {
			r.addSkillToSkillPool(constant.NormalManSkillPool, tempDropItemsMap)
		}
	}

	// 删除临时数据
	delete(task.TempDropItemsMap, int32(choiceItemReq.Type))
	r.taskMap[userId] = task
	// 清除其它技能
	r.ClearSkills(userId)

	return &proto3.ChoiceItemResp{Ret: proto3.ErrEnum_Error_Pass, Type: choiceItemReq.Type, ItemId: choiceItemReq.ItemId}
}

// 放弃的道具如果在技能池里存在 则归还到技能池
func (r *Room) addSkillToSkillPool(skillPoolType int32, tempDropItemMap map[int32][]*proto3.Item) {
	skillPoolMap := make(map[int32]int32)
	if skillPoolType == constant.NormalManSkillPool {
		skillPoolMap = r.normalManSkillPoolMap
	} else if skillPoolType == constant.WolfManSkillPool {
		skillPoolMap = r.wolfManSkillPoolMap
	}
	for _, v := range tempDropItemMap {
		for _, item := range v {
			skillNum, ok := skillPoolMap[item.Id]
			if !ok {
				continue
			}
			skillPoolMap[item.Id] = skillPoolMap[item.Id] + skillNum
		}
	}
	if skillPoolType == constant.NormalManSkillPool {
		r.normalManSkillPoolMap = skillPoolMap
	} else if skillPoolType == constant.WolfManSkillPool {
		r.wolfManSkillPoolMap = skillPoolMap
	}
}

// 处理2号和3号道具选项卡
func (r *Room) DealSkillTab(p *Player, roommates []*proto3.Roommate) {
	second, third := reduceSkillTab(p)
	userId := p.Attr.UserID
	var skillOnOff []proto3.OnOff
	if second {
		r.addSkillTabNum(userId)
		skillOnOff = append(skillOnOff, proto3.OnOff_on)
	} else {
		skillOnOff = append(skillOnOff, proto3.OnOff_off)
	}
	if third {
		r.addSkillTabNum(userId)
		skillOnOff = append(skillOnOff, proto3.OnOff_on)
	} else {
		skillOnOff = append(skillOnOff, proto3.OnOff_off)
	}
	if len(roommates) == 0 {
		return
	}
	for _, v := range roommates {
		if v.UserId == userId {
			v.SkillOnOff = skillOnOff
		}
	}
}

// 增加技能选项卡栏位数
func (r *Room) addSkillTabNum(userId int32) {
	task, ok := r.taskMap[userId]
	if !ok { // 狼人未初始化任务
		skillTabNumMap := make(map[int32]int32)
		skillTabNumMap[userId] = 1
		task = &Task{UserId: userId, SkillTabNumMap: skillTabNumMap}
	} else {
		skillTabNumMap := task.SkillTabNumMap
		if skillTabNumMap == nil {
			task.SkillTabNumMap = make(map[int32]int32)
			task.SkillTabNumMap[userId] = 1
		} else {
			task.SkillTabNumMap[userId]++
		}
	}
	r.taskMap[userId] = task
}

// 扣除2号和3号技能选项卡
func reduceSkillTab(p *Player) (second bool, third bool) {
	second = p.ReduceItem(constant.ItemTypeItem, constant.ItemTabSecond, 1)
	third = p.ReduceItem(constant.ItemTypeItem, constant.ItemTabThird, 1)
	if second {
		items := make([]*proto3.Item, 0)
		items = append(items, &proto3.Item{Id: constant.ItemTabSecond, Num: -1})
		if third {
			items = append(items, &proto3.Item{Id: constant.ItemTabThird, Num: -1})
		}
		cmd := proto3.ProtoCmd_CMD_ItemsResp
		pbData := &proto3.ItemResp{Items: items}
		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	}
	return second, third
}

// 选择道具
func (r *Room) ChoiceItem(reqMap map[int32]*proto3.ChoiceItemReq) {
	var userId int32
	var choiceItemReq *proto3.ChoiceItemReq
	for k, v := range reqMap {
		userId = k
		choiceItemReq = v
	}
	if choiceItemReq == nil {
		return
	}
	choiceItemResp := r.ChoiceOrGiveUpItem(userId, choiceItemReq)
	player := global.GloInstance.GetPlayer(userId)
	if player != nil {
		p := player.(*Player)
		p.Pid.Cast("choiceItemResp", choiceItemResp)
	}
}

// 使用道具
func (r *Room) UseItem(reqMap map[int32]*proto3.UseItemReq) {
	var userId int32
	var useItemReq *proto3.UseItemReq
	for k, v := range reqMap {
		userId = k
		useItemReq = v
	}
	var itemType int32
	var itemId int32
	var useItemedUserId int32
	if useItemReq != nil {
		itemType = int32(useItemReq.Type)
		itemId = useItemReq.ItemId
		useItemedUserId = useItemReq.UseItemedUserId
	}

	player := global.GloInstance.GetPlayer(userId)
	if player != nil {
		p := player.(*Player)
		var errNum proto3.ErrEnum
		if itemType == int32(proto3.DropTypeEnum_item) { // 道具

		} else if itemType == int32(proto3.DropTypeEnum_skill) { // 技能
			errNum = r.reduceSkill(userId, itemId)
			if errNum != proto3.ErrEnum_Error_Pass {
				p.ErrorResponse(proto3.ErrEnum_Error_Goods_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Goods_NotExists)])
				return
			}
			role, randRoomId := r.UseSkill(itemId, userId, useItemedUserId)

			pbData := &proto3.UseItemResp{Ret: int32(proto3.ErrEnum_Error_Pass), ItemId: itemId, UseItemUserId: userId, UseItemedUserId: useItemedUserId, Role: role, RandRoomId: randRoomId}
			r.spreadPlayers(proto3.ProtoCmd_CMD_UseItemResp, pbData)
		}
	}
}

// 出售道具
func (p *Player) SellItem(itemId int32) {
	itemConfig := tableconfig.ItemConfigs.GetItemConfigById(itemId)
	if itemConfig == nil {
		logger.Log.Errorln("SellItem itemConfig is nil")
		return
	}
	if itemConfig.Sell <= 0 {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	item := GetItem(p.Attr.UserID, itemId)
	if item == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Goods_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Goods_NotExists)])
		return
	}
	// 加金币
	acqGold := item.Num * itemConfig.Sell
	p.Attr.Gold += acqGold
	// 扣道具
	p.ReduceItem(itemConfig.Type, itemConfig.Id, item.Num)

	cmd := proto3.ProtoCmd_CMD_SellItemResp
	pbData := &proto3.SellItemResp{ErrNum: proto3.ErrEnum_Error_Pass, Gold: p.Attr.Gold, ItemId: itemId}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
}

// 使用幸运卡
func (r *Room) UseLuckyCard(userId int32, itemId int32, itemNum int32) {
	if itemNum <= 0 {
		logger.Log.Errorf("UseLuckyCard itemNum %v is error", itemNum)
		return
	}
	player := global.GloInstance.GetPlayer(userId)
	if player == nil {
		logger.Log.Errorln("UseLuckyCard player is nil")
		return
	}
	p := player.(*Player)
	item := GetItem(userId, itemId)
	if item == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	if r.taskMap == nil {
		logger.Log.Errorln("UseLuckyCard taskMap is nil")
		return
	}
	task := r.taskMap[userId]
	if task == nil { // 狼人为nil
		task = &Task{UserId: userId}
	}
	itemConfig := tableconfig.ItemConfigs.GetItemConfigById(itemId)
	// 扣道具
	flag := p.ReduceItem(itemConfig.Type, itemId, itemNum)
	if !flag {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	// 设置值
	cardMap := task.LuckCardMap
	if cardMap == nil {
		cardMap = make(map[int32]int32)
		cardMap[itemId] = itemNum
	} else {
		cardMap[itemId] = cardMap[itemId] + itemNum
	}
	task.LuckCardMap = cardMap

	var leftNum int32
	leftItem := GetItem(userId, itemId)
	if leftItem != nil {
		leftNum = item.Num
	}
	cmd := proto3.ProtoCmd_CMD_UseLuckyCardResp
	pbData := &proto3.UseLuckyCardResp{ErrNum: proto3.ErrEnum_Error_Pass, ItemId: itemId, LeftNum: leftNum}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
}

// 使用背包道具
func (p *Player) UseBagItem(itemId int32, itemNum int32) {
	if itemId <= 0 || itemNum <= 0 {
		return
	}
	if p == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	userId := p.Attr.UserID
	item := GetItem(userId, itemId)
	if item == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Goods_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Goods_NotExists)])
		return
	}
	if item.Num < itemNum {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}

	// 扣道具
	itemConfig := tableconfig.ItemConfigs.GetItemConfigById(itemId)
	p.ReduceItem(itemConfig.Type, itemId, itemNum)
	if itemId == constant.ItemIdDoubleCard { // 翻倍卡
		if itemNum > 1 { // 不能超过1个
			p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
			return
		}
		if p.Room == nil {
			p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
			return
		}
		p.Room.DoubleRewards(p.Attr.UserID)
	} else {
		switch itemConfig.Type {
		case constant.ItemTypeFixedRewardBox: // 固定奖励宝箱
			p.AddItems(itemConfig.OpenBox)
		case constant.ItemTypeRandomRewardBox: // 随机奖励宝箱
			openBox := strings.Split(itemConfig.OpenBox, ",")
			openNum, _ := strconv.Atoi(openBox[0])
			reward := ""
			for i := 0; i < openNum; i++ {
				config := GetDropGroupConfigByDropId(util.ToInt(openBox[0]))
				reward += strconv.Itoa(config.ItemId) + "," + strconv.Itoa(config.Num) + "|"
			}
			if len(reward) > 0 {
				reward = reward[:len(reward)-1]
			}
			p.AddItems(reward)
		default:
			break
		}
	}
	cmd := proto3.ProtoCmd_CMD_UseBagItemResp
	pbData := &proto3.UseBagItemResp{ErrNum: proto3.ErrEnum_Error_Pass, Item: &proto3.Item{Id: itemId, Num: item.Num}}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})

	p.SendRedPoint(proto3.RedPointEnum_items_red, GetItemRedDot(p.Attr.UserID))
}

func GetItemRedDot(userId int32) (redDot []int32) {
	itemMap := db.ItemMap[userId] // 加载背包会从数据库取 放入内存 如果内存没有 则玩家没有道具
	if itemMap == nil {
		return
	}
	for itemId, _ := range itemMap {
		item := tableconfig.ItemConfigs.GetItemConfigById(itemId)
		if item == nil {
			continue
		}
		if (item.Type == constant.ItemTypeFixedRewardBox || item.Type == constant.ItemTypeRandomRewardBox) && itemMap[itemId].Num > 0 { // 宝箱类型
			redDot = append(redDot, item.Id)
		}
	}
	return
}
