package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"strings"
)

func (p *Player) initAchievementMap() {
	for _, v := range include.ArchiveTypeList {
		switch v {
		case include.ArchiveTypeKill:
			p.Attr.AchievementMap[v] = p.Attr.KillTotal
		case include.ArchiveTypeVoteOut:
			p.Attr.AchievementMap[v] = p.Attr.BevoteedTotal
		case include.ArchiveTypeVoteWolf:
			p.Attr.AchievementMap[v] = p.Attr.VoteCorrectTotal
		case include.ArchiveTypeFirstOut:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.KeepFirstOut
		case include.ArchiveTypeTask:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalTask
		case include.ArchiveTypeTaskGold:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalGold
		case include.ArchiveTypeSoulTask:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalSoulTask
		case include.ArchiveTypeLevel:
			p.Attr.AchievementMap[v] = p.Attr.Level
		case include.ArchiveTypeNinjaLevel:
			p.Attr.AchievementMap[v] = p.Attr.NinjaID
		case include.ArchiveTypeArcpoint:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalArchive
		case include.ArchiveTypeKillPoor:
			p.Attr.AchievementMap[v] = p.Attr.WolfKillTotal
		case include.ArchiveTypeKillWolf:
			p.Attr.AchievementMap[v] = p.Attr.PoorKillTotal
		case include.ArchiveTypeVoteTrue:
			p.Attr.AchievementMap[v] = p.Attr.VoteCorrectTotal
		case include.ArchiveTypeGold:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalGold
		case include.ArchiveTypeSkin:
			l := strings.Split(p.Attr.GotSkins, ",")
			p.Attr.AchievementMap[v] = int32(len(l))
		case include.ArchiveTypeMatch:
			p.Attr.AchievementMap[v] = p.Attr.MatchGameNum
		case include.ArchiveTypeMatchWin:
			p.Attr.AchievementMap[v] = p.Attr.MatchWinNum
		case include.ArchiveTypeRank:
			p.Attr.AchievementMap[v] = p.Attr.RankID
		case include.ArchiveTypeItem:
		case include.ArchiveTypeDuration:
			p.Attr.AchievementMap[v] = p.Attr.GameDuration
		case include.ArchiveTypeAd:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalAd
		case include.ArchiveTypePhoto:
		case include.ArchiveTypeBorder:
		case include.ArchiveTypeKeepWolf:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.KeepWolf
		case include.ArchiveTypeKeepPoor:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.KeepPoor
		case include.ArchiveTypeKeepNoItem:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.KeepNoItem
		case include.ArchiveTypeWolfDay:
			p.Attr.AchievementMap[v] = p.Attr.UserTitle.TotalWolfDay
		case include.ArchiveTypeVotePoor:
			p.Attr.AchievementMap[v] = p.Attr.VoteFailedTotal
		default:
		}
	}
}

func (p *Player) getArciveValue(archiveType int32) int32 {
	var value int32 = -1
	if v, ok := p.Attr.AchievementMap[archiveType]; ok {
		value = v
	}
	return value
}

func (p *Player) AlterUserTitle() {
	newGotTitle := &proto3.UserNewGotTitleResp{}
	newTitles := ""
	for _, archiveType := range include.ArchiveTypeList {
		archiveValue := p.getArciveValue(archiveType)
		// 称号
		userTitleList := tableconfig.TitleConfigs.GetTitle(include.TitleType, archiveType, archiveValue)
		for i := range userTitleList {
			if userTitleList[i].IsGot != include.GetStatusFalse {
				b := p.IsNewGotTitle(userTitleList[i].ID)
				if b {
					newGotTitle.TitldIds = append(newGotTitle.TitldIds, userTitleList[i].ID)
					if newTitles != "" {
						newTitles += ","
					}
					newTitles += util.ToStr(userTitleList[i].ID)
				}
			}
		}
		// 成就
		archiveList := tableconfig.TitleConfigs.GetTitle(include.ArchiveType, archiveType, archiveValue)
		p.SetUserArchive(archiveType, archiveValue, archiveList)
	}
	if len(newGotTitle.TitldIds) > 0 {
		logger.Log.Infof("user:%s get new title:%v", p.Attr.Username, newGotTitle)
		cmd := proto3.ProtoCmd_CMD_UserNewGotTitleResp
		p.SendMessage(&Message{Cmd: cmd, PbData: newGotTitle})
		p.Attr.UserTitle.TitleRedData = newTitles
		// pbData := &proto3.RedPointResp{RedType: proto3.RedPointEnum_title_red, RedData: newGotTitle.TitldIds}
		// p.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_RedPointResp, PbData: pbData})
		redData := p.GetTitleRedData()
		p.SendRedPoint(proto3.RedPointEnum_title_red, redData)
	}
}

