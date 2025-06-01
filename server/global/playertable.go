package global

import (
	"fmt"
	"sync"
)

type globalList struct {
	playerMap  *sync.Map
	playerLen  int32
	roomList   *roomList
	roomNumAdd int32
	roomNumDel int32
}

var GloInstance *globalList

func NewGlobalPlayers() {
	roomList := initRoomList()

	GloInstance = &globalList{playerMap: &sync.Map{}, roomList: roomList}
}

func (g *globalList) GetPlayerList() *sync.Map {
	return g.playerMap
}

func (g *globalList) GetPlayersNum() (int, int32) {
	length := 0
	g.playerMap.Range(func(_, _ interface{}) bool {
		length++
		return true
	})
	return length, g.playerLen
}

func (g *globalList) GetPlayerIDList() []int32 {
	var listID []int32
	g.playerMap.Range(func(userID, _ interface{}) bool {
		listID = append(listID, userID.(int32))
		return true
	})
	return listID
}

func (g *globalList) AddPlayer(userID int32, player interface{}) {
	g.playerMap.Store(userID, player)
	//g.playerLen += 1	// 会有并发问题
}

func (g *globalList) DelPlayer(userID int32) {
	g.playerMap.Delete(userID)
	//g.playerLen -= 1 // 会有并发问题
}

func (g *globalList) GetPlayer(userID int32) interface{} {
	if player, ok := g.playerMap.Load(userID); ok {
		return player
	}
	return nil
}

func (g *globalList) BroadCast() {
	g.playerMap.Range(func(userID, _ interface{}) bool {
		fmt.Printf("---------------- broadcast userid = %d\n", userID)
		return true
	})
}

func (g *globalList) GetUseableRoomId() int32 {
	return g.roomList.GetUseableRoomId()
}

func (g *globalList) ChangeRoomIdUsed(roomId int32, room interface{}) {
	g.roomNumAdd++
	g.roomList.ChangeRoomIdUsed(roomId, room)
}

func (g *globalList) ChangeRoomIdUnused(roomId int32) {
	g.roomNumDel++
	g.roomList.ChangeRoomIdUnused(roomId)
}

func (g *globalList) GetUsedRoomFaceList() []interface{} {
	return g.roomList.GetUsedRoomFaceList()
}

func (g *globalList) GetRoom(roomID int32) interface{} {
	return g.roomList.GetRoomID(roomID)
}

func (g *globalList) GetRoomIDList() []int32 {
	return g.roomList.GetRoomIDList()
}
