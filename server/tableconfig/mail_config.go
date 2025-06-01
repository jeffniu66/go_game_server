package tableconfig

type MailConfig struct {
	ID        int32  `jsin:"id"`
	Name      string `json:"name"`
	Desc      string `json:"desc"`
	Rewards   string `json:"rewards"`
	UseSwitch string `json:"useing"`
}

type MailConfigCol struct {
	MailConfigList []MailConfig
	MailConfigMap  map[int32]*MailConfig
}

func (m *MailConfigCol) InitMap() {
	m.MailConfigMap = make(map[int32]*MailConfig, 0)
	for i := 0; i < len(m.MailConfigList); i++ {
		v := &m.MailConfigList[i]
		m.MailConfigMap[v.ID] = v
	}
}

func (m *MailConfigCol) GetMail(ID int32) *MailConfig {
	if v, ok := m.MailConfigMap[ID]; !ok {
		return nil
	} else {
		return v
	}
}

var MailConfigs *MailConfigCol