func (p *Player) GetTitleRedData() []int32 {
	var ret []int32
	titleRedList := strings.Split(p.Attr.UserTitle.TitleRedData, ",")
	for _, v := range titleRedList {
		if v != "" {
			ret = append(ret, util.ToInt(v))
		}

	}
	return ret
}

func (p *Player) GetSkinRedData() []int32 {
	var ret []int32
	skinRedList := strings.Split(p.Attr.UserTitle.SkinRedData, ",")
	for _, v := range skinRedList {
		if v != "" {
			ret = append(ret, util.ToInt(v))
		}
	}
	return ret
}

func (p *Player) GetProto3PlayerAttr() (ret *proto3.PlayerAttr) {
	a := p.Attr
	ret = &proto3.PlayerAttr{}
	ret.UserId = a.UserID
	ret.Level = a.Level
	ret.Exp = a.Exp
	ret.MaxExp = a.MaxExp
	ret.Gold = a.Gold
	ret.Gemstone = a.GemStone
	ret.RankId = a.RankID
	ret.Star = a.Star
	ret.StarCount = a.StarCount
	ret.HisRankId = a.HisRankID
	ret.NinjaId = a.NinjaID
	ret.ArchivePoint = a.ArchivePoint
	ret.GameDuration = a.GameDuration
	ret.MatchGameNum = a.MatchGameNum
	ret.MatchWinNum = a.MatchWinNum
	ret.MatchWolfNum = a.MatchWolfNum
	ret.WolfWinNum = a.WolfWinNum
	ret.PoorWinNum = a.PoorWinNum
	ret.OfflineNum = a.OfflineNum
	ret.VoteTotal = a.VoteTotal
	ret.VoteCorrectTotal = a.VoteCorrectTotal
	ret.KillTotal = a.KillTotal
	ret.WolfKillTotal = a.WolfKillTotal
	ret.PoorKillTotal = a.PoorKillTotal
	ret.BekilledTotal = a.BekilledTotal
	ret.BevoteedTotal = a.BevoteedTotal
	ret.UserBorder = a.UserBorder
	ret.UserPhoto = a.UserPhoto
	ret.UseSkin = a.UseSkin
	ret.GotSkins = a.GotSkins
	ret.Username = a.Username
	ret.SexModify = a.SexModify
	ret.Sex = a.Sex
	return
}

// 未领取的成就页
func (p *Player) GetProto3ArchivePage() (ret *proto3.UserArchiveResp) {
	a := p.Attr
	ret = &proto3.UserArchiveResp{}
	redData := []int32{}

	p.initAchievementMap()
	p.AlterUserTitle()
	for k, v := range p.Attr.UserArchiveTypeMap {
		for i := 0; i < len(v); i++ {
			// 可领取 或者 返回类型最后一条
			if v[i].GotStatus != include.GetStatusYes || i == len(v)-1 {
				tmp := proto3.Archivement{
					ArchiveId:     v[i].ArchiveID,
					ArchiveStatus: proto3.CommonStatusEnum(v[i].GotStatus),
					ArchiveNum:    a.AchievementMap[k],
				}
				ret.ArchivementList = append(ret.ArchivementList, &tmp)
				if tmp.ArchiveStatus == proto3.CommonStatusEnum_true {
					redData = append(redData, tmp.ArchiveId)
				}

				// 最后一条
				if i == len(v)-1 {
					// 判断是否达到领取点
					archiveValue := p.getArciveValue(v[i].ArchiveType)
					titleConfig := tableconfig.TitleConfigs.TitleConfigMap[v[i].ArchiveID]
					if archiveValue > titleConfig.NeedValue {
						archiveList := tableconfig.TitleConfigs.GetTitle(include.ArchiveType, v[i].ArchiveType, archiveValue)
						logger.Log.Infof("userTitle:%v", archiveList)
						p.SetUserArchive(v[i].ArchiveType, archiveValue, archiveList)
					}
				}
				break
			}
		}
	}
	ret.NinjaId = a.NinjaID
	ret.ArchivePoint = a.ArchivePoint
	ret.MaxArchivePoint = a.MaxArchivePoint
	ret.NinjaIdGift = a.NinjaIDGift

	pbData := &proto3.RedPointResp{RedType: proto3.RedPointEnum_archive_red, RedData: redData}
	p.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_RedPointResp, PbData: pbData})

	logger.Log.Infof("userArchive:%v archiveMap:%v", *ret, p.Attr.UserArchiveTypeMap)
	return ret
}

