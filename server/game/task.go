package game

import (
	"go_game_server/proto3"
	"go_game_server/server/constant"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"go_game_server/server/tableconfig"
	"go_game_server/server/util"
	"strconv"
	"strings"
	"time"
)

type Task include.Task
type UrgencyTask include.UrgencyTask

func (t *Task) GetFinishTaskNum() int32 {
	l := int32(len(strings.Split(t.PointInfo, ",")))
	return l
}

func (t *Task) GetTaskNum() int32 {
	if t.AssignTaskPoint == "" {
		return 0
	}
	n := strings.Split(t.AssignTaskPoint, ",")
	return int32(len(n))
}

// 初始化技能池
func InitSkillPool(skillPoolType int32) map[int32]int32 {
	poolIdStr := tableconfig.ConstsConfigs.GetValueById(skillPoolType)
	poolId, err := strconv.Atoi(poolIdStr)
	if err != nil {
		logger.Log.Errorln("InitSkill poolId is error!")
		return nil
	}
	dropGroupConfigs := tableconfig.DropGroupConfigs.GetDropGroupMap(poolId)
	if len(dropGroupConfigs) == 0 {
		logger.Log.Errorln("InitSkill dropGroupConfigs is error!")
		return nil
	}
	skillPoolMap := make(map[int32]int32)
	for _, v := range dropGroupConfigs {
		skillPoolMap[int32(v.ItemId)] = int32(v.Num)
	}
	return skillPoolMap
}

type ImportSkill struct {
	dropGroupId int32 // 掉落表id
	useNum      int32 // 使用数量
}

func GetFinishPoints(pointInfo string) (finishPoints []int32) {
	if pointInfo == "" {
		return
	}
	points := strings.Split(pointInfo, ",")
	for _, v := range points {
		point, _ := strconv.Atoi(v)
		finishPoints = append(finishPoints, int32(point))
	}
	return
}

func GetAssignTaskPoints(assignTaskInfo string) (assignTaskPoints []int32) {
	if assignTaskInfo == "" {
		return
	}
	points := strings.Split(assignTaskInfo, ",")
	for _, v := range points {
		point, _ := strconv.Atoi(v)
		if point == 0 {
			continue
		}
		assignTaskPoints = append(assignTaskPoints, int32(point))
	}
	return
}

func (r *Room) RankTaskIds(userId int32) (positions []int32) {
	if r.taskMap == nil {
		r.taskMap = make(map[int32]*Task)
	}
	if r.roomInfo.IsWolfMan(int32(userId)) { // 狼人不能被分配任务
		return
	}
	//taskConfigList := tableconfig.TaskConfigs.GetTaskTypeMap("map_0001")
	taskMapTypeMap := tableconfig.TaskConfigs.TaskMapTypeMap["map_0001"]
	taskConfigList := taskMapTypeMap[constant.TaskTypeNormal]

	//rand.Seed(time.Now().UnixNano()) // 以当前系统时间作为种子参数

	minId := taskConfigList[0].Id
	var maxId int
	for _, v := range taskConfigList {
		if v.Missiontarget != constant.TaskTypeNormal { // 普通任务
			continue
		}
		if v.Id < minId && v.Id != 0 {
			minId = v.Id
		}
		if v.Id > maxId {
			maxId = v.Id
		}
	}
	var posStr string
	for i := 0; i < constant.TaskPoints; i++ {
		randPosition := randPositionRecursion(int32(minId), int32(maxId), posStr)

		posStr += strconv.Itoa(int(randPosition)) + ","

		positions = append(positions, randPosition)
	}
	if len(posStr) > 0 {
		posStr = posStr[:len(posStr)-1]
	}
	r.taskMap[userId] = &Task{UserId: userId, PointInfo: "", AssignTaskPoint: posStr}
	return positions
}

func randPositionRecursion(min int32, max int32, posStr string) int32 {
	randPosition := util.RandInt(min, max)
	if !validPositionSame(posStr, randPosition) {
		return randPosition
	}
	return randPositionRecursion(min, max, posStr)
}

