package game

import (
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/db"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/sdk"
	"math/rand"
	"runtime/debug"
	"time"
)

const saveUserTime = 100
const (
	autoOutNormal  = 0 // 在房房间
	autoOutGameEnd = 1 // 正常退出
	autoOutGaming  = 2 // 游戏中退出
)

type WorldBoxMsgData struct {
	Type   int32 // 0获得,1删除
	LandID int32 // 有宝箱的地块id
}

type Player struct {
	Attr            *db.Attr
	Pid             *global.PidObj
	RoomPid         *global.PidObj
	Room            *Room
	WriteChan       chan interface{} // socket发送信息到client的通道
	ScreenView      int32            // 屏幕视野点
	MessageBox      []string
	IsOptOut        bool // 0-非 1-主动退出
	IsRobot         int32
	LoginRespStatus int32 // 登录返回状态
	IsOffLine       int32 // 掉线
	SessionKey      string
	rankTopId       int32 // 段位旁行榜
}

type LoginStateData struct {
	State     int32                      // 登录状态
	Name      string                     // 玩家账号名
	BanStamp  int32                      // 封禁时间戳
	BanReason string                     // 封禁原因
	BindInfos map[int32]*sdk.SDKBindInfo // sdk绑定信息
}

func CreatePid(writeChan chan interface{}, msg *proto3.LoginReq, ip string) (ret int32, playerPid *global.PidObj, userID int32) {
	player := &Player{WriteChan: writeChan}
	//if len(msg.Username) > constant.MaxNameLen {
	//	player.ErrorResponse(proto3.ErrEnum_Error_Username_OutLen, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Username_OutLen)])
	//	return -1, nil, 0
	//}
	//
	//if msg.Channel != "" && msg.Channel != "dev" {
	//	acctInfo, b := sdk.SdkUtil.GetSDKUserInfo(sdk.SDKChannel(msg.Channel), "", msg.Token)
	//	if acctInfo == nil || !b {
	//		return -1, nil, 0
	//	}
	//	logger.Log.Infof("GetSDKUserInfo = %v", acctInfo)
	//
	//	userIdList := db.GetUserByOpenId(acctInfo.Name)
	//	if len(userIdList) == 0 { // 不存在账号，则创建角色
	//		//if tableconfig.SenWordConfigs.CheckSenWord(msg.Username) {
	//		//	player.ErrorResponse(proto3.ErrEnum_Error_Involving_SenWord, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Involving_SenWord)])
	//		//	return -1, nil, 0
	//		//}
	//		//// 校验名字重复
	//		//if db.JudgeUserName(msg.Username) {
	//		//	player.ErrorResponse(proto3.ErrEnum_Error_UserName_Exists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_UserName_Exists)])
	//		//	return -1, nil, 0
	//		//}
	//		player.CreateUser(msg.Username, ip, msg.Sex, acctInfo.Name, msg.Channel)
	//	} else { // 存在账号
	//		attr := db.GetUserByUserId(userIdList[0])
	//		if attr != nil {
	//			player.Attr = attr
	//			player.ConstructPlayer(attr.UserID)
	//		}
	//	}
	//	sdk.TokenMap[player.Attr.UserID] = acctInfo.Name
	//	if acctInfo.ErrCode == 0 {
	//		player.LoginRespStatus = 1
	//	} else {
	//		player.LoginRespStatus = acctInfo.ErrCode
	//	}
	//	player.SessionKey = acctInfo.SessionKey
	//	player.Attr.Channel = msg.Channel // channel 未入库，V1.0.1版本 后期可以删除
	//} else { // 无渠道登录
	//	attr := db.GetUserByUserName(msg.Username)
	//	if attr == nil {
	//		player.CreateUser(msg.Username, ip, msg.Sex, "", "")
	//	} else {
	//		player.Attr = attr
	//		LoadUserTitle(player.Attr)
	//		LoadUserAchievement(player.Attr)
	//		player.ConstructPlayer(attr.UserID)
	//	}
	//	player.LoginRespStatus = 1
	//}
	//
	//// 先去渠道校验
	///*
	//	if msg.Channel == "wx" || msg.Channel == "qq" {
	//		wxLoginResp := sdk.AppletWeChatLogin(msg.Token, msg.Channel)
	//		if wxLoginResp == nil { // 登录失败
	//			return -1, nil, 0
	//		}
	//		logger.Log.Infof("WxLoginResp = %v", wxLoginResp)
	//		openId := wxLoginResp.OpenId
	//		userIdList := db.GetUserByOpenId(openId)
	//		if len(userIdList) == 0 { // 不存在账号，则创建角色
	//			//if tableconfig.SenWordConfigs.CheckSenWord(msg.Username) {
	//			//	player.ErrorResponse(proto3.ErrEnum_Error_Involving_SenWord, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Involving_SenWord)])
	//			//	return -1, nil, 0
	//			//}
	//			//// 校验名字重复
	//			//if db.JudgeUserName(msg.Username) {
	//			//	player.ErrorResponse(proto3.ErrEnum_Error_UserName_Exists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_UserName_Exists)])
	//			//	return -1, nil, 0
	//			//}
	//			player.CreateUser(msg.Username, ip, msg.Sex, openId, msg.Channel)
	//		} else { // 存在账号
	//			attr := db.GetUserByUserId(userIdList[0])
	//			if attr != nil {
	//				player.Attr = attr
	//				player.ConstructPlayer(attr.UserID)
	//			}
	//		}
	//		sdk.TokenMap[player.Attr.UserID] = openId
	//		if int32(wxLoginResp.ErrCode) == 0 {
	//			player.LoginRespStatus = 1
	//		} else {
	//			player.LoginRespStatus = int32(wxLoginResp.ErrCode)
	//		}
	//		player.SessionKey = wxLoginResp.SessionKey
	//		player.Attr.Channel = msg.Channel // channel 未入库，V1.0.1版本 后期可以删除
	//	} else if msg.Channel == "oppo" {
	//		acctInfo, b := sdk.SdkUtil.GetSDKUserInfo(sdk.SDKChannel(msg.Channel), "", msg.Token)
	//		if acctInfo == nil || !b {
	//			return -1, nil, 0
	//		}
	//		userList := db.GetUserByOpenId(acctInfo.Name)
	//		if len(userList) == 0 { // 不存在账号，则创建角色
	//			//if tableconfig.SenWordConfigs.CheckSenWord(msg.Username) {
	//			//	player.ErrorResponse(proto3.ErrEnum_Error_Involving_SenWord, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_Involving_SenWord)])
	//			//	return -1, nil, 0
	//			//}
	//			//// 校验名字重复
	//			//if db.JudgeUserName(msg.Username) {
	//			//	player.ErrorResponse(proto3.ErrEnum_Error_UserName_Exists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_UserName_Exists)])
	//			//	return -1, nil, 0
	//			//}
	//			player.CreateUser(msg.Username, ip, msg.Sex, acctInfo.Name, msg.Channel)
	//		} else { // 存在账号
	//			attr := db.GetUserByUserId(userList[0])
	//			if attr != nil {
	//				player.Attr = attr
	//				player.ConstructPlayer(attr.UserID)
	//			}
	//			player.Attr.Channel = msg.Channel // channel 未入库，V1.0.1版本
	//		}
	//	} else { // 无渠道登录
	//		attr := db.GetUserByUserName(msg.Username)
	//		if attr == nil {
	//			player.CreateUser(msg.Username, ip, msg.Sex, "", "")
	//		} else {
	//			player.Attr = attr
	//			LoadUserTitle(player.Attr)
	//			LoadUserAchievement(player.Attr)
	//			player.ConstructPlayer(attr.UserID)
	//		}
	//		player.LoginRespStatus = 1
	//	}
	//*/
	//// 创建协程
	//pidName := util.ToStr(player.Attr.UserID)
	//playerPid = global.RegisterPid(pidName, 256, player)
	//player.Pid = playerPid
	//playerPid.SendAfter(include.PlayerSaveDbData, include.PlayerSaveDbData, saveUserTime*1000, nil) // 定时保存玩家数据
	return 0, playerPid, player.Attr.UserID
}