// 领取的成就
func (p *Player) GetProto3GetArchive(req *proto3.GetArchiveReq) (ret *proto3.GetArchiveResp) {
	a := p.Attr
	ret = &proto3.GetArchiveResp{}
	titleConfig := tableconfig.TitleConfigs.TitleConfigMap[req.ArchiveId]
	archiveList := p.Attr.UserArchiveTypeMap[titleConfig.TitleType]
	// 按不同类型找到当前称号 并返回
	for i := 0; i < len(archiveList); i++ {
		// 找到下一个archive
		v := archiveList[i]
		if archiveList[i].ArchiveID == titleConfig.Next {
			ret.Archivement = &proto3.Archivement{
				ArchiveId:     titleConfig.Next,
				ArchiveNum:    a.AchievementMap[titleConfig.TitleType],
				ArchiveStatus: proto3.CommonStatusEnum(v.GotStatus),
			}
		}

		// 当前成就变更为已领取
		if v.ArchiveID == int32(req.ArchiveId) {
			if v.GotStatus == include.GetStatusYes {
				continue
			}
			v.GotStatus = include.GetStatusYes

			// 最后一条领取点 计算出下一条达成的成就
			if i == len(archiveList)-1 {
				// 判断是否达到领取点
				archiveValue := p.getArciveValue(v.ArchiveType)
				if archiveValue >= titleConfig.NeedValue {
					list := tableconfig.TitleConfigs.GetTitle(include.ArchiveType, archiveList[i].ArchiveType, archiveValue)
					logger.Log.Infof("userTitle:%v", archiveList)
					p.SetUserArchive(archiveList[i].ArchiveType, archiveValue, list)

					ret.Archivement = &proto3.Archivement{
						ArchiveId:     titleConfig.Next,
						ArchiveNum:    a.AchievementMap[titleConfig.TitleType],
						ArchiveStatus: proto3.CommonStatusEnum(include.GetStatusNo),
					}
				} else {
					ret.Archivement = &proto3.Archivement{
						ArchiveId:     titleConfig.Next,
						ArchiveNum:    a.AchievementMap[titleConfig.TitleType],
						ArchiveStatus: proto3.CommonStatusEnum(include.GetStatusFalse),
					}
				}
			}

			//计算忍级变化
			p.AddItems(titleConfig.Rewards)
		}
	}
	ret.NinjaId = a.NinjaID
	ret.ArchivePoint = a.ArchivePoint
	ret.MaxArchivePoint = a.MaxArchivePoint
	ret.NinjaIdGift = a.NinjaIDGift
	return ret
}

func (p *Player) UpdateLevel() {
	levelData := &p.Attr.PlayerAttr.UserData
	level, exp, maxExp := tableconfig.LevelConfigs.GetNextLevel(levelData.Level, levelData.Exp)
	if levelData.Level != level {
		//升级发送邮件
		levConf := tableconfig.LevelConfigs.GetLevelConfig(levelData.Level)
		if levConf != nil {
			p.AddNewMail(include.MailTypeLevelUp, levConf.Rewards, proto3.DetailTypeEnum_key_text)
		}
	}
	levelData.Level = level
	levelData.Exp = exp
	levelData.MaxExp = maxExp
}

