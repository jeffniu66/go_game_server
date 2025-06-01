package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"strconv"
)

func saveStore(store *include.Store) {
	if store == nil {
		return
	}
	store.Update = include.Update
	db.StoreMap[store.UserId] = store
}

func getStore(userId int32) *include.Store {
	store, ok := db.StoreMap[userId]
	if !ok {
		store = db.SelectStore(userId)
		if store == nil {
			return nil
		}
		db.StoreMap[userId] = store
	}
	return store
}

// 商店请求
func (p *Player) StoreReq() {
	userId := p.Attr.UserID
	store := getStore(userId)
	var isRefreshMysSkin bool
	if store == nil {
		store = initStore(userId)
		isRefreshMysSkin = true
	} else {
		resetBox(store)
		isRefreshMysSkin = resetMysSkin(store)
		resetSkin(store)
		resetGold(store)
	}
	if isRefreshMysSkin {
		refreshMysSkin(store)
	}
	saveStore(store)

	var freeSkinChipId int32
	skinConfig := getFreeSkinChipId(userId)
	if skinConfig == nil {
		freeSkinChipId = -1
	} else {
		freeSkinChipId = skinConfig.Skinid
	}

	cmd := proto3.ProtoCmd_CMD_StoreResp
	msgData := &proto3.StoreResp{
		AdvanceBoxUseTimes: store.BoxUseAdNum,
		AdvanceBoxEndTime:  util.GetDayTimeStamp(24, 0, 0),
		MysSkinBuyTimes:    store.MysSkinBuyNum,
		MysSkinEndTime:     getMysSkinEndTime(),
		MysSkinChipId:      store.MysSkinChipId,
		FreeSkinBuyTimes:   store.SkinUseAdNum,
		FreeSkinEndTime:    util.GetDayTimeStamp(24, 0, 0),
		FreeSkinChipId:     freeSkinChipId,
		FreeGoldBuyTimes:   store.GoldViewAdNum,
		FreeGoldEndTime:    util.GetDayTimeStamp(24, 0, 0),
	}
	p.SendMessage(&Message{Cmd: cmd, PbData: msgData})
}

// 获取免费皮肤碎片id
func getFreeSkinChipId(userId int32) *tableconfig.SkinConfig {
	skinConfigs := tableconfig.SkinConfigs.SkinTypeMap[constant.AdSkinType]
	for _, skinConfig := range skinConfigs {
		item := GetItem(userId, skinConfig.Skinid)
		if item == nil {
			return &skinConfig
		}
		if item.Num >= skinConfig.Collectnum {
			continue
		}
		return &skinConfig
	}
	return nil
}

func getMysSkinEndTime() int32 {
	cur6Time := util.GetDayTimeStamp(6, 0, 0)
	cur18Time := util.GetDayTimeStamp(18, 0, 0)
	curTime := util.UnixTime()
	if curTime <= cur6Time {
		return cur6Time
	} else if curTime > cur6Time && curTime <= cur18Time {
		return cur18Time
	}
	return util.GetDayTimeStamp(24, 0, 0) + 6*60*60
}

func initStore(userId int32) *include.Store {
	store := &include.Store{
		UserId:                 userId,
		BoxUseAdNum:            0,
		BoxLastUseAdTime:       0,
		MysSkinBuyNum:          0,
		MysSkinLastRefreshTime: 0,
		MysSkinChipId:          0,
		SkinUseAdNum:           0,
		SkinLastRefreshTime:    0,
		GoldViewAdNum:          0,
		GoldLastViewAdTime:     0,
	}
	return store
}

// 每天0点重置金币
func resetGold(store *include.Store) {
	goldLastViewAdTime := store.GoldLastViewAdTime
	curTime := util.UnixTime()
	if !util.IsSameDay(goldLastViewAdTime, curTime) {
		store.GoldViewAdNum = 0
		store.GoldLastViewAdTime = curTime
	}
}

// 每天0点重置皮肤
func resetSkin(store *include.Store) {
	skinLastRefreshTime := store.SkinLastRefreshTime
	curTime := util.UnixTime()
	if !util.IsSameDay(skinLastRefreshTime, curTime) {
		store.SkinUseAdNum = 0
		store.SkinLastRefreshTime = curTime
	}
}