func (p *Player) Start() {
	rand.Seed(time.Now().UnixNano())
	global.GloInstance.AddPlayer(p.Attr.UserID, p)
}

func (p *Player) HandleCall(req global.GenReq) global.Reply {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("player SOCKET_EVENT HandleCall error:%v, userID:%v, method: %v, msg: %v, stack: %v\n ",
				err, p.Attr.UserID, req.Method, req.MsgData, string(debug.Stack()))
			logger.Log.Errorf("player HandleCall panic-------------player:%v", p)
		}
	}()
	switch req.Method {
	case "":

	}
	return nil
}

func (p *Player) HandleCast(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("player SOCKET_EVENT HandleCast error:%v\n, userID:%v\n, method: %v\n, msg: %v\n, stack: %v\n ",
				err, p.Attr.UserID, req.Method, req.MsgData, string(debug.Stack()))
			logger.Log.Errorf("player HandleCast panic-------------player:%v", p)
		}
	}()
	switch req.Method {
	case "SOCKET_EVENT":
		msg := req.MsgData.(*Message)
		Handler.Callback(msg.Cmd, msg.PbData, p)
	case "changeAttr":
		//msg := req.MsgData.(*Message)
		//resource := msg.PbData.([]*dataConfig.CfgCost)
		//p.AddAttr(resource, false)
	case "enterRoomSuccess":
		pbData := req.MsgData.(*proto3.EnterRoomResp)
		// fmt.Println("enter room success : ", pbData.BirthPoint)
		cmd := proto3.ProtoCmd_CMD_EnterRoomResp

		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	case "finishTaskResp":
		cmd := proto3.ProtoCmd_CMD_FinishMissionResp
		pbData := req.MsgData.(*proto3.FinishMissionResp)
		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	case "dropItemResp":
		dropItemMap := req.MsgData.(map[int32][]*proto3.Item)

		cmd := proto3.ProtoCmd_CMD_DropItemResp
		pbData := &proto3.DropItemResp{}
		for k, v := range dropItemMap {
			pbData.DropItem = append(pbData.DropItem, &proto3.DropItem{Type: proto3.DropTypeEnum(k), Items: v})
		}
		// TODO 同类型合并
		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
	case "choiceItemResp":
		msgData := req.MsgData.(*proto3.ChoiceItemResp)

		cmd := proto3.ProtoCmd_CMD_ChoiceItemResp
		p.SendMessage(&Message{Cmd: cmd, PbData: msgData})
		logger.Log.Info("ChoiceItemResp end: ", msgData)
	case "other":
		fmt.Println("other")
	case "quitRoom":
		cmd := proto3.ProtoCmd_CMD_PlayerExitResp
		pbData := &proto3.PlayerExitResp{
			ErrCode: proto3.ErrEnum_Error_Pass,
		}
		p.SendMessage(&Message{Cmd: cmd, PbData: pbData})
		logger.Log.Infof("player :%d", p.Attr.UserID)
		p.ExitRoom()
	default:
		logger.Log.Errorln("playergoroutine HandleCast error: ", p.Attr.UserID, req.Method, req.MsgData)
	}
}