func (p *Player) UpdateArchive() {
	ninjaID, archivePoint, maxArchive := tableconfig.NinjaConfigs.GetLevel(p.Attr.NinjaID, p.Attr.ArchivePoint)
	if ninjaID != p.Attr.NinjaID {
		// 初始化 0阶
		if p.Attr.NinjaID == 0 {
			p.Attr.NinjaIDGift = p.Attr.NinjaID
		}
		if p.Attr.NinjaIDGift == 0 {
			p.Attr.NinjaIDGift = 1
		}
		p.Attr.NinjaID, p.Attr.ArchivePoint, p.Attr.MaxArchivePoint = ninjaID, archivePoint, maxArchive
	}
}

func (p *Player) GetNinjaGift(ninjaIDGift int32) *proto3.GetNinjaGiftResp {
	ret := &proto3.GetNinjaGiftResp{}
	ninjaConfig := tableconfig.NinjaConfigs.GetNinjaConfig(ninjaIDGift)
	if ninjaConfig == nil {
		return ret
	}
	if ninjaIDGift > p.Attr.NinjaID {
		ret.Status = proto3.CommonStatusEnum_false
	} else {
		ret.Status = proto3.CommonStatusEnum_true
		gitTmp := ninjaIDGift // 保存单前领取点
		ret.NinjaIdGift = ninjaConfig.Next
		if ret.NinjaIdGift >= p.Attr.NinjaID || ret.NinjaIdGift == gitTmp { // 如果当前领取点等于最后一条，说明最后一个领取了
			ret.NinjaIdGift = 0
		}
	}
	p.Attr.NinjaIDGift = ret.NinjaIdGift
	p.AddItems(ninjaConfig.Rewards)
	return ret
}

// 获取取得的称号
func (p *Player) GetProto3UserTitle() (ret *proto3.UserTitleResp) {
	a := p.Attr
	ret = &proto3.UserTitleResp{}
	// ret.GotTitles = a.UserTitle.GotTitles
	ret.TitleList = make([]*proto3.UserTitle, 0)
	titles := strings.Split(a.UserTitle.GotTitles, ",")
	if len(titles) <= 0 {
		return ret
	}
	for _, col := range tableconfig.TitleConfigs.TitleConfigList {
		if col.MoudleType != include.TitleType {
			continue
		}
		gotList := make([]int32, 0)
		if p.IsGotTitle(col.ID) {
			gotList = append(gotList, col.ID)
		}
		existIndex := -1
		for i := 0; i < len(ret.TitleList); i++ {
			v := ret.TitleList[i]
			if v.TitleType == col.TitleType {
				existIndex = i
			}
		}
		// 该类型是否存在，存在的index
		if existIndex < 0 {
			tmp := &proto3.UserTitle{}
			tmp.TitleType = col.TitleType
			tmp.TitleNum = p.Attr.AchievementMap[tmp.TitleType]
			tmp.GotTitles = gotList
			ret.TitleList = append(ret.TitleList, tmp)
		} else {
			gots := ret.TitleList[existIndex].GotTitles
			isExist := false
			// 是否重复取得
			for v := range gots {
				if int32(v) == col.ID {
					isExist = true
				}
			}
			if !isExist && p.IsGotTitle(col.ID) {
				gots = append(gots, col.ID)
			}
			ret.TitleList[existIndex].GotTitles = gots
		}
	}
	ret.OutwardInfo = &proto3.OutWardInfo{
		NinjaId:  a.NinjaID,
		RankId:   a.RankID,
		UseTitle: a.UserTitle.UseTitle,
	}
	return ret
}

func (p *Player) IsGotTitle(titleID int32) bool {
	a := p.Attr
	titles := strings.Split(a.UserTitle.GotTitles, ",")
	for i := 0; i < len(titles); i++ {
		if titleID == util.ToInt(titles[i]) {
			return true
		}
	}
	return false
}

func (p *Player) IsGotSkin(skinID int32) bool {
	a := p.Attr
	skins := strings.Split(a.GotSkins, ",")
	for i := 0; i < len(skins); i++ {
		if skinID == util.ToInt(skins[i]) {
			return true
		}
	}
	return false
}

func (p *Player) IsNewGotTitle(titleID int32) bool {
	a := p.Attr
	if !p.IsGotTitle(titleID) {
		s := util.ToStr(titleID)
		a.UserTitle.GotTitles += "," + s
		return true
	}
	return false
}

