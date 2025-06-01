package tableconfig

import (
	"go_game_server/server/constant"
	"go_game_server/server/util"
	"strconv"
	"strings"
)

type ConstsConfig struct {
	Id        int32   `json:"id"`
	Value     string  `json:"value"`
	ValueList []int32 // 数组值
}

type ConstsConfigCol struct {
	ConstList []ConstsConfig          `json:"const_list"`
	ConstMap  map[int32]*ConstsConfig `json:"const_map"`
}

func (c *ConstsConfigCol) GetValueById(id int32) string {
	return c.ConstMap[id].Value
}

func (c *ConstsConfigCol) GetMatchGameParam() (int, int) {
	// 100是匹配赛的id
	v, ok := c.ConstMap[constant.MatchRoomType]
	defaultL, defaultP := 2, 2
	if !ok {
		return defaultL, defaultP
	}
	str := strings.Split(v.Value, ",")
	if len(str) != 2 {
		return defaultL, defaultP
	}
	langRen, err := strconv.Atoi(str[0])
	if err != nil {
		return defaultL, defaultP
	}
	pingMin, err := strconv.Atoi(str[1])
	if err != nil {
		return defaultL, defaultP
	}
	return langRen, pingMin
}

func (c *ConstsConfigCol) GetFPSNumParam() int32 {
	// 100是匹配赛的id
	defFPS := int32(66)
	v, ok := c.ConstMap[constant.ConstFPSNum]
	if !ok {
		return defFPS
	}

	t, err := strconv.Atoi(v.Value)
	if err != nil {
		return defFPS
	}
	return int32(1000 / t)
}

func (c *ConstsConfigCol) GetGameWaitTime() int32 {
	defaultT := int32(5)
	v, ok := c.ConstMap[constant.GameWaitTime]
	if !ok {
		return defaultT
	}
	r, err := strconv.Atoi(v.Value)
	if err != nil {
		return defaultT
	}
	return int32(r)
}

// delete
func (c *ConstsConfigCol) GetChatVoteTime() int32 {
	defaultT := int32(5)
	v, ok := c.ConstMap[constant.ChatVoteTime]
	if !ok {
		return defaultT
	}
	r, err := strconv.Atoi(v.Value)
	if err != nil {
		return defaultT
	}
	return int32(r)
} // delete

func (c *ConstsConfigCol) GetIdValue(ID int32) int32 {
	defaultT := int32(0)
	v, ok := c.ConstMap[ID]
	if !ok {
		return defaultT
	}
	r, err := strconv.Atoi(v.Value)
	if err != nil {
		return defaultT
	}
	return int32(r)
}

func (c *ConstsConfigCol) GetAttackFrozen() int32 {
	defaultT := int32(5)
	v, ok := c.ConstMap[constant.AttackFrozenTime]
	if !ok {
		return defaultT
	}
	r, err := strconv.Atoi(v.Value)
	if err != nil {
		return defaultT
	}
	return int32(r)
}

func (c *ConstsConfigCol) GetIdValueFloat(ID int32) float32 {
	defaultT := float32(5)
	v, ok := c.ConstMap[ID]
	if !ok {
		return defaultT
	}
	r := util.ToFloat64(v.Value, 32)
	return float32(r)
}

func (c *ConstsConfigCol) GetVoteSumEndTime() int32 {
	defaultT := int32(5)
	v, ok := c.ConstMap[constant.VoteSumEndTime]
	if !ok {
		return defaultT
	}
	r, err := strconv.Atoi(v.Value)
	if err != nil {
		return defaultT
	}
	return int32(r)
}

func (c *ConstsConfigCol) GetGameVoteTime() int32 {
	defaultT := int32(30)
	v, ok := c.ConstMap[constant.GameVoteTime]
	if !ok {
		return defaultT
	}
	r, err := strconv.Atoi(v.Value)
	if err != nil {
		return defaultT
	}
	return int32(r)
}

func (c *ConstsConfigCol) GetRobotChatTime() int32 {
	var t, def1, def2 int32 = 5, 5, 10
	v, ok := c.ConstMap[constant.RobotChatTime]
	if !ok {
		return t
	}
	vList := strings.Split(v.Value, ",")
	if len(vList) != 2 {
		return t
	}
	def1 = util.ToInt(vList[0])
	def2 = util.ToInt(vList[1])

	t = util.RandInt(def1, def2)
	return t
}

func (c *ConstsConfigCol) GetRandValue(constID int32) int32 {
	cc, ok := c.ConstMap[constID]
	if !ok {
		return 0
	}
	index := util.RandInt(0, int32(len(cc.ValueList))-1)
	return cc.ValueList[index]
}

func (c *ConstsConfigCol) InitMap() {
	c.ConstMap = make(map[int32]*ConstsConfig)
	for i := range c.ConstList {
		v := &c.ConstList[i]
		if v.Id == constant.RandSkin || v.Id == constant.RandPhoto || v.Id == constant.RandTitle {
			c.ConstList[i].ValueList = make([]int32, 0)
			l := strings.Split(v.Value, ",")
			for _, vv := range l {
				if vv == "" {
					continue
				}
				d := util.ToInt(vv)
				c.ConstList[i].ValueList = append(c.ConstList[i].ValueList, d)
			}
		}
		c.ConstMap[v.Id] = v
	}
}

var ConstsConfigs *ConstsConfigCol