func validPositionSame(posStr string, position int32) bool {
	if posStr == "" {
		return false
	}
	poss := strings.Split(posStr, ",")
	for _, v := range poss {
		posInt, _ := strconv.Atoi(v)
		if int32(posInt) == position {
			return true
		}
	}
	return false
}

// 完成任务
func (r *Room) FinishTask(userId int32, pointId int32) {
	pointIds := r.FinishTaskCore(userId, int(pointId))
	if len(pointIds) <= 0 {
		logger.Log.Errorln("任务点完成小于0")
		return
	}
	player := global.GloInstance.GetPlayer(userId)
	if player != nil {
		p := player.(*Player)
		pbData := &proto3.FinishMissionResp{PointIds: pointIds}
		p.Pid.Cast("finishTaskResp", pbData)

		dropItemMap := r.dropItem(constant.NormalManSkillPool, userId)
		// 即使没掉技能也给客户端下发消息（策划需求）
		p.Pid.Cast("dropItemResp", dropItemMap)
	} else {
		logger.Log.Errorln("完成任务玩家为nil")
	}

	r.totalScore++
	// 判断是否达到最大任务分数
	maxTotalScoreStr := tableconfig.ConstsConfigs.GetValueById(constant.TaskTotalScoreId)
	if maxTotalScoreStr == "" {
		p := player.(*Player)
		p.ErrorResponse(proto3.ErrEnum_Error_ConfigId_NotExists, proto3.ErrEnum_name[int32(proto3.ErrEnum_Error_ConfigId_NotExists)])
		return
	}
	maxTotalScore, _ := strconv.Atoi(maxTotalScoreStr)
	if r.totalScore >= maxTotalScore {
		r.roomInfo.GameWinRole = proto3.PlayerRoleEnum_comm_people
		// 游戏结束逻辑调用
		r.roomPid.Cast("gameEnd", nil)
		logger.Log.Infof("room:%d gameEnd, because totalScore >= maxTotalScore:%v", r.roomId, maxTotalScore)
		return
	}
	// 广播分数
	cmd := proto3.ProtoCmd_CMD_TotalScoreResp
	pbData := &proto3.TotalScoreResp{}
	pbData.TotalScore = int32(r.totalScore)
	r.spreadPlayers(cmd, pbData)
}

func (r *Room) FinishTaskCore(userId int32, positionId int) (pointIds []int32) {
	taskConfig, err := tableconfig.TaskConfigs.GetTaskById(positionId)
	if err != nil {
		logger.Log.Errorln("taskConfig is not exists!", taskConfig)
		return pointIds
	}
	task := r.taskMap[userId]
	pointInfo := task.PointInfo
	if isFinishTask(int32(positionId), pointInfo) {
		logger.Log.Errorln("FinishTaskCore isFinishTask pointInfo: " + pointInfo)
		return
	}
	if isFinishAllTask(pointInfo, task.AssignTaskPoint) {
		logger.Log.Errorf("完成了所有任务 pointInfo: %v, AssignTaskPoint: %v", pointInfo, task.AssignTaskPoint)
		return
	}
	if pointInfo == "" {
		pointInfo = strconv.Itoa(positionId)
	} else {
		points := strings.Split(pointInfo, ",")
		if existsPoint(points, strconv.Itoa(positionId)) {
			logger.Log.Errorln("已完成这个任务: " + pointInfo)
			return pointIds
		}
		pointInfo = pointInfo + "," + strconv.Itoa(positionId)
	}
	task.PointInfo = pointInfo
	if r.taskMap == nil {
		logger.Log.Info("FinishTaskCore taskMap 初始化")
		r.taskMap = make(map[int32]*Task)
	}
	r.taskMap[userId] = task

	points := strings.Split(pointInfo, ",")
	logger.Log.Println("points: ", points)
	for _, v := range points {
		s, err := strconv.Atoi(v)
		if err != nil {
			logger.Log.Errorln("FinishTask is error: ", err)
			continue
		}
		pointIds = append(pointIds, int32(s))
	}
	return pointIds
}

