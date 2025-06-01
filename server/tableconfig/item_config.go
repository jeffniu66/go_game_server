package tableconfig

type ItemConfig struct {
	Id            int32  `json:"id"`
	Name          string `json:"name"`
	Desc          string `json:"desc"`
	Icon          string `json:"icon"`
	Type          int32  `json:"type"`
	Quality       int32  `json:"quality"`
	Bag           int32  `json:"bag"`
	Use           int32  `json:"use"`
	Effect        string `json:"effect"`
	Uselimit      string `json:"uselimit"`
	Useinfo       string `json:"useinfo"`
	Usecount      int32  `json:"usecount"`
	Overlay       int32  `json:"overlay"`
	Skip          string `json:"skip"`
	Uselevel      int32  `json:"uselevel"`
	Jiaobiao      string `json:"jiaobiao"`
	Zidongshiyong int32  `json:"zidongshiyong"`
	Sell          int32  `json:"sell"`
	Addnum        string `json:"addnum"`
	OpenBox       string `json:"openbox"`
	Approach      string `json:"approach"`
}

type ItemConfigCol struct {
	ItemConfigList []ItemConfig
	ItemConfigMap  map[int32]ItemConfig
}

func (i *ItemConfigCol) InitMap() {
	i.ItemConfigMap = make(map[int32]ItemConfig)
	for _, v := range i.ItemConfigList {
		i.ItemConfigMap[v.Id] = v
	}
}

func (i *ItemConfigCol) GetItemConfigById(id int32) *ItemConfig {
	if _, ok := i.ItemConfigMap[id]; !ok {
		return nil
	}
	ret := i.ItemConfigMap[id]
	return &ret
}

var ItemConfigs *ItemConfigCol
