package global

import (
	"go_game_server/server/logger"
	"sync"
)

type roomA struct {
	isUsed int32
	room   interface{}
}

func (r *roomA) GetRoom() interface{} {
	return r.room
}

type roomList struct {
	roomMap *sync.Map // key 表示房间号 value-room
}

func initRoomList() *roomList {
	tmp := &sync.Map{}
	roomNum := MyConfig.ReadInt32("room", "room_max_num")
	// roomNum := int32(10)
	for i := int32(0); i < roomNum; i++ {
		room := &roomA{int32(0), nil}
		tmp.Store(i, room)
	}

	return &roomList{
		roomMap: tmp,
	}
}

func (r *roomList) GetUseableRoomId() int32 {
	var roomId int32 = -1
	r.roomMap.Range(func(k, v interface{}) bool {
		kk := k.(int32)
		vv := v.(*roomA)
		logger.Log.Infof("roomId:%v room:%v", kk, vv)
		if vv.isUsed == 0 {
			roomId = kk
			return false
		}
		return true
	})

	return roomId
}

// roomer Room 指针
func (r *roomList) ChangeRoomIdUsed(roomId int32, roomer interface{}) {
	ro := &roomA{isUsed: int32(1), room: roomer}
	r.roomMap.Store(roomId, ro)
}

func (r *roomList) ChangeRoomIdUnused(roomId int32) {
	ro := &roomA{isUsed: int32(0), room: nil}
	r.roomMap.Store(roomId, ro)
}

func (r *roomList) GetRoomID(roomID int32) interface{} {
	if room, ok := r.roomMap.Load(roomID); ok {
		ro := room.(*roomA)
		return ro.GetRoom()
	}
	return nil
}

func (r *roomList) LoadRoomIAdValue(roomId int32) (interface{}, bool) {
	ro, b := r.roomMap.Load(roomId)
	return ro.(*roomA).room, b
}

func (r *roomList) GetUsedRoomFaceList() []interface{} {
	ret := []interface{}{}

	(*r.roomMap).Range(func(k, v interface{}) bool {
		vv := v.(*roomA)
		logger.Log.Debugf("roomId:%v room:%v", k.(int32), vv)
		if vv.isUsed == 1 && vv.room != nil {
			ret = append(ret, vv.room)
		}
		return true
	})
	return ret
}

func (r *roomList) GetRoomIDList() []int32 {
	var listID []int32
	r.roomMap.Range(func(roomID, _ interface{}) bool {
		listID = append(listID, roomID.(int32))
		return true
	})
	return listID
}