// 是否完成该任务
func isFinishTask(pointId int32, pointInfo string) bool {
	if pointInfo == "" {
		return false
	}
	pointIds := strings.Split(pointInfo, ",")
	for _, v := range pointIds {
		if util.ToInt(v) == pointId {
			return true
		}
	}
	return false
}

func isFinishAllTask(pointInfo string, assignTaskPoint string) bool {
	if pointInfo == "" {
		return false
	}
	pointInfos := strings.Split(pointInfo, ",")
	assignTaskPoints := strings.Split(assignTaskPoint, ",")
	if len(pointInfos) >= len(assignTaskPoints) {
		return true
	}
	return false
}

func existsPoint(points []string, positionId string) (flag bool) {
	for _, v := range points {
		if v == positionId {
			return true
		}
	}
	return false
}

func GetPlayerTaskProgressNum(task *Task) int32 {
	pointInfo := task.PointInfo
	if len(pointInfo) <= 0 {
		return 0
	}
	points := strings.Split(pointInfo, ",")
	return int32(len(points))
}

// ===================================紧急任务===========================================

// 初始化紧急任务信息
func (r *Room) InitUrgencyTaskInfo(userId, keepTime, triggerPoint int32) (errEnum proto3.ErrEnum) {
	taskConfig, err := tableconfig.TaskConfigs.GetTaskById(int(triggerPoint))
	if err != nil {
		logger.Log.Errorf("this task:%d don't set in mission.json")
		return proto3.ErrEnum_Error_Operation_Fail
	}

	var first bool
	urgencyTask := r.urgencyTask
	if urgencyTask == nil {
		urgencyTask = &UrgencyTask{LastTriggerTime: time.Now().Unix(), IngPoint: triggerPoint}
		r.urgencyTask = urgencyTask
		first = true
		// return proto3.ErrEnum_Error_Pass
	}
	// 冷却时间校验
	urgencyTaskCd := tableconfig.ConstsConfigs.GetValueById(constant.UrgencyTaskCdId)
	lastCallTime := urgencyTask.LastTriggerTime
	cd, err := strconv.Atoi(urgencyTaskCd)
	if err != nil {
		logger.Log.Errorf("AssignUrgencyTaskPoint is error: %v", err)
		return proto3.ErrEnum_Error_Operation_Fail
	}
	if time.Now().Unix()-lastCallTime <= int64(cd) && !first {
		return proto3.ErrEnum_Error_Operation_Fail
	}
	urgencyTask.LastTriggerTime = time.Now().Unix()
	urgencyTask.IngPoint = triggerPoint
	urgencyTask.IngUserId = userId
	r.urgencyTask = urgencyTask
	if taskConfig.Fatal == int(proto3.CommonStatusEnum_true) { // 致命任务
		logger.Log.Infof("room_id:%d task:%d start 致命任务。", r.roomId, triggerPoint)
		r.roomPid.SendAfter("validGameEnd", "validGameEnd", keepTime*1000, nil)
	}
	return proto3.ErrEnum_Error_Pass
}

