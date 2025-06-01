package tableconfig

type ActionConfig struct {
	Id   int32 `json:"id"`
	Type int32 `json:"type"`
}

type ActionConfigCol struct {
	ActionList []ActionConfig
	ActionMap  map[int32]ActionConfig
}

var ActionConfigs *ActionConfigCol

func (t *ActionConfigCol) InitMap() {
	t.ActionMap = make(map[int32]ActionConfig)
	for _, v := range t.ActionList {
		t.ActionMap[v.Id] = v
	}
}

func (t *ActionConfigCol) GetActionById(id int32) ActionConfig {
	action, ok := t.ActionMap[id]
	if !ok {
		return ActionConfig{}
	}
	return action
}
