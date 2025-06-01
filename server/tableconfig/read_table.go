package tableconfig

import (
	"bufio"
	"encoding/json"
	"go_game_server/server/logger"
	"io"
	"os"
)

type excelFace interface {
	InitMap()
}

type excelFiles struct {
	Path     string      `json:"path"`
	Data     interface{} `json:"data"`
	InitData interface{}
}

func initExcel() []excelFiles {
	ret := []excelFiles{}

	taskConfigCol := TaskConfigCol{}
	r := excelFiles{
		Path:     "./server/tableconfig/json/missionsConfig.json",
		Data:     &taskConfigCol.TaskList,
		InitData: &taskConfigCol,
	}
	TaskConfigs = &taskConfigCol
	ret = append(ret, r)

	constsConfigCol := ConstsConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/constsConfig.json",
		Data:     &constsConfigCol.ConstList,
		InitData: &constsConfigCol,
	}
	ConstsConfigs = &constsConfigCol
	ret = append(ret, r)

	dropGroupConfigCol := DropGroupConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/dropgroupsConfig.json",
		Data:     &dropGroupConfigCol.DropGroupList,
		InitData: &dropGroupConfigCol,
	}
	DropGroupConfigs = &dropGroupConfigCol
	ret = append(ret, r)

	itemGroupConfigCol := ItemConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/itemsConfig.json",
		Data:     &itemGroupConfigCol.ItemConfigList,
		InitData: &itemGroupConfigCol,
	}
	ItemConfigs = &itemGroupConfigCol
	ret = append(ret, r)

	quaLevelConfs := QualityingConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/qualifyinglevelsConfig.json",
		Data:     &quaLevelConfs.QuaLevelCofList,
		InitData: &quaLevelConfs,
	}
	QuaLevelConfs = &quaLevelConfs
	ret = append(ret, r)

	quaConfs := QuaConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/qualifyingsConfig.json",
		Data:     &quaConfs.QuaConfigList,
		InitData: &quaConfs,
	}
	QuaConfigs = &quaConfs
	ret = append(ret, r)

	levelConfigs := LevelConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/levelsConfig.json",
		Data:     &levelConfigs.LevelConfigList,
		InitData: &levelConfigs,
	}
	LevelConfigs = &levelConfigs
	ret = append(ret, r)

	titleConfigs := TitleConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/titlesConfig.json",
		Data:     &titleConfigs.TitleConfigList,
		InitData: &titleConfigs,
	}
	TitleConfigs = &titleConfigs
	ret = append(ret, r)

	ninjaConfigs := NinjaConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/ninjasConfig.json",
		Data:     &ninjaConfigs.NinjaConfigList,
		InitData: &ninjaConfigs,
	}
	NinjaConfigs = &ninjaConfigs
	ret = append(ret, r)

	avatarConfigs := AvatarConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/avatarsConfig.json",
		Data:     &avatarConfigs.AvatarConfigList,
		InitData: &avatarConfigs,
	}
	AvatarConfigs = &avatarConfigs
	ret = append(ret, r)

	mailConfigs := MailConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/mailsConfig.json",
		Data:     &mailConfigs.MailConfigList,
		InitData: &mailConfigs,
	}
	MailConfigs = &mailConfigs
	ret = append(ret, r)

	storeConfigs := StoreConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/storesConfig.json",
		Data:     &storeConfigs.StoreList,
		InitData: &storeConfigs,
	}
	StoresConfigs = &storeConfigs
	ret = append(ret, r)

	skinConfigs := SkinConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/skinsConfig.json",
		Data:     &skinConfigs.SkinList,
		InitData: &skinConfigs,
	}
	SkinConfigs = &skinConfigs
	ret = append(ret, r)

	nameZhConfigs := NameZhConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/name_zhsConfig.json",
		Data:     &nameZhConfigs.NameZhConfigList,
		InitData: &nameZhConfigs,
	}
	NameZhConfigs = &nameZhConfigs
	ret = append(ret, r)

	skillConfigs := SkillConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/skillsConfig.json",
		Data:     &skillConfigs.SkillList,
		InitData: &skillConfigs,
	}
	SkillConfigs = &skillConfigs
	ret = append(ret, r)

	senWordConfigs := SenWordConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/sensitivewordssConfig.json",
		Data:     &senWordConfigs.SenWordConfigList,
		InitData: &senWordConfigs,
	}
	SenWordConfigs = &senWordConfigs
	ret = append(ret, r)

	quickTextCols := QuickTextCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/quicktextsConfig.json",
		Data:     &quickTextCols.QuickTextConfigList,
		InitData: &quickTextCols,
	}
	QuickTextCols = &quickTextCols
	ret = append(ret, r)

	actionConfigCols := ActionConfigCol{}
	r = excelFiles{
		Path:     "./server/tableconfig/json/newuserBDsConfig.json",
		Data:     &actionConfigCols.ActionList,
		InitData: &actionConfigCols,
	}
	ActionConfigs = &actionConfigCols
	ret = append(ret, r)

	return ret
}

func ReadTable() error {
	excelList := initExcel()
	for _, v := range excelList {
		err := MissionTableToStruct(v.Path, v.Data)
		if err != nil {
			logger.Log.Errorf("read table %s failed, err:%v\n", v.Path, err)
			return err
		}
		if res, ok := v.InitData.(excelFace); ok {
			res.InitMap()
		}
	}
	return nil
}

func MissionTableToStruct(path string, dataList interface{}) error {
	//打开文件
	inputFile, err := os.Open(path)
	if err != nil {
		// logger.Log.Errorf("MissionTableToJson err = %v", err)
		return err
	}
	defer inputFile.Close()

	//按行读取文件
	var s string
	inputReader := bufio.NewReader(inputFile)
	for {
		inputString, _, readerError := inputReader.ReadLine()
		if readerError == io.EOF {
			break
		}
		s = s + string(inputString)
	}

	err = json.Unmarshal([]byte(s), dataList)
	if err != nil {
		logger.Log.Errorf("MissionTableToJson err = %v", err)
		return err
	}
	return nil
}