// 使用取得的称号
func (p *Player) UseTitle(titleID int32) bool {
	a := p.Attr
	useFlag := p.IsGotTitle(titleID)
	if useFlag {
		a.UserTitle.UseTitle = titleID

		titleRedList := strings.Split(p.Attr.UserTitle.TitleRedData, ",")
		title := util.ToStr(titleID)
		for i, v := range titleRedList {
			if v == title {
				titleRedList = append(titleRedList[:i], titleRedList[i+1:]...)
				break
			}
		}

		p.Attr.UserTitle.TitleRedData = strings.Join(titleRedList, ",")
	}

	redData := p.GetTitleRedData()
	p.SendRedPoint(proto3.RedPointEnum_title_red, redData)
	return useFlag
}

func (p *Player) GetProto3UserSkin() *proto3.UserSkinResp {
	a := p.Attr
	ret := &proto3.UserSkinResp{}
	ret.OutwardInfo = &proto3.OutWardInfo{
		NinjaId:  a.NinjaID,
		RankId:   a.RankID,
		UseTitle: a.UserTitle.UseTitle,
	}
	skins := strings.Split(a.GotSkins, ",")
	for _, v := range skins {
		vv := util.ToInt(v)
		ret.GotSkins = append(ret.GotSkins, vv)
	}
	ret.UseSkin = a.UseSkin

	ret.PieceList = GetSkinPiece(a.UserID)
	logger.Log.Info("skin:", ret)
	return ret
}

// 使用取得的皮肤
func (p *Player) UsedSkin(skinID int32) bool {
	a := p.Attr
	useFlag := p.IsGotSkin(skinID)
	if useFlag {
		a.UseSkin = skinID
	}
	if p.Room != nil {
		p.Room.ChangeSkin(p.Attr.UserID, skinID)
	}

	skinRedList := strings.Split(p.Attr.UserTitle.SkinRedData, ",")
	skin := util.ToStr(skinID)
	for i, v := range skinRedList {
		if v == skin {
			skinRedList = append(skinRedList[:i], skinRedList[i+1:]...)
			break
		}
	}

	p.Attr.UserTitle.SkinRedData = strings.Join(skinRedList, ",")
	redData := p.GetSkinRedData()
	p.SendRedPoint(proto3.RedPointEnum_skin_red, redData)
	return useFlag
}

// 使用取得的皮肤
func (p *Player) UnlockSkin(avatarId int32) *proto3.UnlockSkinResp {
	ret := &proto3.UnlockSkinResp{}
	avatarConf := tableconfig.AvatarConfigs.GetConfig(avatarId)

	b := p.ReduceItem(int32(proto3.ItemTypeEnum_skin_item), avatarConf.ItemId, avatarConf.NeedNum)
	if !b {
		// p.ErrorResponse(proto3.ErrEnum_Error_Goods_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Goods_NotExists)])
		// todo
		return nil
	}

	ret.AvatarId = avatarId
	ret.PieceId = avatarConf.ItemId
	item := GetItem(p.Attr.UserID, avatarConf.ItemId)
	if item == nil {
		ret.PieceNum = 0
	} else {
		ret.PieceNum = item.Num
	}
	allPowerPiece := GetItem(p.Attr.UserID, constant.ItemIdAllPowerChip)
	if allPowerPiece == nil {
		ret.AllPowerPiece = 0
	} else {
		ret.AllPowerPiece = allPowerPiece.Num
	}
	p.addSkin(avatarId)
	return ret
}

func (p *Player) addSkin(skinID int32) {
	if p.IsGotSkin(skinID) {
		return
	}
	skin := util.ToStr(skinID)
	if len(p.Attr.GotSkins) <= 0 {
		p.Attr.GotSkins = "1" + "," + skin
	} else {
		p.Attr.GotSkins += "," + skin
	}

	if p.Attr.UserTitle.SkinRedData != "" {
		p.Attr.UserTitle.SkinRedData += ","
	}
	p.Attr.UserTitle.SkinRedData += skin

	skinRedData := p.GetSkinRedData()
	p.SendRedPoint(proto3.RedPointEnum_skin_red, skinRedData)
}

// 使用取得的头像
func (p *Player) SetPhoto(photoID int32) bool {
	a := p.Attr
	a.UserPhoto = photoID
	logger.Log.Info(p.Attr.UserPhoto)
	return true
}

// 使用取得的边框
func (p *Player) SetBorder(borderID int32) bool {
	a := p.Attr
	a.UserBorder = borderID
	return true
}

