package global

import "testing"

func TestRoomList(t *testing.T) {
	tm := initRoomList()
	id := tm.GetUseableRoomId()
	t.Log(id)
	// t.Log(tm.LoadRoomIdValue(id))
	// // tm.ChangeRoomIdUsed(id)
	// t.Log(tm.LoadRoomIdValue(id))
	// tm.ChangeRoomIdUnused(id)
	// t.Log(tm.LoadRoomIdValue(id))
}
