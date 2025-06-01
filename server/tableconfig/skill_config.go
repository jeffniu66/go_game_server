package tableconfig

type SkillConfig struct {
	Id        int32  `json:"id"`
	Skilltype int32  `json:"skilltype"`
	Link      int32  `json:"link"`
	Active    int32  `json:"active"`
	Target    int32  `json:"target"`
	Time      string `json:"time"`
	Text      string `json:"text"`
	Num       int32  `json:"num"`
}

type SkillConfigCol struct {
	SkillList []SkillConfig
	SkillMap  map[int32]SkillConfig
}

var SkillConfigs *SkillConfigCol

func (t *SkillConfigCol) InitMap() {
	t.SkillMap = make(map[int32]SkillConfig)
	for _, v := range t.SkillList {
		t.SkillMap[v.Id] = v
	}
}

func (t *SkillConfigCol) GetSkillById(id int32) SkillConfig {
	skill, ok := t.SkillMap[id]
	if !ok {
		return SkillConfig{}
	}
	return skill
}
