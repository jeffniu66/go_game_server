package tableconfig

import (
	"go_game_server/server/util"
	"strings"
)

type AvatarConfig struct {
	ID         int32  `json:"id"`
	AvatarType int32  `json:"type"`
	UnLockType int32  `json:"unlocakType"`
	Items      string `json:"num"`
	ReturnNum  string `json:"returnnum"`
	ItemId     int32
	NeedNum    int32
}
type AvatarConfigCol struct {
	AvatarConfigList []AvatarConfig
	AvatarConfigMap  map[int32]*AvatarConfig
}

func (a *AvatarConfigCol) GetConfig(ID int32) *AvatarConfig {
	if v, ok := a.AvatarConfigMap[ID]; ok {
		return v
	}
	return nil
}
func (a *AvatarConfigCol) InitMap() {
	a.AvatarConfigMap = make(map[int32]*AvatarConfig, 0)
	for i := range a.AvatarConfigList {
		v := &a.AvatarConfigList[i]
		a.AvatarConfigMap[v.ID] = v
		items := strings.Split(v.Items, ",")
		if len(items) == 2 {
			v.ItemId = util.ToInt(items[0])
			v.NeedNum = util.ToInt(items[1])
		}
	}
}

func (a *AvatarConfigCol) GetEnoughItem(itemId, num int32) int32 {
	for i := range a.AvatarConfigList {
		v := a.AvatarConfigList[i]
		if v.ItemId == itemId && v.NeedNum <= num {
			return v.ID
		}
	}
	return -1
}

var AvatarConfigs *AvatarConfigCol
