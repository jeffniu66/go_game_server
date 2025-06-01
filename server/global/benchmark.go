package global

import (
	"go_game_server/server/logger"
	"runtime"
)

func RoomNumIsHealthy() bool {
	warnNum := MyConfig.ReadInt32("bench", "room_warn_num")
	num := GloInstance.roomNumAdd - GloInstance.roomNumDel
	if num > warnNum {
		logger.Log.Warnf("this room num is :%d out %d", num, warnNum)
		return false
	}
	return true
}

func MemIsHealthy() bool {
	var m runtime.MemStats

	runtime.ReadMemStats(&m)
	useMem := m.Sys / 1024 // KB
	memWarn := uint64(MyConfig.ReadInt32("bench", "mem_warn_num")) * 1024 * 1024
	if useMem > memWarn && memWarn > 0 {
		logger.Log.Warnf("this mem is :%d kb out %d kb", useMem, memWarn)
		return false
	}
	return true
}

func CpuIsHealthy() bool {
	return true
}
