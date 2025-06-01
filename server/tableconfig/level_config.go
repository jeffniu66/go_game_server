package tableconfig

type LevelConfig struct {
	Level   int32  `json:"id"`
	NeedNum int32  `json:"neednum"`
	Rewards string `json:"reward"`
}

type LevelConfigCol struct {
	LevelConfigList []LevelConfig
	LevelConfigMap  map[int32]*LevelConfig // level连续等级
}

func (l *LevelConfigCol) GetLevelConfig(level int32) *LevelConfig {
	if v, ok := l.LevelConfigMap[level]; ok {
		return v
	}
	return nil
}
func (l *LevelConfigCol) InitMap() {
	l.LevelConfigMap = make(map[int32]*LevelConfig, 0)
	for i := 0; i < len(l.LevelConfigList); i++ {
		v := &l.LevelConfigList[i]
		l.LevelConfigMap[v.Level] = v
	}
}

func (l *LevelConfigCol) GetNextLevel(level, exp int32) (retLevel, retExp, maxExp int32) {
	if level < 0 {
		return
	}
	// 最大等级不升级
	if int(level) == len(l.LevelConfigMap) {
		retLevel = level
		retExp = exp
		maxExp = l.LevelConfigMap[level].NeedNum
		return
	}

	i := level
	if i == 0 {
		i = 1
	}
	for {
		levelInfo, ok := l.LevelConfigMap[i]
		if !ok {
			break
		}
		if exp < levelInfo.NeedNum {
			retLevel = levelInfo.Level
			maxExp = levelInfo.NeedNum
			retExp = exp
			break
		} else {
			exp -= levelInfo.NeedNum
		}
		i++
		if i > 1000 {
			break
		}
	}
	return
}

var LevelConfigs *LevelConfigCol
