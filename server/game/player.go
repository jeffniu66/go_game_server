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
	"runtime/debug"
	"strings"
	"sync"
)

var birthLock sync.Mutex

func LoadNames() {
	names := db.SelectNames()
	if len(names) == 0 {
		return
	}
	for _, v := range names {
		db.NameMap[v] = v
	}
}

func (p *Player) CreateUser(name, ip string, sex int32, openid, channel string) {
	name, sex, _ = GetRankName(tableconfig.NameZhConfigs)
	p.Attr = initUser(name, sex, openid, channel)
	p.Attr.Ip = ip            // 暂时没用到 20201218
	p.ScreenView = -100       // 暂时没用到 20201218
	if p.Attr.Channel == "" { // 如果渠道为空，兼容旧的账号表示dev渠道
		p.Attr.Channel = "dev"
	}

	// 初始化成就数据
	p.initAchievementMap()

	// 第一次登陆,角色初始化数据
	if p.Attr.FirstLogin {
		// 初始化道具
		initItems := tableconfig.ConstsConfigs.GetValueById(constant.InitItems)
		p.AddItems(initItems)

		go p.initNewRole()
		db.NameMap[p.Attr.Username] = p.Attr.Username // 加载到内存
	}
}
func (p *Player) KickOldClient(oldPlayer *Player) {
	oldPlayer.Pid.Stop()
}

func (p *Player) ConstructPlayer(userId int32) {
	// 重复登录的玩家,需要先摧毁之前的旧对象
	OldPlayer := global.GloInstance.GetPlayer(userId)
	if OldPlayer != nil {
		player := OldPlayer.(*Player)

		// 通知下线
		pbData := &proto3.LogoutResp{Status: proto3.CommonStatusEnum_true}
		player.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_LogoutResp, PbData: pbData})

		// 此处异步stop，需要先将旧的未保存的数据保留给新的连接
		p.Attr = player.Attr
		p.KickOldClient(player)
		// player.Pid.Stop()
	}

	// 初始化成就数据
	p.initAchievementMap()
	if p.Attr.Channel == "" { // 如果渠道为空，兼容旧的账号表示dev渠道
		p.Attr.Channel = "dev"
	}

	// 第一次登陆,角色初始化数据
	//if p.Attr.FirstLogin {
	//	p.initNewRole()
	//}
}

func LoadUserTitle(a *db.Attr) {
	a.UserTitle = &db.UserTitle{}
	a.UserTitle.UserID = a.UserID
	a.UserTitle.GetUserTitle(a.UserID)
}

func LoadUserAchievement(a *db.Attr) {
	a.UserArchiveTypeMap = db.GetUserArchivement(a.UserID)
	a.AchievementMap = make(map[int32]int32, 0)
}

func initUser(name string, sex int32, openid, channel string) *db.Attr {
	attrs := new(db.Attr)
	attr := attrs.InitData(name, "", "-", sex, openid, channel)
	return attr
}

func (p *Player) initNewRole() {
	defer func() {
		if err := recover(); err != nil {
			logger.Log.Errorf("initNewRole error!!!\n ", err)
			debug.PrintStack()
		}
	}()
	// 发送新人礼包
	p.AddNewMail(include.MailTypeRegister, "", proto3.DetailTypeEnum_key_text)
}

func (p *Player) AddNewMail(mailType int32, rewards string, detailType proto3.DetailTypeEnum) {
	mail := tableconfig.MailConfigs.GetMail(mailType)
	if rewards == "" {
		rewards = mail.Rewards
	}
	if mail != nil && mail.UseSwitch == "1" {
		userMail := &db.UserMail{
			UserID:     p.Attr.UserID,
			Title:      mail.Name,
			Content:    mail.Desc,
			AnnexItems: rewards,
			DetailType: int32(detailType),
		}
		db.SaveMail(userMail)
		p.SendRedPoint(proto3.RedPointEnum_mail_red, []int32{1})
	}
}

func (p *Player) SaveDBData() {
	// p.Attr.SaveData()
	if p.IsRobot == 1 {
		return
	}
	// TODO 任务的示例
	userID := p.Attr.UserID
	db.SaveItemsData(userID)  // 背包
	db.SaveStoreData(userID)  // 商店
	db.SaveActionData(userID) // 行为
	//p.Task.TaskMap.SaveData(userID)

	//p.SaveRedisData()
	p.QuaChange() // 赛季变更
	p.Attr.SaveData()
}

func (p *Player) SaveRedisData() {
	// 活跃数据写到redis，方便离线操作
	SetHashFields(GetRedisUserKey(p.Attr.UserID), map[string]interface{}{
		include.UserFieldLordName: p.Attr.Username,
		include.UserFieldCountry:  p.Attr.Country,
		include.UserFieldPower:    p.Attr.Power,
	})
}