func (p *Player) IsGotArchivement(archiveType, archiveID int32) *db.UserArchive {
	a := p.Attr
	v := a.UserArchiveTypeMap[archiveType]
	for i := 0; i < len(v); i++ {
		if v[i].ArchiveID == archiveID {
			return v[i]
		}
	}
	return nil
}

func (p *Player) SetUserArchive(archiveType, nowValue int32, archiveList []tableconfig.UserTitle) {
	isRedPoint := false
	for _, archive := range archiveList {
		a := p.Attr
		isNil := true
		v, ok := a.UserArchiveTypeMap[archiveType]
		// 初始化
		if !ok {
			a.UserArchiveTypeMap[archiveType] = make([]*db.UserArchive, 0)
		}
		for i := 0; i < len(v); i++ {
			if v[i].ArchiveID == archive.ID {
				isNil = false
				if v[i].GotStatus != include.GetStatusYes {
					v[i].GotStatus = include.GetStatusNo
				}
			}
		}
		if isNil {
			isRedPoint = true
			next := tableconfig.TitleConfigs.TitleConfigMap[archive.ID].Next
			tmp := &db.UserArchive{
				ArchiveID:   archive.ID,
				ArchiveNext: next,
				UserID:      a.UserID,
				ArchiveType: archiveType,
				GotStatus:   include.GetStatusNo,
			}

			a.UserArchiveTypeMap[archiveType] = append(a.UserArchiveTypeMap[archiveType], tmp)
		}
	}
	if isRedPoint {
		pbData := &proto3.RedPointResp{RedType: proto3.RedPointEnum_archive_red}
		p.SendMessage(&Message{Cmd: proto3.ProtoCmd_CMD_RedPointResp, PbData: pbData})
	}
	return
}

func (p *Player) GetProto3UserMail() *proto3.UserMailResp {
	a := p.Attr
	ret := &proto3.UserMailResp{}
	mailList := db.GetMailList(a.UserID)
	for i := range mailList {
		tmp := mailList[i].ToProto3UserMail()
		ret.UserMailList = append(ret.UserMailList, tmp)
	}
	return ret
}

func (p *Player) GetSkinRedPoint() []int32 {
	ret := make([]int32, 0)
	// 皮肤
	skinData := strings.Split(p.Attr.UserTitle.SkinRedData, ",")
	for _, v := range skinData {
		if v == "" {
			continue
		}
		d := util.ToInt(v)
		ret = append(ret, d)
	}

	// 皮肤碎片p.Attr.UserID
	skin := p.GetProto3UserSkin()
	var powerNum int32 = 0
	allPowerItem := GetItem(p.Attr.UserID, constant.ItemIdAllPowerChip) // 万能碎片
	if allPowerItem != nil {
		powerNum += allPowerItem.Num
	}
	for _, v := range skin.PieceList {
		pieceNum := powerNum + v.PieceNum
		skinId := tableconfig.AvatarConfigs.GetEnoughItem(v.Id, pieceNum)
		if skinId > 0 {
			ret = append(ret, skinId)
		}
	}
	return ret
}

func (p *Player) GetTitleRedPoint() []int32 {
	ret := make([]int32, 0)
	titleData := strings.Split(p.Attr.UserTitle.TitleRedData, ",")
	for _, v := range titleData {
		if v == "" {
			continue
		}
		d := util.ToInt(v)
		ret = append(ret, d)
	}
	return ret
}

func (p *Player) GetArchiveRedPoint() []int32 {
	ret := make([]int32, 0)
	archiveList := p.GetProto3ArchivePage()
	for _, v := range archiveList.ArchivementList {
		if v.ArchiveStatus == proto3.CommonStatusEnum_true { // 待领取
			ret = append(ret, v.ArchiveId)
		}
	}
	return ret
}

func (p *Player) GetMailRedPoint() []int32 {
	ret := make([]int32, 0)
	mailList := p.GetProto3UserMail()
	for _, v := range mailList.UserMailList {
		if (v.ItemIds != "" && v.IsOpen != proto3.CommonStatusEnum_true) || v.IsRead != proto3.CommonStatusEnum_true { // 待领取
			ret = append(ret, v.Id)
			break
		}
	}
	return ret
}
