package tableconfig

type NameZhConfig struct {
	ID      int32  `json:"id"`
	Prefix  string `json:"prefix"`
	Surname string `json:"surname"`
	Name    string `json:"name"`
	Sex     int32  `json:"type"`
}

type NameZhConfigCol struct {
	NameZhConfigList []NameZhConfig
	NameZhConfigMap  map[int32]*NameZhConfig
}

func (n *NameZhConfigCol) InitMap() {
	n.NameZhConfigMap = make(map[int32]*NameZhConfig, 0)
	for i := 1; i < len(n.NameZhConfigList); i++ {
		v := &n.NameZhConfigList[i]
		n.NameZhConfigMap[v.ID] = v
	}
}

var NameZhConfigs *NameZhConfigCol