func (p *Player) initRedisUser() {
	SetHashFields(GetRedisUserKey(p.Attr.UserID), map[string]interface{}{
		include.UserFieldLordName: p.Attr.Username,
		include.UserFieldCountry:  p.Attr.Country,
		include.UserFieldPower:    p.Attr.Power,
	})
}

func (p *Player) SendMessage(msg interface{}) {
	if p.IsRobot == 1 || p.IsOffLine == 1 {
		return
	}
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("player:%v SendMessage error =%s debug:%v ", p.Attr.Username, err, string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
		}
	}()
	if p.WriteChan == nil {
		logger.Log.Infof("userID:%d writeChan is nil", p.Attr.UserID)
		return
	}
	p.WriteChan <- msg
}

func (p *Player) GetLevel() int32 {
	return p.Attr.Level
}

func (p *Player) SetLevel(level int32) {
	p.Attr.Level = level
}

func (p *Player) calcPlayerAttr(curNum, addNum, maxNum int32, canLimit bool) int32 {
	if !canLimit {
		if curNum < maxNum {
			curNum = util.Min(curNum+addNum, maxNum)
		} else if addNum < 0 {
			curNum += addNum
		}
	} else {
		curNum += addNum
	}
	curNum = util.Max(curNum, 0)
	return curNum
}

func GetUserName(userID int32) string {
	if userID <= 0 {
		return ""
	}

	player := global.GloInstance.GetPlayer(userID)
	if player != nil {
		p := player.(*Player)
		return p.Attr.Username
	} else {
		lordName, err := GetRedisUserField(userID, include.UserFieldLordName)
		if nil == err {
			return lordName
		}

		return db.GetUsername(userID)
	}
}

func (p *Player) GetCommonFlag(flag uint32) int32 {
	return p.Attr.CommonFlag >> flag & 0x1
}

func (p *Player) ErrorResponse(errNum proto3.ErrEnum, errMsg string) {
	cmd := proto3.ProtoCmd_CMD_ErrResp
	pbData := &proto3.ErrResp{ErrCode: errNum, ErrMsg: errMsg}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
}

//func (p *Player) PushKick(kickType proto3.KickTypeEnum, extend int32, reason string) {
//	p.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_PushKick,
//		//PbData: &proto3.PushKick{KickType: kickType, Extend: extend, Reason: reason}})
//}

//func GetRedisUser(userId int32) *RedisUser {
//	return GetRedisUserByKey(GetRedisUserKey(userId))
//}

func (p *Player) IsInRoom() (room *Room) {
	// 根据room_id 判断
	if p.Attr.RoomID > 0 {
		roomFace := global.GloInstance.GetRoom(p.Attr.RoomID)
		if roomFace != nil {
			room := roomFace.(*Room)
			if room.IsInRoom(p.Attr.UserID) {
				return room
			} else {
				return nil
			}
		}
	}

	/* 全局搜索room消耗性能
	if p.RoomPid != nil {
		return p.Room
	}
	roomList := global.GloInstance.GetUsedRoomFaceList()
	logger.Log.Infof("rooms:%v", roomList)
	for _, roomFace := range roomList {
		if roomFace != nil {
			ro := roomFace.(*Room)
			if ro.IsInRoom(p.Attr.UserID) {
				room = ro
				break
			}
		}
	}
	*/
	return nil
}

func (p *Player) ExitRoom() {
	if p.RoomPid != nil && p.RoomPid.PidName != "" {
		p.RoomPid.Cast("playerExit", p.Attr.UserID)
	}
	p.Room = nil
	p.RoomPid = nil
}

func (p *Player) GameEnd(roomID int32) {
	if p.Attr.RoomID == roomID {
		p.Attr.RoomID = 0 // 清理玩家房间
		p.Room = nil
		p.RoomPid = nil
	}
}

func (p *Player) SendRedPoint(pointType proto3.RedPointEnum, redData []int32) {
	cmd := proto3.ProtoCmd_CMD_RedPointResp
	pbData := &proto3.RedPointResp{RedType: pointType, RedData: redData}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
}

func GetPlayerAiList(startNum, endNum int32, ch string) []*Player {
	ret := make([]*Player, 0)
	for i := startNum; i < endNum; i++ {
		tmp := NewAIPlayer(i, ch)
		ret = append(ret, tmp)
	}
	return ret
}

func NewAIPlayer(userID int32, ch string) *Player {
	p := &Player{
		Attr: &db.Attr{},
	}
	p.Attr.UserID = userID
	username, _, _ := GetRankName(tableconfig.NameZhConfigs)
	p.Attr.Username = username
	p.Attr.StarCount = 0
	p.Attr.Channel = ch
	p.IsRobot = 1
	LoadUserTitle(p.Attr)
	LoadUserAchievement(p.Attr)
	p.Attr.UseSkin = tableconfig.ConstsConfigs.GetRandValue(constant.RandSkin)
	p.Attr.UserPhoto = tableconfig.ConstsConfigs.GetRandValue(constant.RandPhoto)
	p.Attr.UserTitle.UseTitle = tableconfig.ConstsConfigs.GetRandValue(constant.RandTitle)
	global.GloInstance.AddPlayer(userID, p)
	return p
}