// 每天6点 18点重置
func resetMysSkin(store *include.Store) bool {
	cur6Time := util.GetDayTimeStamp(6, 0, 0)
	cur18Time := util.GetDayTimeStamp(18, 0, 0)
	curTime := util.UnixTime()
	mysSkinLastRefreshTime := store.MysSkinLastRefreshTime
	if curTime > cur6Time && curTime < cur18Time {
		// 上次刷新时间小于6点
		if mysSkinLastRefreshTime < cur6Time {
			store.MysSkinBuyNum = 0
			store.MysSkinLastRefreshTime = curTime
			return true
		}
	} else if curTime > cur18Time {
		// 上次刷新时间小于18点
		if mysSkinLastRefreshTime < cur18Time {
			store.MysSkinBuyNum = 0
			store.MysSkinLastRefreshTime = curTime
			return true
		}
	}
	return false
}

// 每天0点重置宝箱
func resetBox(store *include.Store) {
	boxLastUseAdTime := store.BoxLastUseAdTime
	curTime := util.UnixTime()
	if !util.IsSameDay(boxLastUseAdTime, curTime) {
		store.BoxUseAdNum = 0
		store.BoxLastUseAdTime = curTime
	}
}

// 开宝箱
func (p *Player) OpenBox(boxId int32) {
	store := getStore(p.Attr.UserID)
	if store == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	if boxId == constant.ItemIdNormalBox { // 普通宝箱
		openBoxNeedGold := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.OpenBoxNeedGold))
		gold := p.Attr.Gold
		if gold < openBoxNeedGold { // 金币不足
			p.ErrorResponse(proto3.ErrEnum_Error_Gold_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Gold_Not_Enough)])
			return
		}
		// 扣金币
		p.Attr.Gold -= openBoxNeedGold
	} else { // 高级宝箱

		resetBox(store)

		boxMaxAdNumPerDay := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.BoxMaxAdNumPerDay))
		if store.BoxUseAdNum >= boxMaxAdNumPerDay { // 广告次数不足
			p.ErrorResponse(proto3.ErrEnum_Error_Ad_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Ad_Not_Enough)])
			return
		}
		// 校验广告时间是否正常
		if !checkAdTimeInterval(store) {
			return
		}

		// 增加广告使用次数
		store.BoxUseAdNum++
		store.LastViewAdTime = util.UnixTime()

		// 记录广告数据
		RecordAdData(p, constant.AdvanceBoxAd)
	}
	reward := getBoxReward(boxId)
	if len(reward) == 0 {
		return
	}
	p.AddItems(reward)

	saveStore(store)

	cmd := proto3.ProtoCmd_CMD_OpenBoxResp
	msgData := &proto3.OpenBoxResp{ErrNum: proto3.ErrEnum_Error_Pass, Gold: p.Attr.Gold, UseAdNum: store.BoxUseAdNum}
	p.SendMessage(&Message{Cmd: cmd, PbData: msgData})
}

func getBoxReward(boxId int32) string {
	itemConfig := tableconfig.ItemConfigs.GetItemConfigById(boxId)
	if itemConfig == nil {
		logger.Log.Errorln("OpenBox itemConfig is nil!")
		return ""
	}
	dropGroupId := util.ToInt(itemConfig.Effect)
	dropGroupList := tableconfig.DropGroupConfigs.GetDropGroupMap(int(dropGroupId))
	if dropGroupList == nil || len(dropGroupList) == 0 {
		logger.Log.Errorln("OpenBox dropGroupList is nil!")
		return ""
	}
	var dropNumConst int32
	if boxId == constant.ItemIdNormalBox {
		dropNumConst = constant.NormalBoxDropNum
	} else {
		dropNumConst = constant.AdvanceBoxDropNum
	}
	normalBoxDropNum := util.ToInt(tableconfig.ConstsConfigs.GetValueById(dropNumConst))
	var reward = ""
	for i := 0; i < int(normalBoxDropNum); i++ {
		dropGroupConfig := GetDropGroupConfig(dropGroupList)
		if dropGroupConfig == nil {
			continue
		}
		reward = reward + strconv.Itoa(dropGroupConfig.ItemId) + "," + strconv.Itoa(dropGroupConfig.Num) + "|"
	}
	if len(reward) > 0 {
		reward = reward[:len(reward)-1]
	}
	return reward
}

// 看广告得金币
func (p *Player) ViewAdGetGold() {
	store := getStore(p.Attr.UserID)
	if store == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	maxAdNum := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.GoldMaxAdNumPerDay))
	if store.GoldViewAdNum >= maxAdNum {
		p.ErrorResponse(proto3.ErrEnum_Error_Ad_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Ad_Not_Enough)])
		return
	}
	// 校验广告时间间隔
	if !checkAdTimeInterval(store) {
		return
	}
	viewAdGetGoldNum := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.ViewAdGetGoldNum))
	p.Attr.Gold += viewAdGetGoldNum
	store.GoldViewAdNum++
	store.LastViewAdTime = util.UnixTime()

	saveStore(store)

	cmd := proto3.ProtoCmd_CMD_FreeItemResp
	msgData := &proto3.FreeItemResp{ErrNum: proto3.ErrEnum_Error_Pass, ItemType: proto3.ItemTypeEnum_gold_item, Gold: p.Attr.Gold, UseAdNum: store.GoldViewAdNum}
	p.SendMessage(&Message{Cmd: cmd, PbData: msgData})

	// 记录广告数据
	RecordAdData(p, constant.FreeGoldAd)
}

