package tableconfig

type StoreConfig struct {
	Id      int32 `json:"id"`
	Itemid  int32 `json:"itemid"`
	Salenum int32 `json:"salenum"`
	Type    int32 `json:"type"`
	Price   int32 `json:"price"`
}

type StoreConfigCol struct {
	StoreList []StoreConfig
	StoreMap  map[int32]StoreConfig
}

var StoresConfigs *StoreConfigCol

func (t *StoreConfigCol) InitMap() {
	t.StoreMap = make(map[int32]StoreConfig)
	for _, v := range t.StoreList {
		t.StoreMap[v.Id] = v
	}
}

func (t *StoreConfigCol) GetStoreById(id int32) StoreConfig {
	store, ok := t.StoreMap[id]
	if !ok {
		return StoreConfig{}
	}
	return store
}
