package server_robot

// u3d选择3国家进入,然后3,4,5国家分别添加198,149,149人
//func TestAddCountryNum(num, country int32) {
//	logger.Log.Warnln("国家已出生人数: ", db.CountryNumMap)
//	db.SetCountryNum(country, num)
//	//db.CountryNumMap[country] += num
//	logger.Log.Warnln("设置后国家已出生人数: ", db.CountryNumMap)
//}

// 测试流程
//func TestAddRobotNum(num, country int32) {
//	for i := 0; i < int(num); i++ {
//		acctName := "robot" + util.GetRandName()
//		writeChan := make(chan interface{}, 1024)
//		playerPid, userID := game.CreatePid(writeChan, acctName, "192.168.1.1", nil)
//		robot := &game.ServerRobot{Pid: playerPid, UserID: userID, AcctName: acctName}
//		robot.CreateCountry(country)
//	}
//}