func (p *Player) OpenAllMail() {
	items := db.GetNoOpenMail(p.Attr.UserID)
	if len(items) <= 0 {
		return
	}
	itemMap := make(map[string]string, 0)
	for _, v := range items {
		vList := strings.Split(v, "|")
		for _, it := range vList {
			vv := strings.Split(it, ",")
			if len(vv) != 2 {
				continue
			}
			if iv, ok := itemMap[vv[0]]; ok {
				tmp := util.ToInt(iv) + util.ToInt(vv[1])
				itemMap[vv[0]] = util.ToStr(tmp)
			} else {
				itemMap[vv[0]] = vv[1]
			}
		}
	}
	item := ""
	i := 0
	for k, v := range itemMap {
		if i == 0 {
			item = k + "," + v
		} else {
			item += "|" + k + "," + v
		}
		i++
	}
	p.AddItems(item)
	db.OpenAllMail(p.Attr.UserID)
}

func SpreadWorldPlayer(player *Player, chatType proto3.ChatTypeEnum, chatData string, quickType, quickId int32, subsText []string) {
	cmd := proto3.ProtoCmd_CMD_HomeChatResp
	pbData := &proto3.HomeChatResp{}

	pbData.UserId = player.Attr.UserID
	pbData.ChatData = chatData
	pbData.ChatType = chatType
	pbData.Username = player.Attr.Username
	pbData.UserPhoto = player.Attr.UserPhoto
	pbData.QuickType = quickType
	pbData.QuickId = quickId
	pbData.SubsText = subsText

	// 广播
	go func() {
		for _, uid := range global.GloInstance.GetPlayerIDList() {
			p := global.GloInstance.GetPlayer(uid)
			if p == nil {
				logger.Log.Errorf("this uid:%d is nil player")
				continue
			}
			player := p.(*Player)
			if player.IsInRoom() == nil {
				player.SendMessage(&Message{Cmd: cmd, PbData: pbData})
			}
		}
	}()
}

// QuaChange 赛季变更
func (p *Player) QuaChange() {
	if p.IsRobot == 1 {
		return
	}
	oldQuaConfig, b := tableconfig.QuaConfigs.IsNowQua(p.Attr.UpdateTime)
	if b && oldQuaConfig == nil {
		logger.Log.Infof("oldQuaConfig is nil, please set qualfy.excel")
		return
	}
	if oldQuaConfig != nil && b { // 旧赛季存在，第一次更新-当前不在新赛季 发送旧赛季奖励
		cmd := proto3.ProtoCmd_CMD_QualfyChangeResp
		pbData := &proto3.QualfyChangeResp{
			RankId: tableconfig.QuaLevelConfs.StartID,
		}
		if p.Attr.RankID == oldQuaConfig.NeedLevel1 {
			p.AddItems(oldQuaConfig.Rewards1)
		}
		if p.Attr.RankID == oldQuaConfig.NeedLevel2 {
			p.AddItems(oldQuaConfig.Rewards2)
		}
		pbData.OldRankId = p.Attr.RankID
		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})

		levelConfs := tableconfig.QuaLevelConfs.GetQuaConfig(p.Attr.RankID)
		if levelConfs == nil {
			logger.Log.Info("this rank:%d quaLevelConfigs is nil", p.Attr.RankID)
		} else {
			p.Attr.RankID = levelConfs.ChangeLevel
		}
		if p.Attr.HisRankID < p.Attr.RankID {
			p.Attr.HisRankID = p.Attr.RankID
		}
	}
}

// 新手引导
func NewGuide(p *Player, guideId int32) {
	if guideId < 0 || guideId > 999 {
		return
	}
	guide := db.SelectGuide(p.Attr.UserID)
	if guide == nil {
		guide = &include.Guide{UserId: p.Attr.UserID, GuideIds: util.ToStr(guideId)}
		db.InsertGuide(guide)
	} else {
		guide.GuideIds = guide.GuideIds + "," + util.ToStr(guideId)
	}
	db.InsertGuide(guide)

	cmd := proto3.ProtoCmd_CMD_NewGuideResp
	pbData := &proto3.NewGuideResp{ErrNum: proto3.ErrEnum_Error_Pass}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
}

