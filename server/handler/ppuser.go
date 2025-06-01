package handler

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/db"
	"go_game_server/server/game"
	"go_game_server/server/logger"
)

func init() {
	//game.Handler.RegistHandler(proto3.ProtoCmd_CMD_SdkBindReq, &proto3.SdkBindReq{}, handleSDKBind)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UserPageReq, &proto3.UserPageReq{}, handleUserPage)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_GetArchiveReq, &proto3.GetArchiveReq{}, handleGetArchive)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_GetNinjaGiftReq, &proto3.GetNinjaGiftReq{}, handleGetNinjaGift)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UseTitleReq, &proto3.UseTitleReq{}, handleUseTitle)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UseSkinReq, &proto3.UseSkinReq{}, handleUseSkin)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UnlockSkinReq, &proto3.UnlockSkinReq{}, handleUnlockSkin)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_SetPhotoReq, &proto3.SetPhotoReq{}, handleSetPhoto)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_SetBorderReq, &proto3.SetBorderReq{}, handleSetBorder)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UserReadMailReq, &proto3.UserReadMailReq{}, handleUserReadMail)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_StoreReq, &proto3.StoreReq{}, handleStoreReq)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_OpenBoxReq, &proto3.OpenBoxReq{}, handleOpenBox)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_FreeItemReq, &proto3.FreeItemReq{}, handleFreeItemReq)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_BuyItemsReq, &proto3.BuyItemsReq{}, handleBuyItems)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_SellItemReq, &proto3.SellItemReq{}, handleSellItem)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UseBagItemReq, &proto3.UseBagItemReq{}, handleUseBagItem)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_NewGuideReq, &proto3.NewGuideReq{}, handleNewGuide)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_ModifyNameReq, &proto3.ModifyNameReq{}, handleModifyName)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_GetFreshGiftReq, &proto3.GetFreshGiftReq{}, handleGetFreshGift)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_TopBoardReq, &proto3.TopBoardReq{}, handleTopBoard)
	game.Handler.RegistHandler(proto3.ProtoCmd_CMD_UserTopReq, &proto3.UserTopReq{}, handleTopUser)
}

//func handleSDKBind(req interface{}, player *game.Player) interface{} {
//	//msg := req.(*proto3.SdkBindReq)
//	pbData := &proto3.SdkBindResp{}
//	//if player.Other.Bind(msg.BindType, msg.Name) {
//	//	pbData.Name = msg.Name
//	//	pbData.BindType = msg.BindType
//	//}
//
//	player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_SdkBindResp, PbData: pbData})
//	return nil
//}

func handleUserPage(req interface{}, player *game.Player) interface{} {
	msg, b := req.(*proto3.UserPageReq)
	if !b {
		return nil
	}
	var cmd proto3.ProtoCmd
	var pbData interface{}
	if msg.PageType == proto3.UserPageEnum_home_page { // 个人主页
		cmd = proto3.ProtoCmd_CMD_UserInfoResp
		tmp := &proto3.UserInfoResp{}
		tmp.PlayerAttr = player.GetProto3PlayerAttr()
		pbData = tmp
	} else if msg.PageType == proto3.UserPageEnum_archive_page { // 成就
		cmd = proto3.ProtoCmd_CMD_UserArchiveResp
		pbData = player.GetProto3ArchivePage()

		game.RecordAction(player.Attr.UserID, constant.ClickAchievement)
	} else if msg.PageType == proto3.UserPageEnum_title_page { // 称号
		cmd = proto3.ProtoCmd_CMD_UserTitleResp
		pbData = player.GetProto3UserTitle()
	} else if msg.PageType == proto3.UserPageEnum_skin_page { // 皮肤
		cmd = proto3.ProtoCmd_CMD_UserSkinResp
		pbData = player.GetProto3UserSkin()

		game.RecordAction(player.Attr.UserID, constant.ClickRole)
	} else if msg.PageType == proto3.UserPageEnum_mail_page { // 邮件
		cmd = proto3.ProtoCmd_CMD_UserMailResp
		pbData = player.GetProto3UserMail()
	} else {
		return nil
	}
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	logger.Log.Infof("username:%s page cmd:%v, pbData:%v", player.Attr.Username, cmd, pbData)
	return nil
}

func handleGetArchive(req interface{}, player *game.Player) interface{} {
	msg, b := req.(*proto3.GetArchiveReq)
	if !b {
		return nil
	}
	// cmd := proto3.ProtoCmd_CMD_GetArchiveResp
	_ = player.GetProto3GetArchive(msg)
	cmd := proto3.ProtoCmd_CMD_UserArchiveResp
	pbData := player.GetProto3ArchivePage()
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	logger.Log.Infof("username:%s page cmd:%v, pbData:%v", player.Attr.Username, cmd, pbData)
	return nil
}

func handleGetNinjaGift(req interface{}, player *game.Player) interface{} {
	msg, b := req.(*proto3.GetNinjaGiftReq)
	if !b {
		return nil
	}
	cmd := proto3.ProtoCmd_CMD_GetNinjaGiftResp
	pbData := player.GetNinjaGift(msg.NinjaIdGift)
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	logger.Log.Infof("username:%s page cmd:%v, pbData:%v", player.Attr.Username, cmd, pbData)
	return nil
}