// 购买道具
func (p *Player) BuyItems(tid int32) {
	storeConfig := tableconfig.StoresConfigs.GetStoreById(tid)
	if storeConfig.Id == 0 {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	currencyType := storeConfig.Type
	if currencyType == constant.ItemIdGold { // 金币
		if p.Attr.Gold < storeConfig.Price {
			p.ErrorResponse(proto3.ErrEnum_Error_Gold_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Gold_Not_Enough)])
			return
		}
		p.Attr.Gold = p.Attr.Gold - storeConfig.Price
	} else if currencyType == constant.ItemIdGemStone { // 宝石
		if p.Attr.GemStone < storeConfig.Price {
			p.ErrorResponse(proto3.ErrEnum_Error_GemStore_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_GemStore_Not_Enough)])
			return
		}
		p.Attr.GemStone = p.Attr.GemStone - storeConfig.Price
	} else {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	reward := util.ToStr(storeConfig.Itemid) + "," + util.ToStr(storeConfig.Salenum)
	p.AddItems(reward)

	cmd := proto3.ProtoCmd_CMD_BuyItemsResp
	msgData := &proto3.BuyItemsResp{ErrNum: proto3.ErrEnum_Error_Pass, ItemType: proto3.ItemTypeEnum_prop_item, Gold: p.Attr.Gold, Gemstone: p.Attr.GemStone}
	p.SendMessage(&Message{Cmd: cmd, PbData: msgData})
}

// 获取常驻广告皮肤
func (p *Player) ViewAdGetSkin() {
	store := getStore(p.Attr.UserID)
	if store == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}

	resetSkin(store)

	skinMaxAdNumPerDay := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.SkinMaxAdNumPerDay))
	if store.SkinUseAdNum >= skinMaxAdNumPerDay { // 广告数不足
		p.ErrorResponse(proto3.ErrEnum_Error_Ad_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Ad_Not_Enough)])
		return
	}
	// 校验广告时间间隔
	if !checkAdTimeInterval(store) {
		return
	}
	skinConfig := getFreeSkinChipId(p.Attr.UserID)
	if skinConfig == nil {
		cmd := proto3.ProtoCmd_CMD_FreeItemResp
		msgData := &proto3.FreeItemResp{ErrNum: proto3.ErrEnum_Error_Pass, ItemType: proto3.ItemTypeEnum_skin_item, Gold: p.Attr.Gold, UseAdNum: store.SkinUseAdNum, NextChipId: -1}
		p.SendMessage(&Message{Cmd: cmd, PbData: msgData})
		return
	}
	dropGroupList := tableconfig.DropGroupConfigs.GetDropGroupMap(int(skinConfig.Dropid))
	if dropGroupList == nil || len(dropGroupList) == 0 {
		logger.Log.Errorln("ViewAdGetSkin dropGroupList is null")
		return
	}

	dropGroupConfig := GetDropGroupConfig(dropGroupList)
	if dropGroupConfig == nil {
		logger.Log.Errorln("ViewAdGetSkin dropGroupConfig is null")
		return
	}

	store.SkinUseAdNum++
	store.LastViewAdTime = util.UnixTime()

	p.AddItems(util.ToStr(skinConfig.Skinid) + "," + util.ToStr(int32(dropGroupConfig.Num)))
	saveStore(store)

	cmd := proto3.ProtoCmd_CMD_FreeItemResp
	msgData := &proto3.FreeItemResp{ErrNum: proto3.ErrEnum_Error_Pass, ItemType: proto3.ItemTypeEnum_skin_item, Gold: p.Attr.Gold, UseAdNum: store.SkinUseAdNum, NextChipId: skinConfig.Skinid}
	p.SendMessage(&Message{Cmd: cmd, PbData: msgData})

	// 记录广告数据
	RecordAdData(p, constant.FreeSkinAd)
}

