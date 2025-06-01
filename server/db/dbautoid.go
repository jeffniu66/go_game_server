package db

import (
	"sync"
)

type increaseID struct {
	userID int32
	mutex  sync.Mutex
}

type IncreaseItemUID struct {
	itemUid int32
	mutex   sync.Mutex
}

var IDInstance *increaseID           // 全局共享实例,id分配
var ItemUidInstance *IncreaseItemUID // 道具自增uid

func NewIncrease() {
	userID := GetUserMaxID()
	itemUid := GetItemMaxID()
	//logger.Log.Info(">>>>>>>>>>>> init increaseid success = ", userID)
	IDInstance = &increaseID{userID: userID}
	ItemUidInstance = &IncreaseItemUID{itemUid: itemUid}
}

func (i *increaseID) GetNewUserID() (userID int32) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	userID = i.userID
	i.userID += 1
	return
}

func (i *IncreaseItemUID) GetNewItemUid() (itemUid int32) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	itemUid = i.itemUid
	i.itemUid += 1
	return
}
