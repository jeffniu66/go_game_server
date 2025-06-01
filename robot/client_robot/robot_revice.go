package client_robot

import (
	"fmt"
	"go_game_server/proto3"
)

func InitReceive() {
	Handler.RegistHandler(proto3.ProtoCmd_CMD_ErrResp, &proto3.ErrResp{}, receiveErr)
	Handler.RegistHandler(proto3.ProtoCmd_CMD_LoginResp, &proto3.LoginResp{}, receiveLogin)
	// Handler.RegistHandler(proto3.ProtoCmd_CMD_CreateNameResp, &proto3.CreateNameResp{}, receiveCreateName)
	// Handler.RegistHandler(proto3.ProtoCmd_CMD_CreateCountryResp, &proto3.CreateCountryResp{}, receiveCreateCountry)

}

func receiveLogin(req interface{}, robot *ClientRobot) interface{} {
	// msg := req.(*proto3.LoginResp)
	// robot.BirthPoint = msg.PlayerAttr.BirthPoint
	// robot.AllyID = msg.PlayerAttr.AllyID
	// robot.UserID = msg.PlayerAttr.UserID
	if robot.BirthPoint == 0 { // 新账号
		robot.createName()
	} else {
		// robot.LordName = msg.PlayerAttr.LordName
		fmt.Printf("old机器人名字：%v, 机器人出生点：%v\n", robot.LordName, robot.BirthPoint)
		robot.robotRun()
	}
	return nil
}

func receiveCreateName(req interface{}, robot *ClientRobot) interface{} {
	fmt.Printf("new机器人名字：%v, 新创建的机器人出生点：%v\n", robot.LordName, robot.BirthPoint)
	robot.createCountry()
	return nil
}

func receiveCreateCountry(req interface{}, robot *ClientRobot) interface{} {
	msg := req.(*proto3.CreateCountryResp)
	robot.BirthPoint = msg.BirthPoint
	fmt.Printf("新创建的机器人名字：%v, 新创建的机器人出生点：%v\n", robot.LordName, robot.BirthPoint)
	robot.robotRun()
	return nil
}

// 出征返回错误，继续打
func receiveErr(req interface{}, robot *ClientRobot) interface{} {
	// msg := req.(*proto3.ErrResp)
	// switch msg.Cmd {
	// case proto3.ProtoCmd_CMD_WorldLandReq:
	// 	// 继续出征
	// 	fmt.Printf("机器人名字：%v, 出征失败原因：%v, 继续出征\n", robot.LordName, msg.ErrCode)
	// 	robot.testStep10()
	// case proto3.ProtoCmd_CMD_CreateNameReq:
	// 	robot.createName()
	// default:
	// 	fmt.Println("协议返回报错：", msg.Cmd, msg.ErrCode, msg.ErrMsg)
	// }
	return nil

}
