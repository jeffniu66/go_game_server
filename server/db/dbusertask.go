package db

//
//import (
//	"encoding/json"
//	"fmt"
//	"go_game_server/server/include"
//	"go_game_server/server/util"
//)
//
//type Task = include.Task
//
//type TasksMap map[int32]*Task
//
//func (tMap TasksMap) InitData(userId int32) TasksMap {
//	nowStamp := util.UnixTime()
//	rows, err := DB.Query("select task_id, type, subtype, target, task_state, create_time, expire_time from t_user_task where user_id = ?", userId)
//	defer rows.Close()
//	util.CheckErr(err)
//	for rows.Next() {
//		t := &Task{}
//		var targetStr string
//		err = rows.Scan(&t.TaskId, &t.Type, &t.SubType, &targetStr, &t.State, &t.CreateTime, &t.ExpireTime)
//		util.CheckErr(err)
//
//		var targetList []*include.Target
//		_ = json.Unmarshal([]byte(targetStr), &targetList)
//		t.TargetList = targetList
//		// 将过期的任务删除
//		if t.ExpireTime > 0 && t.ExpireTime < nowStamp {
//			t.Update = include.DelUpdate
//		}
//		tMap[t.TaskId] = t
//	}
//	return tMap
//}
//
//func (tMap TasksMap) SaveData(userId int32) interface{} {
//	for id, task := range tMap {
//		if task.Update == include.Update {
//			InsertTask(task, userId)
//			task.Update = include.UnUpdate
//		}
//		if task.Update == include.DelUpdate {
//			DeleteTask(task, userId)
//			delete(tMap, id)
//		}
//	}
//	return nil
//}
//
//// 添加任务
//func InsertTask(tasks *Task, userId int32) {
//	TargetStr, err := json.Marshal(tasks.TargetList)
//	if err != nil {
//		fmt.Println(err)
//	}
//	ExecDB(UserDBType, userId, "replace into t_user_task(user_id, task_id, type, subtype, target, task_state, create_time, expire_time) "+
//		"values(?, ?, ?, ?, ?, ?, ?, ?)",
//		userId, tasks.TaskId, tasks.Type, tasks.SubType, TargetStr, tasks.State, tasks.CreateTime, tasks.ExpireTime)
//}
//
//// 删除任务
//func DeleteTask(task *Task, userId int32) {
//	ExecDB(UserDBType, userId, "delete from t_user_task where user_id = ? and task_id = ? ", userId, task.TaskId)
//}
//
//func GetUserDBTaskMap(userId int32) TasksMap {
//	tasks := make(TasksMap)
//	tasksMap := tasks.InitData(userId)
//	return tasksMap
//}
//
//func GetUserDBTaskMapByType(userId, subType int32) map[int32]*Task {
//	tasks := make(map[int32]*Task)
//	rows, err := DB.Query("select task_id, type, subtype, target, task_state, create_time, expire_time from t_user_task"+
//		" where user_id = ? and subtype = ? and task_state = ?", userId, subType, include.TaskStateNotFinish)
//	defer rows.Close()
//	util.CheckErr(err)
//	for rows.Next() {
//		t := &Task{}
//		var targetStr string
//		err = rows.Scan(&t.TaskId, &t.Type, &t.SubType, &targetStr, &t.State, &t.CreateTime, &t.ExpireTime)
//		util.CheckErr(err)
//
//		var targetList []*include.Target
//		_ = json.Unmarshal([]byte(targetStr), &targetList)
//		t.TargetList = targetList
//		tasks[t.TaskId] = t
//	}
//	return tasks
//}