// 获取新手引导
func (p *Player) GetGuides() (guides []int32) {
	guide := db.SelectGuide(p.Attr.UserID)
	if guide == nil {
		return
	}
	guideStr := guide.GuideIds
	if len(guideStr) == 0 {
		return
	}
	ss := strings.Split(guideStr, ",")
	for _, v := range ss {
		guides = append(guides, util.ToInt(v))
	}
	return
}

// 修改角色名
func (p *Player) ModifyName(name string, sex int32) {
	if len(name) > constant.MaxNameLen {
		p.ErrorResponse(proto3.ErrEnum_Error_Username_OutLen, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Username_OutLen)])
		return
	}
	if tableconfig.SenWordConfigs.CheckSenWord(name) {
		p.ErrorResponse(proto3.ErrEnum_Error_Involving_SenWord, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Involving_SenWord)])
		return
	}
	// 校验名字重复
	if db.JudgeUserName(name) {
		p.ErrorResponse(proto3.ErrEnum_Error_UserName_Exists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_UserName_Exists)])
		return
	}
	item := GetItem(p.Attr.UserID, constant.ItemIdModifyName)
	if item == nil || item.Num <= 0 {
		p.ErrorResponse(proto3.ErrEnum_Error_Goods_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Goods_NotExists)])
		return
	}

	p.ReduceItem(constant.ItemTypeItem, constant.ItemIdModifyName, 1)

	p.Attr.Username = name
	if p.Attr.SexModify == 1 {
		p.Attr.Sex = sex
		p.Attr.SexModify = 2
	}
	p.Attr.SaveData()

	cmd := proto3.ProtoCmd_CMD_ModifyNameResp
	pbData := &proto3.ModifyNameResp{ErrNum: proto3.ErrEnum_Error_Pass}
	pbData.Name = p.Attr.Username
	pbData.Sex = p.Attr.Sex
	userItem := GetItem(p.Attr.UserID, constant.ItemIdModifyName)
	if userItem == nil {
		pbData.LeftModifyCard = 0
	} else {
		pbData.LeftModifyCard = userItem.Num
	}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
}

// param1 姓名 param2 性别 param3 索引
func GetRankName(n *tableconfig.NameZhConfigCol) (string, int32, string) {
	var preIndex, surIndex, nameIndex int32 = -1, -1, -1
	prefix, surname, name := "", "", ""
	sex := int32(0)
	l := int32(len(n.NameZhConfigList))
	t := 0 // 避免死循环
	// 姓
	for {
		t++
		if t > 100 {
			break
		}
		index := util.RandInt(0, l-1)
		nameConf := n.NameZhConfigList[int(index)]
		if nameConf.Surname != "" {
			surname = nameConf.Name
			surIndex = nameConf.ID
			break
		}
	}
	// 前缀
	t = 0
	for {
		t++
		if t > 100 {
			break
		}
		index := util.RandInt(0, l-1)
		nameConf := n.NameZhConfigList[int(index)]
		if nameConf.Prefix != "" {
			prefix = nameConf.Prefix
			preIndex = nameConf.ID
			break
		}
	}

	// 名
	t = 0
	for {
		t++
		if t > 100 {
			break
		}
		index := util.RandInt(0, l-1)
		nameConf := n.NameZhConfigList[int(index)]
		if nameConf.Name != "" {
			name = nameConf.Name
			nameIndex = nameConf.ID
			sex = nameConf.Sex
			break
		}
	}

	model := util.RandInt(1, 3)
	if model == include.NameTypeSurname {
		prefix = ""
		preIndex = -1
	} else if model == include.NameTypePrefixName {
		surname = ""
		surIndex = -1
	} else {

	}
	indexStr := util.ToStr(preIndex) + "," + util.ToStr(surIndex) + "," + util.ToStr(nameIndex)
	fullName := prefix + surname + name
	if db.NameMap[fullName] != "" {
		return GetRankName(n)
	}
	return fullName, sex, indexStr
}

func (p *Player) GetRedData() (redDots []*proto3.RedDot) {
	userId := p.Attr.UserID
	redDots = append(redDots, &proto3.RedDot{RedType: proto3.RedPointEnum_store_red, RedData: GetStoreRedDot(userId)})
	redDots = append(redDots, &proto3.RedDot{RedType: proto3.RedPointEnum_items_red, RedData: GetItemRedDot(userId)})
	redDots = append(redDots, &proto3.RedDot{RedType: proto3.RedPointEnum_skin_red, RedData: p.GetSkinRedPoint()})
	redDots = append(redDots, &proto3.RedDot{RedType: proto3.RedPointEnum_title_red, RedData: p.GetTitleRedPoint()})
	redDots = append(redDots, &proto3.RedDot{RedType: proto3.RedPointEnum_archive_red, RedData: p.GetArchiveRedPoint()})
	redDots = append(redDots, &proto3.RedDot{RedType: proto3.RedPointEnum_mail_red, RedData: p.GetMailRedPoint()})
	return
}