// 狼人触发紧急任务
func (r *Room) TriggerUrgencyTask(userId int32, triggerPointId int32) {
	if !r.roomInfo.IsWolfMan(userId) { // 非狼人不能召唤紧急任务
		return
	}
	// 关门期间不能触发紧急任务
	urgencyTask := r.urgencyTask
	closeDoorTime := r.closeDoorTime
	if urgencyTask != nil && closeDoorTime != 0 {
		curTime := time.Now().Unix()
		closeDoorDurationTime := util.ToInt(tableconfig.ConstsConfigs.GetValueById(constant.CloseDoorDurationTime))
		if curTime < closeDoorTime+int64(closeDoorDurationTime) {
			logger.Log.Errorln("TriggerUrgencyTask is ing")
			return
		}
	}
	taskConfig, _ := tableconfig.TaskConfigs.GetTaskByKey(int(triggerPointId))
	keepTime := taskConfig.Keeptime
	taskEndTime := time.Now().Unix() + int64(keepTime)
	cd := tableconfig.ConstsConfigs.GetValueById(constant.UrgencyTaskCdId)
	cdInt, _ := strconv.Atoi(cd)
	cdEndTime := time.Now().Unix() + int64(cdInt)

	errNum := r.InitUrgencyTaskInfo(userId, int32(keepTime), triggerPointId)
	if errNum != proto3.ErrEnum_Error_Pass {
		player := global.GloInstance.GetPlayer(userId)
		if player != nil {
			p := player.(*Player)
			p.ErrorResponse(errNum, proto3.ErrEnum_name[int32(errNum)])
		}
		return
	}

	triggerNumMap := r.urgencyTask.TriggerNumMap
	if triggerNumMap == nil {
		triggerNumMap = make(map[int32][]int32)
		triggerNumMap[userId] = []int32{triggerPointId}
	} else {
		triggerNumMap[userId] = append(triggerNumMap[userId], triggerPointId)
	}
	r.urgencyTask.TriggerNumMap = triggerNumMap

	pbData := &proto3.UrgencyTaskResp{TriggerPoint: triggerPointId, TaskEndTime: taskEndTime, CdEndTime: cdEndTime}
	r.spreadPlayers(proto3.ProtoCmd_CMD_UrgencyTaskResp, pbData)
}

// 平民完成紧急任务
func (r *Room) FinishUrgencyTask(userId int32, pointId int32) {
	pointIds, errNum := r.finishUrgencyTaskCore(userId, pointId)
	if errNum != proto3.ErrEnum_Error_Pass {
		player := global.GloInstance.GetPlayer(userId)
		if player != nil {
			p := player.(*Player)
			p.ErrorResponse(errNum, proto3.ErrEnum_name[int32(errNum)])
			return
		}
	}
	pbData := &proto3.FinishUrgencyTaskResp{PointIds: pointIds}
	r.spreadPlayers(proto3.ProtoCmd_CMD_FinishUrgencyTaskResp, pbData)
}

func (r *Room) GetTriggerUrgencyNum(userId int32) int32 {
	if !r.roomInfo.IsWolfMan(userId) { // 非狼人不能召唤紧急任务
		return 0
	}
	urgencyTask := r.urgencyTask
	if urgencyTask == nil {
		return 0
	}
	triggerNum := urgencyTask.TriggerNumMap[userId]
	return int32(len(triggerNum))
}

func (r *Room) CancelUrgencyTask() {
	if r.urgencyTask == nil {
		return
	}

	userTask := r.urgencyTask.TriggerNumMap[r.urgencyTask.IngUserId]
	for i := 0; i < len(userTask); i++ {
		if userTask[i] == r.urgencyTask.IngPoint {
			userTask = append(userTask[:i], userTask[i+1:]...)
			break
		}
	}
	r.urgencyTask.IngPoint = 0 // 取消紧急任务
	r.urgencyTask.TriggerNumMap[r.urgencyTask.IngUserId] = userTask
}

func (r *Room) validGameEnd() {
	if r.urgencyTask == nil {
		logger.Log.Errorln("validGameEnd r.urgencyTask is nil!!!")
		return
	}
	if r.roomInfo.RoomStatus == proto3.RoomStatus_voting {
		r.CancelUrgencyTask()
		return
	}
	taskConfig, err := tableconfig.TaskConfigs.GetTaskByKey(int(r.urgencyTask.IngPoint))
	if err != nil {
		logger.Log.Errorf("Urgency task :%d validGameEnd taskConfig is null: %v", r.urgencyTask.IngPoint, err)
		return
	}
	if taskConfig.Fatal != constant.UrgencyTaskTypeDead { // 非致命任务不处理
		return
	}
	taskPoints := r.urgencyTask.TaskPoints
	if len(taskPoints) != 0 {
		// 完成了紧急任务
		if isFinishTaskPoint(taskPoints, r.urgencyTask.IngPoint) {
			return
		}
	}
	// 游戏结束
	r.roomInfo.GameWinRole = proto3.PlayerRoleEnum_wolf_man
	r.roomPid.Cast("gameEnd", nil)
}

