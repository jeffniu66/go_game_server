package tableconfig

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
	"strings"
)

type TitleConfig struct {
	ID           int32  `json:"id"`
	MoudleType   int32  `json:"moudletype"`
	TitleType    int32  `json:"titletype"`
	Prev         int32  `json:"prev"`
	Next         int32  `json:"next"`
	NeedValue    int32  `json:"num"`
	Rewards      string `json:"rewards"`
	ArchivePoint int32
	Gold         int32
}
type TitleTypeAnnex struct {
	TitleConfigMap map[int32]*TitleConfig
	MaxValue       int32
	MinValue       int32
	MaxID          int32
	MinID          int32
}

type TitleConfigCol struct {
	TitleConfigList    []TitleConfig
	TitleConfigMap     map[int32]*TitleConfig              // key-id
	TitleAnnexTitleMap map[int32]map[int32]*TitleTypeAnnex // key1-modletype key2-titleType
}

type UserTitle struct {
	ID    int32
	IsGot int32 // include gotstatus 0-达成 已领取，1-达成，未领取，2-达成
}

// taType == include.TitleType include.ArchiveType
// GetTitle return 成就id，成就值，是否达成
func (t *TitleConfigCol) GetTitle(taType int32, titleType, value int32) []UserTitle {
	userTitle := make([]UserTitle, 0)
	if value < 0 {
		return userTitle
	}
	titleMap, ok := t.TitleAnnexTitleMap[taType]
	if !ok {
		return userTitle
	}

	titleAnnex, ok := titleMap[titleType]
	if !ok {
		return userTitle
	}

	curIndex := titleAnnex.MinID
	if _, ok := titleAnnex.TitleConfigMap[curIndex]; !ok {
		return userTitle
	}
	nextIndex := titleAnnex.TitleConfigMap[curIndex].Next
	for {
		// 当前等于下一个，达到最大称号
		if value < titleAnnex.TitleConfigMap[curIndex].NeedValue {
			break
		}
		if value >= titleAnnex.TitleConfigMap[curIndex].NeedValue {
			tmp := UserTitle{ID: curIndex, IsGot: include.GetStatusNo}
			userTitle = append(userTitle, tmp)
		}
		// 如果当前== 最大
		if curIndex == titleAnnex.MaxID {
			break
		}
		curIndex = nextIndex
		nextIndex = t.TitleConfigMap[curIndex].Next
	}
	return userTitle
}

func (t *TitleConfigCol) InitMap() {
	t.TitleAnnexTitleMap = make(map[int32]map[int32]*TitleTypeAnnex, 0)
	t.TitleConfigMap = make(map[int32]*TitleConfig, 0)
	titleMap := make(map[int32][]*TitleConfig, 0)   // key-titleTYpe
	archiveMap := make(map[int32][]*TitleConfig, 0) // key-titleTYpe
	for i := range t.TitleConfigList {
		tc := &t.TitleConfigList[i]
		if tc.MoudleType == include.TitleType {
			titleMap[tc.TitleType] = append(titleMap[tc.TitleType], tc)
		} else if tc.MoudleType == include.ArchiveType {
			archiveMap[tc.TitleType] = append(archiveMap[tc.TitleType], tc)
		}
		// 处理奖励
		rewards := strings.Split(tc.Rewards, "|") // 3,10|1,100
		for _, v := range rewards {
			values := strings.Split(v, ",")
			if len(values) == 2 {
				if util.ToInt(values[0]) == 1 {
					tc.Gold = util.ToInt(values[1])
				} else if util.ToInt(values[0]) == 3 {
					tc.ArchivePoint = util.ToInt(values[1])
				}
			}
		}
		t.TitleConfigMap[tc.ID] = tc
	}

	// 设置称号附加值
	titleTypeAnnex := make(map[int32]*TitleTypeAnnex, 0)
	for titleType := range titleMap {
		titleAnnex := &TitleTypeAnnex{}
		titleAnnex.TitleConfigMap = make(map[int32]*TitleConfig, 0)
		for ii := range titleMap[titleType] {
			title := titleMap[titleType][ii]
			titleAnnex.TitleConfigMap[title.ID] = title

			if title.ID == title.Prev {
				titleAnnex.MinValue = title.NeedValue
				titleAnnex.MinID = title.ID
			}
			if title.ID == title.Next {
				titleAnnex.MaxValue = title.NeedValue
				titleAnnex.MaxID = title.ID
			}
		}
		titleTypeAnnex[titleType] = titleAnnex
	}
	t.TitleAnnexTitleMap[include.TitleType] = titleTypeAnnex
	// 成就附加值
	archiveTypeAnnex := make(map[int32]*TitleTypeAnnex, 0)
	for archiveType := range archiveMap {
		archiveAnnex := &TitleTypeAnnex{}
		archiveAnnex.TitleConfigMap = make(map[int32]*TitleConfig, 0)
		for ii := range archiveMap[archiveType] {
			title := archiveMap[archiveType][ii]
			archiveAnnex.TitleConfigMap[title.ID] = title
			if title.ID == title.Prev {
				archiveAnnex.MinValue = title.NeedValue
				archiveAnnex.MinID = title.ID
			}
			if title.ID == title.Next {
				archiveAnnex.MaxValue = title.NeedValue
				archiveAnnex.MaxID = title.ID
			}
		}
		archiveTypeAnnex[archiveType] = archiveAnnex
	}
	t.TitleAnnexTitleMap[include.ArchiveType] = archiveTypeAnnex
}

var TitleConfigs *TitleConfigCol
