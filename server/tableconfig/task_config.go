package tableconfig

import (
	"errors"
)

type TaskConfig struct {
	Id            int    `json:"id"`
	Missiontarget int    `json:"missiontarget"`
	Map           string `json:"map"`
	Name          string `json:"name"`
	MissionType   int    `json:"missiontype"`
	Fatal         int    `json:"fatal"`
	Keeptime      int    `json:"keeptime"`
}

type TaskConfigCol struct {
	TaskList       []TaskConfig                    `json:"task_list"`
	TaskMap        map[int]TaskConfig              `json:"task_map"`
	TaskMapMap     map[string][]TaskConfig         // 地图
	TaskMapTypeMap map[string]map[int][]TaskConfig // 地图->类型
}

func (t *TaskConfigCol) GetTaskById(id int) (*TaskConfig, error) {
	for _, v := range t.TaskList {
		if v.Id == id {
			return &v, nil
		}
	}
	return nil, errors.New("task is null")
}

func (t *TaskConfigCol) GetTaskByKey(key int) (*TaskConfig, error) {
	if _, ok := t.TaskMap[key]; !ok {
		return nil, errors.New("task is null")
	}
	ret := t.TaskMap[key]
	return &ret, nil
}

func (t *TaskConfigCol) GetTaskTypeMap(key string) []TaskConfig {
	return t.TaskMapMap[key]
}

func (t *TaskConfigCol) InitMap() {
	t.TaskMap = make(map[int]TaskConfig)
	t.TaskMapMap = make(map[string][]TaskConfig)
	t.TaskMapTypeMap = make(map[string]map[int][]TaskConfig)
	for _, v := range t.TaskList {
		t.TaskMap[v.Id] = v

		typeMap := t.TaskMapMap[v.Map]
		if typeMap == nil {
			typeMap = []TaskConfig{}
		}
		typeMap = append(typeMap, v)
		t.TaskMapMap[v.Map] = typeMap

		mapTypeMap := t.TaskMapTypeMap[v.Map]
		if mapTypeMap == nil {
			mapTypeMap = make(map[int][]TaskConfig)
			t.TaskMapTypeMap[v.Map] = mapTypeMap
		}
		typeMap2 := mapTypeMap[v.Missiontarget]
		if typeMap2 == nil {
			typeMap2 = []TaskConfig{}
		}
		typeMap2 = append(typeMap2, v)
		mapTypeMap[v.Missiontarget] = typeMap2
		t.TaskMapTypeMap[v.Map] = mapTypeMap
	}
}

var TaskConfigs *TaskConfigCol