// 购买神秘皮肤
func (p *Player) BuyMysSkin(itemId int32) {
	store := getStore(p.Attr.UserID)
	if store == nil {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	if itemId != store.MysSkinChipId {
		p.ErrorResponse(proto3.ErrEnum_Error_Operation_Fail, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Operation_Fail)])
		return
	}
	mysSkinConfigList, _ := tableconfig.SkinConfigs.SkinTypeMap[constant.MysSkinType]
	mysSkinConfig := mysSkinConfigList[0]
	price := mysSkinConfig.Price
	if p.Attr.GemStone < price {
		p.ErrorResponse(proto3.ErrEnum_Error_GemStore_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_GemStore_Not_Enough)])
		return
	}
	// 次数判断
	if store.MysSkinBuyNum >= mysSkinConfig.Limit {
		p.ErrorResponse(proto3.ErrEnum_Error_MysSkin_Times_Not_Enough, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_MysSkin_Times_Not_Enough)])
		return
	}

	dropGroupConfig := getDropGroupConfig(mysSkinConfig.Dropid, itemId)

	// 加道具
	p.AddItems(util.ToStr(itemId) + "," + util.ToStr(int32(dropGroupConfig.Num)))

	p.Attr.GemStone = p.Attr.GemStone - price
	store.MysSkinBuyNum = store.MysSkinBuyNum + 1
	saveStore(store)

	cmd := proto3.ProtoCmd_CMD_BuyItemsResp
	msgData := &proto3.BuyItemsResp{ErrNum: proto3.ErrEnum_Error_Pass, ItemType: proto3.ItemTypeEnum_skin_item, Gold: p.Attr.Gold, Gemstone: p.Attr.GemStone,
		MysSkinBuyTimes: store.MysSkinBuyNum}
	p.SendMessage(&Message{Cmd: cmd, PbData: msgData})
}

func getDropGroupConfig(dropId int32, itemId int32) *tableconfig.DropGroupConfig {
	dropGroupConfigs := tableconfig.DropGroupConfigs.DropGroupMap[int(dropId)]
	for _, dropGroupConfig := range dropGroupConfigs {
		if int32(dropGroupConfig.ItemId) == itemId {
			return &dropGroupConfig
		}
	}
	return nil
}

// 刷新神秘皮肤
func refreshMysSkin(store *include.Store) {
	mysSkinChipId := store.MysSkinChipId
	mysSkinConfigList, _ := tableconfig.SkinConfigs.SkinTypeMap[constant.MysSkinType]
	mysSkinConfig := mysSkinConfigList[0]
	dropGroupList := tableconfig.DropGroupConfigs.GetDropGroupMap(int(mysSkinConfig.Dropid))
	if dropGroupList == nil || len(dropGroupList) == 0 {
		logger.Log.Errorln("refreshMysSkin dropGroupList is nil!")
		return
	}

	chipId := randMysSkinRecursion(dropGroupList, mysSkinChipId)
	store.MysSkinChipId = chipId
}

func randMysSkinRecursion(dropGroupList []tableconfig.DropGroupConfig, mysSkinChipId int32) int32 {
	dropGroupConfig := GetDropGroupConfig(dropGroupList)
	if dropGroupConfig == nil {
		logger.Log.Errorln("refreshMysSkin dropGroupConfig is nil!")
		return 0
	}
	if !validMysSkinChipIdSame(dropGroupConfig, mysSkinChipId) {
		return int32(dropGroupConfig.ItemId)
	}
	return randMysSkinRecursion(dropGroupList, mysSkinChipId)
}

func validMysSkinChipIdSame(dropGroupConfig *tableconfig.DropGroupConfig, mysSkinChipId int32) bool {
	if int32(dropGroupConfig.ItemId) == mysSkinChipId {
		return true
	}
	return false
}

// 校验广告时间间隔
func checkAdTimeInterval(store *include.Store) bool {
	if store == nil {
		return false
	}
	curTime := util.UnixTime()
	adTimeCheckInterval := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.AdTimeCheckInterval))
	if curTime-store.LastViewAdTime < adTimeCheckInterval {
		return false
	}
	return true
}

func GetStoreRedDot(userId int32) (redDots []int32) {
	store := getStore(userId)
	if store == nil {
		store = initStore(userId)
	} else {
		resetBox(store)
		resetSkin(store)
		resetGold(store)
	}
	saveStore(store)

	boxMaxAdNum := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.BoxMaxAdNumPerDay))
	if boxMaxAdNum-store.BoxUseAdNum > 0 {
		redDots = append(redDots, 1)
	}
	skinMaxAdNum := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.SkinMaxAdNumPerDay))
	if skinMaxAdNum-store.SkinUseAdNum > 0 {
		redDots = append(redDots, 2)
	}
	maxAdNum := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.GoldMaxAdNumPerDay))
	if maxAdNum-store.GoldViewAdNum > 0 {
		redDots = append(redDots, 3)
	}
	return redDots
}
