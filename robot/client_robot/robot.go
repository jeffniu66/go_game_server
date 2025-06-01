package client_robot

import (
	"bufio"
	"fmt"
	"math/rand"
	"net"
	"os"
	"runtime/debug"
	"go_game_server/proto3"
	"go_game_server/server/util"
	"time"
)

type RobotCfg struct {
	Address     string
	Name        string
	Num         int32 // 机器人数量
	Country     int32
	RobotTime   time.Duration
	LotteryTime time.Duration
	GlobalTime  time.Duration
	LastTime    int32 // 抽奖持续时间
	Attack      int32
}

var RobotInfo *RobotCfg

var RobotMap map[int]*ClientRobot

type robotData struct {
	Cmd    proto3.ProtoCmd
	PbData interface{}
}

type ClientRobot struct {
	CurIndex   int32
	UserID     int32
	BirthPoint int32
	Country    int32
	LordName   string
	AllyID     int32
	Socket     net.Conn
	MsgList    []robotData
}

func Start() {
	rand.Seed(time.Now().UnixNano())
	RobotInfo = &RobotCfg{}
	input := bufio.NewScanner(os.Stdin)

	//RobotInfo.Num = 100
	//RobotInfo.LotteryTime = time.Duration(rand.Intn(10) * 60)
	RobotInfo.Address = "170.106.66.44:10008"
	//fmt.Println("请输入ip和端口，例如192.168.1.187:8005(默认170.106.66.44:10008)")
	//input.Scan()
	//address := input.Text()
	//if address == "" {
	//	RobotInfo.Address = "127.0.0.1:8001"
	//} else {
	//	RobotInfo.Address = address
	//}

	RobotInfo.Num = 600
	//fmt.Println("请输入机器人数量：")
	//input.Scan()
	//robotNum := input.Text()
	//if robotNum == "" {
	//	RobotInfo.Num = 50
	//} else {
	//	RobotInfo.Num = util.ToInt(robotNum)
	//}

	RobotInfo.LotteryTime = time.Duration(util.RandInt(1, 5) * 60)

	fmt.Println("请输入国家选择：")
	input.Scan()
	num := input.Text()
	if num == "" {
		RobotInfo.Country = 2
	} else {
		RobotInfo.Country = util.ToInt(num)
	}
	RobotInfo.Attack = 1
	//fmt.Println("是否打地：输入1为打地")
	//input.Scan()
	//attack := input.Text()
	//if attack == "" {
	//	RobotInfo.Attack = 0
	//} else {
	//	RobotInfo.Attack = util.ToInt(attack)
	//}
	//
	//fmt.Println("请输入机器人几分钟后开始抽奖，例如5")
	//input.Scan()
	//minute := input.Text()
	//if minute == "" {
	//	RobotInfo.LotteryTime = time.Duration(rand.Intn(10) * 60)
	//} else {
	//	t := time.Duration(util.ToInt(minute))
	//	RobotInfo.LotteryTime = t * 60
	//}
	RobotInfo.LotteryTime = time.Duration(rand.Intn(10)) * 30

	RobotInfo.Name = "cli_robot"
	//RobotInfo.Country = 4
	RobotInfo.RobotTime = time.Duration(rand.Intn(5) + 3)
	RobotInfo.GlobalTime = 25
	RobotInfo.LastTime = 300

	fmt.Printf("------------欢迎进入大保健机器人测试------------\n"+
		"机器人名字：%v\n"+
		"机器人数量：%v\n"+
		"机器人思考时间：%v\n"+
		"机器人%v后开始抽奖\n", RobotInfo.Name, RobotInfo.Num, RobotInfo.RobotTime, RobotInfo.LotteryTime)
	fmt.Println("测试项： 登录 - 打开邮件 - 加入同盟 - gm指令 - 地块列表 - 出征 - 抽卡")
	f := func(t int) {
		for i := 0; i <= t; i++ {
			time.Sleep(1 * time.Second)
			fmt.Println("倒计时：", t-i)
		}
	}
	f(5)

	RobotMap = make(map[int]*ClientRobot)
	go globalTimer()

	for i := 0; i < int(RobotInfo.Num); i++ {
		time.Sleep(300 * time.Millisecond)
		fmt.Printf("----------------- address:%v", RobotInfo.Address)
		conn, _ := net.Dial("tcp", RobotInfo.Address)
		c := &ClientRobot{Socket: conn}
		go c.robotTimer(i)
		RobotMap[i] = c
		go c.receive(conn)
		c.robotLogin(i) //请求登录
		//go c.robotLogin(i) // TODO 测试并发登录

	}
	clientServer := make(chan struct{})
	<-clientServer
}

func (c *ClientRobot) receive(socket net.Conn) {
	defer func() { // 必须要先声明defer，否则不能捕获到panic异常
		if err := recover(); err != nil {
			fmt.Println("tcp_server receive packet err: ", err, string(debug.Stack())) // 这里的err其实就是panic传入的内容，55
			debug.PrintStack()
		}
	}()
	for {
		cmd, _, pbData, _ := RecvLTVPacket(socket, c)
		if pbData != nil {
			Handler.Callback(cmd, pbData, c)
		}
		if cmd == proto3.ProtoCmd_CMD_PASS {
			fmt.Println("没有心跳，断开连接")
			return
		}
	}
}

func (c *ClientRobot) robotTimer(i int) {
	ticker := time.NewTicker(RobotInfo.RobotTime * time.Second)
	for {
		select {
		case <-ticker.C:
			if len(c.MsgList) > 0 {
				msgData := c.MsgList[0]
				SendLTVPacket(c.Socket, msgData.Cmd, msgData.PbData)
				c.MsgList = c.MsgList[1:]
			}
		}
	}
}

func globalTimer() {
	ticker := time.NewTicker(RobotInfo.GlobalTime * time.Second)
	for {
		select {
		case <-ticker.C:

		}
	}
}