func handleUseTitle(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UseTitleReq)
	if msg.TitleId < 0 {
		return nil
	}
	b := player.UseTitle(msg.TitleId)
	if b {
		pbData := &proto3.UseTitleResp{Id: msg.TitleId}
		player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_UseTitleResp, PbData: pbData})
	}
	return nil
}

func handleUseSkin(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UseSkinReq)
	if msg.SkinId < 0 {
		return nil
	}
	b := player.UsedSkin(msg.SkinId)
	if b {
		pbData := &proto3.UseSkinResp{Id: msg.SkinId}
		player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_UseSkinResp, PbData: pbData})
	} else {
		player.ErrorResponse(proto3.ErrEnum_Error_NotGet, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_NotGet)])
	}
	return nil
}

func handleUnlockSkin(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UnlockSkinReq)
	if msg.AvatarId < 0 {
		return nil
	}
	pbData := player.UnlockSkin(msg.AvatarId)
	player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_UnlockSkinResp, PbData: pbData})
	return nil
}

func handleUserReadMail(req interface{}, player *game.Player) interface{} {
	msg, ok := req.(*proto3.UserReadMailReq)
	if !ok {
		return nil
	}
	cmd := proto3.ProtoCmd_CMD_UserReadMailResp
	pbData := &proto3.UserReadMailResp{}
	if msg.OptType == proto3.MailEnum_mail_read {
		db.ReadUserMail(msg.MailId, player.Attr.UserID, int32(proto3.CommonStatusEnum_true))
		pbData.MailId = msg.MailId
	} else if msg.OptType == proto3.MailEnum_mail_getone {
		mail := db.GetMail(msg.MailId)
		if mail == nil || mail.AnnexOpen == int32(proto3.CommonStatusEnum_true) {
			return nil
		}
		player.AddItems(mail.AnnexItems)
		db.OpenOneMail(msg.MailId, player.Attr.UserID, int32(proto3.CommonStatusEnum_true))
		pbData.MailId = msg.MailId
	} else if msg.OptType == proto3.MailEnum_mail_getall {
		player.OpenAllMail()
	} else if msg.OptType == proto3.MailEnum_mail_delallread {
		db.DelAllRead(player.Attr.UserID)
	}

	pbData.Success = proto3.CommonStatusEnum_true
	pbData.OptType = msg.OptType
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})

	player.SendRedPoint(proto3.RedPointEnum_mail_red, player.GetMailRedPoint())
	return nil
}

func handleSetPhoto(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.SetPhotoReq)
	if msg.PhotoId < 0 {
		return nil
	}
	b := player.SetPhoto(msg.PhotoId)
	if b {
		pbData := &proto3.SetPhotoResp{PhotoId: msg.PhotoId}
		player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_SetPhotoResp, PbData: pbData})
	} else {
		player.ErrorResponse(proto3.ErrEnum_Error_NotGet, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_NotGet)])
	}
	return nil
}

func handleSetBorder(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.SetBorderReq)
	if msg.BorderId < 0 {
		return nil
	}
	b := player.SetBorder(msg.BorderId)
	if b {
		pbData := &proto3.SetBorderResp{BorderId: msg.BorderId}
		player.SendMessage(&game.Message{Cmd: proto3.ProtoCmd_CMD_SetBorderResp, PbData: pbData})
	} else {
		player.ErrorResponse(proto3.ErrEnum_Error_NotGet, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_NotGet)])
	}
	return nil
}

func handleStoreReq(req interface{}, player *game.Player) interface{} {
	player.StoreReq()

	game.RecordAction(player.Attr.UserID, constant.ClickStore)
	return nil
}

func handleOpenBox(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.OpenBoxReq)
	player.OpenBox(msg.BoxId)
	return nil
}

func handleFreeItemReq(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.FreeItemReq)
	if msg == nil {
		return nil
	}
	if msg.ItemType == proto3.ItemTypeEnum_gold_item {
		player.ViewAdGetGold()
	} else if msg.ItemType == proto3.ItemTypeEnum_skin_item {
		player.ViewAdGetSkin()
	}
	return nil
}

func handleBuyItems(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.BuyItemsReq)
	if msg.ItemType == proto3.ItemTypeEnum_prop_item {
		player.BuyItems(msg.ItemId)
	} else {
		player.BuyMysSkin(msg.ItemId)
	}
	return nil
}

func handleSellItem(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.SellItemReq)
	player.SellItem(msg.ItemId)
	return nil
}

func handleUseBagItem(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UseBagItemReq)
	player.UseBagItem(msg.Item.Id, msg.Item.Num)
	return nil
}

func handleNewGuide(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.NewGuideReq)
	game.NewGuide(player, msg.GuideId)
	return nil
}

func handleModifyName(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.ModifyNameReq)
	player.ModifyName(msg.Name, msg.Sex)
	return nil
}

func handleGetFreshGift(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.GetFreshGiftReq)
	player.GetFreshGiftSkin(msg)
	return nil
}

func handleTopBoard(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.TopBoardReq)
	pbData := player.GetTopBoard(msg.TopType)
	cmd := proto3.ProtoCmd_CMD_TopBoardResp
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}

func handleTopUser(req interface{}, player *game.Player) interface{} {
	msg := req.(*proto3.UserTopReq)
	pbData := player.GetUserTop(msg.TopType, msg.UserId)
	cmd := proto3.ProtoCmd_CMD_UserTopResp
	player.SendMessage(&game.Message{Cmd: cmd, PbData: pbData})
	return nil
}
