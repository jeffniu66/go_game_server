package game

import (
	"go_game_server/proto3"
	"go_game_server/server/db"
	"go_game_server/server/include"
	"go_game_server/server/logger"
)

type TopBoard struct {
	RankTopList []*include.UserTop // 排行榜从1开始
	Length      int32
}

func (c *TopBoard) Sort() {
	// toDo 优先根据数据库进行排行
}

func (c *TopBoard) UpdateBoard(p *Player) *include.UserTop {
	if p.IsRobot == 1 {
		return nil
	}
	var index int32 = -1
	uS := p.Attr.StarCount
	uL := p.Attr.Level
	uE := p.Attr.Exp
	curIndex := p.rankTopId - 1
	if curIndex <= 100 && curIndex >= 0 { // 在排行榜内
		for i := curIndex - 1; i >= 0; i-- { // 往上走
			v := c.RankTopList[i]
			if uS > v.StarCount || uS == v.StarCount && uL > v.Level || uS == v.RankId && uL == v.Level && uE >= v.Exp {
				index = i
				// 右移留出位置
				c.RankTopList[index].TopId += 1
				c.RankTopList[index+1] = c.RankTopList[index]
				logger.Log.Info(">>>>>>>>>>>", index, v)
			} else {
				break
			}
		}
		logger.Log.Info(">>>>>>>>>>>pppppp", index)
		if index == -1 {
			for i := curIndex + 1; i < c.Length; i++ { // 往下走
				v := c.RankTopList[i]
				if uS < v.StarCount || uS == v.StarCount && uL < v.Level || uS == v.RankId && uL == v.Level && uE <= v.Exp {
					index = i
					// 左移留出位置
					if index < topLength-1 {
						c.RankTopList[index].TopId -= 1
						c.RankTopList[index-1] = c.RankTopList[index]
						logger.Log.Info(">>>>>>>>>>>", index, v)
					}
				} else {
					if i == c.Length-1 && c.Length < topLength {
						index = c.Length
						c.Length += 1
					}
					break
				}
			}
		}
		if index == -1 {
			index = p.rankTopId - 1
		}
	} else {
		for i := c.Length - 1; i >= 0; i-- {
			v := c.RankTopList[i]
			// 段位、等级、当前经验 谁后更新谁入榜
			if uS > v.StarCount || uS == v.StarCount && uL > v.Level || uS == v.RankId && uL == v.Level && uE >= v.Exp {
				index = i
				if index < topLength-1 {
					c.RankTopList[index].TopId += 1
					c.RankTopList[index+1] = c.RankTopList[index]
				}
			} else {
				if i == c.Length-1 && c.Length < topLength {
					index = c.Length
					c.Length += 1
				}
				break
			}
		}
	}

	nowTop := new(include.UserTop)
	nowTop.UserId = p.Attr.UserID
	nowTop.Username = p.Attr.Username
	nowTop.UserPhoto = p.Attr.UserPhoto
	nowTop.TopId = -1
	nowTop.RankId = p.Attr.RankID
	nowTop.Star = p.Attr.Star
	nowTop.StarCount = p.Attr.StarCount
	nowTop.Level = p.Attr.Level
	nowTop.Exp = p.Attr.Exp
	nowTop.UpdateTime = p.Attr.UpdateTime
	if index < topLength && index > 0 {
		nowTop.TopId = index + 1
		c.RankTopList[index] = nowTop
		p.rankTopId = index + 1
	}
	return nowTop
}

var topBoard *TopBoard
var topLength int32 = 100

func InitTopBoard() {
	var length int32 = topLength
	c := new(TopBoard)

	c.Length = length
	c.RankTopList = make([]*include.UserTop, int(length), int(length))
	db.SelectRankBoard(length, c.RankTopList, &c.Length)
	topBoard = c
	return
}

func (p *Player) GetTopBoard(topType proto3.TopBoardEnum) *proto3.TopBoardResp {
	ret := new(proto3.TopBoardResp)
	ret.SelfTop = &proto3.UserTop{
		UserId:   p.Attr.UserID,
		Username: p.Attr.Username,
	}
	ret.TopType = topType
	if topType == proto3.TopBoardEnum_rank_top {
		for i := 0; i < int(topBoard.Length); i++ {
			v := topBoard.RankTopList[i]
			tmp := &proto3.UserTop{
				UserId:    v.UserId,
				Username:  v.Username,
				UserPhoto: v.UserPhoto,
				TopId:     v.TopId,
				RankId:    v.RankId,
				Star:      v.Star,
				Level:     v.Level,
				Exp:       v.Exp,
			}
			ret.TopList = append(ret.TopList, tmp)
		}
		ret.SelfTop.TopId = p.rankTopId
	}
	return ret
}

func (p *Player) GetUserTop(topType proto3.TopBoardEnum, userID int32) *proto3.UserTopResp {
	ret := new(proto3.UserTopResp)
	ret.PlayerAttr = new(proto3.PlayerAttr)
	ret.TopType = topType
	// userInfo := global.GloInstance.GetPlayer(userID) 直接查库
	a := db.GetUserByUserId(userID)
	ret.PlayerAttr.UserId = a.UserID
	ret.PlayerAttr.Level = a.Level
	ret.PlayerAttr.Exp = a.Exp
	ret.PlayerAttr.MaxExp = a.MaxExp
	ret.PlayerAttr.Gold = a.Gold
	ret.PlayerAttr.Gemstone = a.GemStone
	ret.PlayerAttr.RankId = a.RankID
	ret.PlayerAttr.Star = a.Star
	ret.PlayerAttr.StarCount = a.StarCount
	ret.PlayerAttr.HisRankId = a.HisRankID
	ret.PlayerAttr.NinjaId = a.NinjaID
	ret.PlayerAttr.ArchivePoint = a.ArchivePoint
	ret.PlayerAttr.GameDuration = a.GameDuration
	ret.PlayerAttr.MatchGameNum = a.MatchGameNum
	ret.PlayerAttr.MatchWinNum = a.MatchWinNum
	ret.PlayerAttr.MatchWolfNum = a.MatchWolfNum
	ret.PlayerAttr.WolfWinNum = a.WolfWinNum
	ret.PlayerAttr.PoorWinNum = a.PoorWinNum
	ret.PlayerAttr.OfflineNum = a.OfflineNum
	ret.PlayerAttr.VoteTotal = a.VoteTotal
	ret.PlayerAttr.VoteCorrectTotal = a.VoteCorrectTotal
	ret.PlayerAttr.KillTotal = a.KillTotal
	ret.PlayerAttr.WolfKillTotal = a.WolfKillTotal
	ret.PlayerAttr.PoorKillTotal = a.PoorKillTotal
	ret.PlayerAttr.BekilledTotal = a.BekilledTotal
	ret.PlayerAttr.BevoteedTotal = a.BevoteedTotal
	ret.PlayerAttr.UserBorder = a.UserBorder
	ret.PlayerAttr.UserPhoto = a.UserPhoto
	ret.PlayerAttr.UseSkin = a.UseSkin
	ret.PlayerAttr.GotSkins = a.GotSkins
	ret.PlayerAttr.Username = a.Username
	ret.PlayerAttr.SexModify = a.SexModify
	ret.PlayerAttr.Sex = a.Sex
	ret.PlayerAttr.UseTitle = a.UserTitle.UseTitle
	return ret
}

func (p *Player) CheckUpTop() {
	p.rankTopId = -1
	for i := 0; i < int(topBoard.Length); i++ {
		v := topBoard.RankTopList[i]
		if v.UserId == p.Attr.UserID {
			p.rankTopId = v.TopId
		}
	}
}
