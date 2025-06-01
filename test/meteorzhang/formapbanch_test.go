package main

import (
	"math/rand"
	"testing"
)

type userInfo struct {
	userID   int32
	username string
}

// go test -bench='.' .\formapbanch_test.go  -benchmem

var UserList []userInfo
var UserMap map[int32]userInfo

var matNum int32 = 10

func init() {
	UserMap = make(map[int32]userInfo, matNum)
	UserList = make([]userInfo, matNum)
	for i := int32(0); i < matNum; i++ {
		userinfo := userInfo{userID: i, username: "username"}
		UserList[i] = userinfo
		UserMap[i] = userinfo
	}
}
func listIsExist(userID int32) bool {
	for _, v := range UserList {
		if v.userID == userID {
			return true
		}
	}
	return false
}

func listAdd(userID int32) {
	if listIsExist(userID) {
		user := userInfo{userID: userID, username: "username"}
		UserList = append(UserList, user)
	}
}

func getUserIDList() []int32 {
	ret := make([]int32, 0)
	for _, v := range UserList {
		ret = append(ret, v.userID)
	}
	return ret
}

func clearList() {
	UserList = make([]userInfo, 0)
}
func delUser(userID int32) {
	tmp := make([]userInfo, 0)
	for _, v := range UserList {
		if v.userID != userID {
			tmp = append(tmp, v)
		}
	}
	UserList = tmp
}

func forEnterMatch(num int32) {
	for i := int32(0); i < num; i++ {
		listAdd(i)
		if i%matNum == 9 {
			userIDL := getUserIDList()
			_ = userIDL
			clearList()
		}
	}
}

func twoForEnterMatch(num int32) {
	for i := int32(0); i < num; i++ {
		listAdd(i)
	}
	for i := int32(0); i < num; i++ {
		listAdd(i)
		if i%matNum == 9 {
			userIDL := getUserIDList()
			_ = userIDL
			clearList()
		}
	}
}

func mapGetUserID() []int32 {
	ret := make([]int32, 0)
	for _, v := range UserMap {
		ret = append(ret, v.userID)
	}
	return ret
}

func mapEnterMatch(num int32) {
	for i := int32(0); i < num; i++ {
		user := userInfo{userID: i, username: "username"}
		UserMap[i] = user
		if len(UserMap) == int(matNum) {
			userIDL := mapGetUserID()
			// _ = userIDL
			for k := range userIDL {
				delete(UserMap, userIDL[k])
			}
		}
	}
}

func BenchmarkForEnterMatch10(b *testing.B) {
	var num int32 = 10
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		forEnterMatch(num)
	}
}

func BenchmarkTwoForEnterMatch10(b *testing.B) {
	var num int32 = 10
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		twoForEnterMatch(num)
	}
}

func BenchmarkMapEnterMatch10(b *testing.B) {
	var num int32 = 10
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapEnterMatch(num)
	}
}

func BenchmarkForEnterMatch1000(b *testing.B) {
	var num int32 = 1000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		forEnterMatch(num)
	}
}

func BenchmarkTwoForEnterMatch1000(b *testing.B) {
	var num int32 = 1000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		twoForEnterMatch(num)
	}
}

func BenchmarkMapEnterMatch1000(b *testing.B) {
	var num int32 = 1000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		mapEnterMatch(num)
	}
}

func AddUserList(num int32) {
	UserList = make([]userInfo, num)
	for i := int32(0); i < num; i++ {
		// UserList = append(UserList, userInfo{userID: i, username: "username"})
		UserList[i] = userInfo{userID: i, username: "username"}
	}
}

func AddMap(num int32) {
	for i := int32(0); i < num; i++ {
		UserMap[i] = userInfo{userID: i, username: "username"}
	}
}

func SearchUserList(num int32) {
	for _, v := range UserList {
		if v.userID == num {
			_ = v.userID
		}
	}
}

func SearchMap(num int32) {
	if v, k := UserMap[num]; k {
		_ = v
	}
}

func BenchmarkAddUserList1000(b *testing.B) {
	var num int32 = 10000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AddUserList(num)
	}
}

func BenchmarkAddMap1000(b *testing.B) {
	var num int32 = 10000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AddMap(num)
	}
}

func BenchmarkSearchUserList1000(b *testing.B) {
	num := int32(rand.Int63n(int64(matNum)))
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SearchUserList(num)
	}
}

func BenchmarkSearchMap1000(b *testing.B) {
	num := int32(rand.Int63n(int64(matNum)))
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SearchMap(num)
	}
}

func ForList() {
	for i, v := range UserList {
		_ = v
		_ = i
	}
}

func ForMap() {
	for i, v := range UserMap {
		_ = v
		_ = i
	}
}

func BenchmarkForList1000(b *testing.B) {
	num := int32(rand.Int63n(int64(matNum)))
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SearchUserList(num)
	}
}

func BenchmarkForMap1000(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		ForMap()
	}
}
