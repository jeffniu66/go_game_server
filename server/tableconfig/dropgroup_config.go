package tableconfig

type DropGroupConfig struct {
	Id          int `json:"id"`
	DropId      int `json:"dropid"`
	ItemType    int `json:"itemtype"`
	ItemId      int `json:"itemid"`
	Num         int `json:"num"`
	Probability int `json:"probability"`
}

type DropGroupConfigCol struct {
	DropGroupList []DropGroupConfig
	DropMap       map[int]DropGroupConfig
	DropGroupMap  map[int][]DropGroupConfig // key: DropId
}

func (t *DropGroupConfigCol) InitMap() {
	t.DropMap = make(map[int]DropGroupConfig)
	t.DropGroupMap = make(map[int][]DropGroupConfig)
	for _, v := range t.DropGroupList {
		t.DropMap[v.Id] = v

		groupMap := t.DropGroupMap[v.DropId]
		if groupMap == nil {
			groupMap = []DropGroupConfig{}
		}
		groupMap = append(groupMap, v)
		t.DropGroupMap[v.DropId] = groupMap
	}
}

func (t *DropGroupConfigCol) GetDropGroupById(id int) *DropGroupConfig {
	if _, ok := t.DropMap[id]; !ok {
		return nil
	}
	ret := t.DropMap[id]
	return &ret
}

func (t *DropGroupConfigCol) GetDropGroupMap(dropGroupId int) []DropGroupConfig {
	return t.DropGroupMap[dropGroupId]
}

var DropGroupConfigs *DropGroupConfigCol
