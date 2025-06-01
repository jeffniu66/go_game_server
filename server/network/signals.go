package network

import (
	"fmt"
	"go_game_server/server/game"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/sdk"
	"go_game_server/server/tableconfig"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	SIGRTMIN = syscall.Signal(0x22) // sign real time min 34
)

func Signals(cls chan os.Signal) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGINT+28, syscall.SIGINT+29, SIGRTMIN+1,
		SIGRTMIN+2, SIGRTMIN+3, SIGRTMIN+4, SIGRTMIN+5, SIGRTMIN+6, SIGRTMIN+7, SIGRTMIN+8, SIGRTMIN+9)
	logger.Log.Infof("服务进程 pid = %d\n", syscall.Getpid())
	_ = savePid(syscall.Getpid())

	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			logger.Log.Errorln("signals receive err: ", err, string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
			Signals(cls)
		}
	}()

	for {
		sig := <-sigs
		switch sig {
		case syscall.SIGINT:
			// Server.Connect.Close() // 这里要先关闭服务器端口
			logger.Log.Info("sig=", sig)
			_ = os.Remove("../script/pid.txt")
			logger.Log.Info(">>>>>>>>>>>> 服务器清理玩家, 延时3秒：<<<<<<<<<<<<<<")
			KickAllPlayer()
			time.Sleep(3 * time.Second)
			// close(Server.Done)
		case syscall.SIGINT + 28: // 热更配置 windows兼容修改
			// 服务器配置
			logger.Log.Info("更新ServerConfig")
			global.MyConfig.InitConfig("./config/pro_server.config")
			logLevel := global.MyConfig.Read("logger", "level")
			id, err := strconv.Atoi(logLevel)
			if err != nil {
				logger.Log.Errorln(err)
			}
			logger.InitLogger(logrus.Level(id))
			sdk.NewSDKUtil()
		case syscall.SIGINT + 29:
			logger.Log.Info("更新tableConfigs")
			// 读策划表数据
			if err := tableconfig.ReadTable(); err != nil {
				logger.Log.Errorf("read table failed, err:%v", err)
				return
			}
		case SIGRTMIN + 1:
			global.OrderConfig.OpenOrder(global.OrderGm, true)
		case SIGRTMIN + 2:
			global.OrderConfig.OpenOrder(global.OrderGm, false)
		case SIGRTMIN + 3:
			global.OrderConfig.OpenOrder(global.OrderWhite, true)
		case SIGRTMIN + 4:
			global.OrderConfig.OpenOrder(global.OrderWhite, false)
		case SIGRTMIN + 5:
			fmt.Println("kill all player")
			KickAllPlayer()
		case SIGRTMIN + 6: // warning level
			logger.Log.SetLevel(3)
		case SIGRTMIN + 7: // info level
			logger.Log.SetLevel(4)
		case SIGRTMIN + 8:
			global.OrderConfig.OpenOrder(global.OrderCountry, true)
		case SIGRTMIN + 9:
			global.OrderConfig.OpenOrder(global.OrderCountry, false)
		default:
			fmt.Println("unknown signal.")
		}
		cls <- sig
	}
}

func savePid(pid int) error {
	fileName := "../script/pid.txt"
	_ = os.Remove(fileName)
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_EXCL, os.ModePerm) // 如果已经存在,则失败
	defer file.Close()
	if err != nil {
		return err
	}
	_, _ = file.WriteString(strconv.Itoa(pid))
	return nil
}

func KickAllPlayer() {
	for _, userID := range global.GloInstance.GetPlayerIDList() {
		player := global.GloInstance.GetPlayer(userID)
		if player != nil {
			p := player.(*game.Player)
			if p.Pid == nil {
				logger.Log.Info("KickAllPlayer p.Pid is nil")
				continue
			}
			fmt.Println("userID = ", p.Attr.UserID)
			p.Pid.Stop()
		}
	}
}