func (p *Player) HandleInfo(req global.GenReq) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorf("player SOCKET_EVENT HandleInfo error:%v\n, userID:%v\n, method: %v\n, msg: %v\n, stack: %v\n ",
				err, p.Attr.UserID, req.Method, req.MsgData, string(debug.Stack()))
			logger.Log.Errorf("player HandleInfo panic-------------player:%v", p)
		}
	}()
	switch req.Method {
	case include.PlayerAddBasicAttr:
		//p.CaleResourceGrowUp()
	case include.PlayerSaveDbData:
		p.SaveDBData()
		p.Pid.SendAfter(include.PlayerSaveDbData, include.PlayerSaveDbData, saveUserTime*1000, nil)
	case "player_stop":
		p.Pid.CastStop()
	default:
		logger.Log.Errorln("playergoroutine HandleInfo error: ", p.Attr.UserID, req.Method, req.MsgData)
	}
}

func (p *Player) Terminate() {
	logger.Log.Warnf("player pid terminated,save DB data!!! id:%d, userName:%s", p.Attr.UserID, p.Attr.Username)
	global.GloInstance.DelPlayer(p.Attr.UserID)
	p.IsOffLine = 1
	close(p.WriteChan)
	p.SaveDBData()
	if p.RoomPid != nil {
		p.RoomPid.Cast("playerOffline", p.Attr.UserID)
	}
}

func (p *Player) CheckWaiting() bool {
	if p.Room == nil || p.RoomPid == nil {
		return false
	}
	return p.Room.CheckRoomGameStatus(proto3.RoomStatus_wait_game)
}

func (p *Player) CheckGaming() bool {
	if p.Room == nil || p.RoomPid == nil {
		return false
	}
	return p.Room.CheckRoomGameStatus(proto3.RoomStatus_gameing)
}

func (p *Player) CheckVoting() bool {
	if p.Room == nil || p.RoomPid == nil {
		return false
	}
	return p.Room.CheckRoomGameStatus(proto3.RoomStatus_voting)
}

func (p *Player) CheckPlayerGameStatus(status proto3.PlayerGameStatus) bool {
	if p.Room == nil || p.RoomPid == nil {
		return false
	}

	return p.Room.CheckPlayerGameStatus(p.Attr.UserID, status)
}
