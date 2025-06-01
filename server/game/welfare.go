package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
)

func (p *Player) GetFreshGiftStep() {
	if p.Attr.FreshGiftStep >= constant.FreshGiftStepFinish {
		return
	}
	var pbData interface{}
	var cmd proto3.ProtoCmd
	if p.Attr.FreshGiftStep < constant.FreshGiftStepTiming {
		cmd = proto3.ProtoCmd_CMD_FreshGiftResp
		tmp := &proto3.FreshGiftResp{}
		tmp.GiftType = 1 // 写死新手礼包第一个
		tmp.NeedNum = tableconfig.ConstsConfigs.GetIdValue(constant.FreshGiftNeedNum)
		tmp.NowNum = p.Attr.MatchGameNum
		tmp.Rewards = tableconfig.ConstsConfigs.GetValueById(constant.FreshGiftRewards)
		pbData = tmp
	} else {
		cmd = proto3.ProtoCmd_CMD_FreshEndTimeResp
		tmp := &proto3.FreshEndTimeResp{}
		tmp.GiftType = 2 // 写死新手礼包第二个
		tmp.EndTime = p.Attr.FreshEndTime
		tmp.NowTime = util.UnixTime()
		tmp.Reward = tableconfig.ConstsConfigs.GetValueById(constant.FreshGiftEndReward)
		pbData = tmp
	}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	return
}

func (p *Player) GetFreshGiftSkin(msg *proto3.GetFreshGiftReq) {
	if msg == nil {
		return
	}
	if msg.GiftType == 1 {
		if p.Attr.FreshGiftStep != constant.FreshGiftStepStart || p.Attr.MatchGameNum < tableconfig.ConstsConfigs.GetIdValue(constant.FreshGiftNeedNum) {
			return
		}

		p.Attr.FreshGiftStep = constant.FreshGiftStepTiming
		p.Attr.FreshEndTime = util.UnixTime() + tableconfig.ConstsConfigs.GetIdValue(constant.FreshGiftEndTime)*3600

		p.GetFreshGiftStep()

		RecordAction(p.Attr.UserID, constant.AcqMatchReward)
	} else if msg.GiftType == 2 {
		if p.Attr.FreshGiftStep != constant.FreshGiftStepTiming {
			return
		}
		if util.UnixTime() < p.Attr.FreshEndTime {
			return
		}

		p.Attr.FreshGiftStep = constant.FreshGiftStepFinish
	}
	avatarId := util.ToInt(msg.Reward)
	p.addSkin(avatarId)
	cmd := proto3.ProtoCmd_CMD_GetFreshGiftResp
	pbData := &proto3.GetFreshGiftResp{
		GiftType: msg.GiftType,
		Reward:   msg.Reward,
	}
	p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	return
}
