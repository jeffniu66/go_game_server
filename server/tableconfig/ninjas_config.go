package tableconfig

type NinjaConfig struct {
	Id         int32  `json:"id"`
	NinjaLevel int32  `json:"ninjalevel"`
	Prev       int32  `json:"prev"`
	Next       int32  `json:"next"`
	NinjaStar  int32  `json:"star"`
	NeedNum    int32  `json:"neednum"`
	Rewards    string `json:"reward"`
}
type NinjaConfigCol struct {
	StartID         int32
	EndID           int32
	NinjaConfigList []NinjaConfig
	NinjaConfigMap  map[int32]*NinjaConfig // key-string(level+star)
}

func (n *NinjaConfigCol) GetNinjaConfig(ninjaID int32) *NinjaConfig {
	v, ok := n.NinjaConfigMap[ninjaID]
	if !ok {
		return nil
	}
	return v
}

func (n *NinjaConfigCol) GetLevel(ninjaID, archivePoint int32) (retID, retPoint, maxPoint int32) {
	retID, retPoint = ninjaID, archivePoint
	count := 0
	for {
		count++
		if count > 1000 {
			return
		}
		conf, ok := n.NinjaConfigMap[retID]
		if !ok {
			return
		}
		maxPoint = conf.NeedNum
		if retID == conf.Next || conf.NeedNum < maxPoint {
			return
		}
		if retPoint >= conf.NeedNum {
			conf2, ok := n.NinjaConfigMap[conf.Next]
			if !ok {
				retID = ninjaID
				retPoint = archivePoint
				maxPoint = conf.NeedNum
			}

			retID = conf2.Id
			retPoint -= conf.NeedNum
			if retPoint <= 0 {
				retPoint = 0
			}
		}
	}
}

func (n *NinjaConfigCol) InitMap() {
	n.NinjaConfigMap = make(map[int32]*NinjaConfig, 0)
	for i := 0; i < len(n.NinjaConfigList); i++ {
		v := &n.NinjaConfigList[i]
		if v.Id == v.Prev {
			n.StartID = v.Id
		}
		if v.Id == v.Next {
			n.EndID = v.Id
		}
		n.NinjaConfigMap[v.Id] = v
	}
}

var NinjaConfigs *NinjaConfigCol
