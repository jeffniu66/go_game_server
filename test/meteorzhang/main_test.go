package main

import (
	"fmt"
	"go_game_server/server/util"
	"sync"
	"testing"
	"time"
)

func BenchmarkOnlyMatch1Q(b *testing.B) {
	maxNum = 1000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		OnlyMatch()
	}
}

func BenchmarkMoreMatch1Q(b *testing.B) {
	maxNum = 1000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MoreMatch()
	}
}

func BenchmarkOnlyMatch1W(b *testing.B) {
	maxNum = 10000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		OnlyMatch()
	}
}

func BenchmarkMoreMatch1W(b *testing.B) {
	maxNum = 10000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MoreMatch()
	}
}

func BenchmarkOnlyMatch1SW(b *testing.B) {
	maxNum = 100000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		OnlyMatch()
	}
}

func BenchmarkMoreMatch1SW(b *testing.B) {
	maxNum = 100000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MoreMatch()
	}
}

func BenchmarkOnlyMatch5SW(b *testing.B) {
	maxNum = 500000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		OnlyMatch()
	}
}

func BenchmarkMoreMatch5SW(b *testing.B) {
	maxNum = 500000
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		MoreMatch()
	}
}

type roomList struct {
	roomMap *sync.Map // key 表示房间号 value-room
}

type room struct {
	IsUsed int32
	Room   interface{}
}

func TestSyncMapRangeBreak(t *testing.T) {
	mp := sync.Map{}
	ro := roomList{
		roomMap: &mp,
	}

	for i := int32(0); i < 7; i++ {
		room := &room{int32(0), nil}
		if i == 1 || i == 4 {
			room.IsUsed = 1
		}
		ro.roomMap.Store(i, room)
	}
	ro.roomMap.Range(
		func(k, v interface{}) bool {
			// false 退出循环
			kk := k.(int32)
			vv := v.(*room)
			fmt.Printf("roomId:%v room:%v\n", kk, vv)
			if vv.IsUsed == 1 {
				// t.Logf("roInId:%v room:%v", kk, vv)
				return false
			}
			return true
		})
}

func TestChanWait(t *testing.T) {
	ch := make(chan int, 1)
	for i := 0; i < 3; i++ {
		go func(a int) {
			t.Log("write ", a)
			time.Sleep(1 * time.Second)
			ch <- a
			t.Log("write ", a, " end")
		}(i)
	}
	t.Log("reading wait time 2 second")
	time.Sleep(2 * time.Second)
	for i := 0; i < 3; i++ {
		select {
		case a := <-ch:
			t.Logf("read:%d", a)
		}
	}
	select {}
}

func TestP(t *testing.T) {
	ro := room{
		IsUsed: 12,
	}
	l := &ro
	// 指针类型修改原数据
	l.IsUsed = 2
	t.Log(ro)
}

func TestRandMp(t *testing.T) {
	mp := make(map[int]int, 0)
	mp[0] = 0
	mp[1] = 1
	mp[2] = 2
	mp[3] = 3
	mp[4] = 4
	l := int32(len(mp))
	i := util.RandInt(0, l-1)
	fmt.Println(i, l)
	n := int32(0)
	for k := range mp {
		fmt.Println("k:", k, n)
		if n == i {
			break
		}
		n++
	}
}
