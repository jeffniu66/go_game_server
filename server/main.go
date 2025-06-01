package main

import (
	"fmt"
	_ "go_game_server/config/dataConfig"
	"go_game_server/server/db"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/goredis"
	_ "go_game_server/server/handler"
	"go_game_server/server/logger"
	"go_game_server/server/network"
	"go_game_server/server/sdk"
	"go_game_server/server/tableconfig"
	_ "net/http/pprof"
	"strconv"

	"github.com/sirupsen/logrus"
)

func main() {
	// 服务器配置
	global.InitServerConfig()
	global.NewGlobalPlayers()

	logLevel := global.MyConfig.Read("logger", "level")
	id, err := strconv.Atoi(logLevel)
	if err != nil {
		logger.Log.Errorln(err)
	}
	logger.InitLogger(logrus.Level(id))
	logger.Log.Infof(">>>>>>>>>>>> 服务器正在启动中 ... <<<<<<<<<<<<<<\n\n")
	sdk.NewSDKUtil()

	// 读策划表数据
	if err := tableconfig.ReadTable(); err != nil {
		logger.Log.Errorf("read table failed, err:%v", err)
		return
	}

	// 初始化数据库连接池
	db.InitMysql()
	db.NewIncrease()
	goredis.InitRedis()
	game.InitMatchMgr() // 单独协程用于匹配管理器
	game.LoadNames()
	game.NewGlobalTimer()

	go game.InitTopBoard() // 初始化排行榜
	go network.StartHttpServer()

	startServer()
	logger.Log.Info(">>>>>>>>>>>> 服务器已关闭 ... <<<<<<<<<<<<<<\n\n")
}

func startServer() {
	netType := global.MyConfig.Read("server", "net_type")
	switch netType {
	case "tcp":
		network.StartTcpServer()
	case "websocket":
		network.StartWsServer()
	default:
		fmt.Println("net work type err ", netType)
	}
}
