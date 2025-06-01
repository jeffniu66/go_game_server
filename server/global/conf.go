package global

import (
	"bufio"
	"flag"
	"fmt"
	"go_game_server/proto3"
	"go_game_server/server/logger"
	"go_game_server/server/util"
	"io"
	"os"
	"strings"
)

const middle = "========="

type configf struct {
	mymap  map[string]string
	strcet string
}

type OrderCfg struct {
	GMMap       map[string]bool
	WhiteMap    map[string]bool
	OpenGM      bool
	OpenWhite   bool
	OpenCountry bool
	GmIndex     int32 // 所有机器人打这块地
}

const (
	OrderGm = iota
	OrderWhite
	OrderCountry
)

var MyConfig *configf
var OrderConfig *OrderCfg

func InitServerConfig() {
	MyConfig = new(configf)
	ServerPort = flag.String("port", "", "端口")
	Game2Battle = flag.String("game2battle", "", "game2battle rabbitmq队列名称")
	Battle2Game = flag.String("battle2game", "", "battle2game rabbitmq队列名称")
	flag.Parse()
	if *ServerPort == "" {
		MyConfig.InitConfig("./config/pro_server.config")
		port := MyConfig.Read("server", "address") // 服务器地址端口
		game2battle := MyConfig.Read("kafka", "game_2_battle")
		battle2game := MyConfig.Read("kafka", "battle_2_game")
		*ServerPort = port
		*Game2Battle = game2battle
		*Battle2Game = battle2game
	} else {
		MyConfig.InitConfig("./config/dev_server.config")
		//portStr := strings.Split(*ServerPort, ":")
		//ServerNum = util.ToInt(portStr[1])
	}
	fmt.Println("init config success! ", *ServerPort, *Battle2Game, *Game2Battle, ServerNum)
	OrderConfig = &OrderCfg{OpenGM: MyConfig.ReadBool("order", "open_GM"),
		OpenWhite: MyConfig.ReadBool("order", "open_white"), OpenCountry: MyConfig.ReadBool("order", "open_country")}
	OrderConfig.InitConfig()
}

func (c *configf) InitConfig(path string) {
	c.mymap = make(map[string]string)

	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		b, _, err := r.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			panic(err)
		}

		s := strings.TrimSpace(string(b))
		//fmt.Println(s)
		if strings.Index(s, "#") == 0 {
			continue
		}

		n1 := strings.Index(s, "[")
		n2 := strings.LastIndex(s, "]")
		if n1 > -1 && n2 > -1 && n2 > n1+1 {
			c.strcet = strings.TrimSpace(s[n1+1 : n2])
			continue
		}

		if len(c.strcet) == 0 {
			continue
		}
		index := strings.Index(s, "=")
		if index < 0 {
			continue
		}

		frist := strings.TrimSpace(s[:index])
		if len(frist) == 0 {
			continue
		}
		second := strings.TrimSpace(s[index+1:])

		pos := strings.Index(second, "\t#")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " #")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, "\t//")
		if pos > -1 {
			second = second[0:pos]
		}

		pos = strings.Index(second, " //")
		if pos > -1 {
			second = second[0:pos]
		}

		if len(second) == 0 {
			continue
		}

		key := c.strcet + middle + frist
		c.mymap[key] = strings.TrimSpace(second)
	}
}

func (c configf) getKey(node, key string) string {
	return node + middle + key
}

func (c configf) Read(node, key string) string {
	k := c.getKey(node, key)
	v, found := c.mymap[k]
	if !found {
		return ""
	}
	return v
}

func (c configf) ReadInt32(node, key string) int32 {
	k := c.getKey(node, key)
	v, found := c.mymap[k]
	if !found {
		return 0
	}
	return util.ToInt(v)
}

func (c configf) ReadBool(node, key string) bool {
	k := c.getKey(node, key)
	v, found := c.mymap[k]
	if !found {
		return false
	}
	return util.ToBool(v)
}

//////////////////////////////////////////////////////////
func (o *OrderCfg) InitConfig() {
	o.GMMap = make(map[string]bool)
	o.WhiteMap = make(map[string]bool)

	f, err := os.Open("./config/order.config")
	if err != nil {
		return
	}

	defer f.Close()

	t := OrderWhite
	r := bufio.NewReader(f)
	for {
		a, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}

		str := string(a)
		if len(str) == 0 {
			continue
		}

		if "[gm_order]" == str {
			t = OrderGm
			continue
		} else if "[white_order]" == str {
			t = OrderWhite
			continue
		}

		if t == OrderGm {
			o.GMMap[str] = true
		} else if t == OrderWhite {
			o.WhiteMap[str] = true
		}
	}
}

func (o *OrderCfg) existGmOrder(name string) bool {
	if _, ok := o.GMMap[name]; ok {
		return true
	}

	return false
}

func (o *OrderCfg) existWhiteOrder(name string) bool {
	if _, ok := o.WhiteMap[name]; ok {
		return true
	}

	return false
}

func (o *OrderCfg) OpenOrder(orderType int, isOpen bool) {
	logger.Log.Warnln(">>>>>>>>>>>>>>>>>> OpenOrder ", orderType, isOpen)
	if orderType == OrderGm {
		o.OpenGM = isOpen
	} else if orderType == OrderWhite {
		o.OpenWhite = isOpen
	} else if orderType == OrderCountry {
		o.OpenCountry = isOpen
	}

	// reload
	o.InitConfig()
}

func (o *OrderCfg) IsOpen(orderType int) bool {
	switch orderType {
	case OrderGm:
		return o.OpenGM
	case OrderWhite:
		return o.OpenWhite
	}

	return false
}

func (o *OrderCfg) CanLogin(name string, ip string) bool {
	// 没开白名单都能登陆，开了白名单ip符合则可以进，或者有该账号
	if !o.OpenWhite {
		return true
	} else if o.existWhiteOrder(ip) || o.existWhiteOrder(name) {
		return true
	}

	return false
}

func (o *OrderCfg) CanUseGm(name string) bool {
	if (o.OpenGM && o.existGmOrder(name)) || MyConfig.ReadBool("order", "dev_GM") {
		return true
	}
	if len(name) >= 9 && name[:9] == "cli_robot" {
		return true
	}

	return false
}

// 屏蔽国内玩家
func (o *OrderCfg) CanCountry(name, ip string, language proto3.LanguageEnum) bool {
	// 开关
	if !o.OpenCountry || MyConfig.ReadBool("order", "dev_GM") {
		return true
	}
	if language != proto3.LanguageEnum_ChineseSimple {
		return true
	}
	logger.Log.Warnln("shield country ip or name: ", ip, name, o.WhiteMap)
	return o.existWhiteOrder(ip) || o.existWhiteOrder(name)
}

func (o *OrderCfg) InitConfig1(path string) {
	o.GMMap = make(map[string]bool)
	o.WhiteMap = make(map[string]bool)

	f, err := os.Open(path)
	if err != nil {
		return
	}

	defer f.Close()

	t := OrderWhite
	r := bufio.NewReader(f)
	for {
		a, _, c := r.ReadLine()
		if c == io.EOF {
			break
		}

		str := string(a)
		if len(str) == 0 {
			continue
		}

		if "[gm_order]" == str {
			t = OrderGm
			continue
		} else if "[white_order]" == str {
			t = OrderWhite
			continue
		}

		if t == OrderGm {
			o.GMMap[str] = true
		} else if t == OrderWhite {
			o.WhiteMap[str] = true
		}
	}
}