// 是否完成任务点
func isFinishTaskPoint(taskPoints []int32, taskPoint int32) bool {
	if len(taskPoints) == 0 {
		return false
	}
	for _, v := range taskPoints {
		if v == taskPoint {
			return true
		}
	}
	return false
}

func (r *Room) finishUrgencyTaskCore(userId int32, pointId int32) (pointIds []int32, errEnum proto3.ErrEnum) {
	if pointId <= 0 {
		return pointIds, proto3.ErrEnum_Error_Operation_Fail
	}
	urgencyTask := r.urgencyTask
	if urgencyTask == nil {
		return pointIds, proto3.ErrEnum_Error_Operation_Fail
	}
	// 不是正在进行的任务点
	if pointId != urgencyTask.IngPoint {
		return pointIds, proto3.ErrEnum_Error_Operation_Fail
	}
	taskPoints := urgencyTask.TaskPoints
	// 任务点是否在配置表里
	if !existsUrgencyPoint(GetAssignUrgencyTaskPoints(), pointId) {
		return taskPoints, proto3.ErrEnum_Error_Operation_Fail
	}
	// 狼人不能做紧急任务
	if r.roomInfo.IsWolfMan(userId) {
		return taskPoints, proto3.ErrEnum_Error_Operation_Fail
	}
	if len(taskPoints) == 0 {
		taskPoints = append(taskPoints, pointId)
	} else { // 只会包含1个元素
		taskPoints = []int32{pointId}
	}
	urgencyTask.TaskPoints = taskPoints
	urgencyTask.IngPoint = 0 // 完成后恢复
	r.urgencyTask = urgencyTask

	return taskPoints, proto3.ErrEnum_Error_Pass
}

func existsUrgencyPoint(points []int32, pointId int32) (flag bool) {
	for _, v := range points {
		if v == pointId {
			return true
		}
	}
	return false
}

// 获取紧急任务信息-断线重连
func (r *Room) GetUrgencyTaskInfoResp() *proto3.UrgencyTaskResp {
	urgencyTask := r.urgencyTask
	if urgencyTask == nil {
		return &proto3.UrgencyTaskResp{}
	}
	taskConfig, _ := tableconfig.TaskConfigs.GetTaskByKey(int(urgencyTask.IngPoint))
	if taskConfig == nil { // 没有正在进行的紧急任务点
		return &proto3.UrgencyTaskResp{}
	}
	keepTime := taskConfig.Keeptime
	taskEndTime := urgencyTask.LastTriggerTime + int64(keepTime)

	cd := tableconfig.ConstsConfigs.GetValueById(constant.UrgencyTaskCdId)
	cdInt, _ := strconv.Atoi(cd)
	cdEndTime := urgencyTask.LastTriggerTime + int64(cdInt)
	resp := &proto3.UrgencyTaskResp{TriggerPoint: urgencyTask.IngPoint, TaskEndTime: taskEndTime, CdEndTime: cdEndTime}
	return resp
}

// 从配置表获取紧急任务点-匹配成功后
func GetAssignUrgencyTaskPoints() []int32 {
	taskMapTypeMap := tableconfig.TaskConfigs.TaskMapTypeMap["map_0001"]
	taskConfigList := taskMapTypeMap[constant.TaskTypeUrgency]

	var taskConfigIds []int32
	for _, v := range taskConfigList {
		taskConfigIds = append(taskConfigIds, int32(v.Id))
	}
	return taskConfigIds
}

func GetFinishTaskNum(taskInfo string) int32 {
	if len(taskInfo) == 0 {
		return 0
	}
	tasks := strings.Split(taskInfo, ",")
	return int32(len(tasks))
}
