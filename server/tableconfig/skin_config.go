package tableconfig

type SkinConfig struct {
	Id         int32 `json:"id"`
	Storetype  int32 `json:"storetype"`
	Dropid     int32 `json:"dropid"`
	Skinid     int32 `json:"skinid"`
	Collectnum int32 `json:"collectnum"`
	Type       int32 `json:"type"`
	Limit      int32 `json:"limit"`
	Price      int32 `json:"price"`
}

type SkinConfigCol struct {
	SkinList    []SkinConfig
	SkinMap     map[int32]*SkinConfig
	SkinTypeMap map[int32][]SkinConfig
}

var SkinConfigs *SkinConfigCol

func (t *SkinConfigCol) InitMap() {
	t.SkinMap = make(map[int32]*SkinConfig)
	t.SkinTypeMap = make(map[int32][]SkinConfig)
	for _, v := range t.SkinList {
		t.SkinMap[v.Id] = &v

		typeList := t.SkinTypeMap[v.Storetype]
		if typeList == nil {
			typeList = []SkinConfig{}
		}
		typeList = append(typeList, v)
		t.SkinTypeMap[v.Storetype] = typeList
	}
}

func (t *SkinConfigCol) GetSkinById(id int32) *SkinConfig {
	skin, ok := t.SkinMap[id]
	if !ok {
		return nil
	}
	return skin
}
