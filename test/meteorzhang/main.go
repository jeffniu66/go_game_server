package main

import (
	"fmt"
	_ "net/http/pprof"
	"runtime"
	"sync"
	"time"
)

var maxNum = 10000
var procs = 4

func main() {
	runtime.GOMAXPROCS(procs)
	match := initmatch()
	var wg sync.WaitGroup
	currentTime := time.Now().UnixNano()
	fmt.Println("match Now:", currentTime)
	for i := 0; i < maxNum; i++ {
		wg.Add(1)
		go func(userID int) {
			// fmt.Println("r:", userID)
			defer wg.Done()
			match.matchMgrAddPlayer(userID)
		}(i)
	}
	wg.Wait()
	currentTime1 := time.Now().UnixNano()
	matchTime := currentTime1 - currentTime
	fmt.Printf("match end and rank match Now:%d mast:%d us\n", currentTime1, matchTime)

	wgg := sync.WaitGroup{}
	for i := 0; i < maxNum; i++ {
		wgg.Add(1)
		go func(userId int) {
			defer wgg.Done()
			// fmt.Println(userId)
			match.runRankMatch(userId%7, userId)
		}(i)
	}
	wgg.Wait()
	currentTime2 := time.Now().UnixNano()
	rmatchTime := currentTime2 - currentTime1
	fmt.Printf("rank match end:%d mast:%d us\n", currentTime2, rmatchTime)

	d := (matchTime - rmatchTime) / 1000000
	fmt.Printf("match - rank match:%d us, %d ms", matchTime-rmatchTime, d)

	// http.ListenAndServe("0.0.0.0:6060", nil)
}

func OnlyMatch() {
	var wg sync.WaitGroup
	runtime.GOMAXPROCS(procs)
	match := initmatch()
	for i := 0; i < maxNum; i++ {
		wg.Add(1)
		go func(userID int) {
			// fmt.Println("r:", userID)
			defer wg.Done()
			match.matchMgrAddPlayer(userID)
		}(i)
	}
	wg.Wait()
}

// MoreMatch afe
func MoreMatch() {
	runtime.GOMAXPROCS(procs)
	match := initmatch()
	wgg := sync.WaitGroup{}
	for i := 0; i < maxNum; i++ {
		wgg.Add(1)
		go func(userId int) {
			defer wgg.Done()
			// fmt.Println(userId)
			match.runRankMatch(userId%7, userId)
		}(i)
	}
	wgg.Wait()
}
func initmatch() *matchMgr {
	tmp := make(map[int]*rankMatch)

	match := &matchMgr{
		matchNum:     10,
		rankMatchMap: tmp,
	}
	for i := 0; i < 7; i++ {
		ch := make(chan int, 12)
		rankMatch := &rankMatch{
			matchNum: 10,
			req:      ch,
		}
		go rankMatch.rankMatchMgrAddPlayer()
		match.rankMatchMap[i] = rankMatch
	}
	return match
}

type matchMgr struct {
	playerIdList []int
	matchNum     int // 策划配置
	rankMatchMap map[int]*rankMatch
}

func (m *matchMgr) matchMgrAddPlayer(userID int) {
	m.playerIdList = append(m.playerIdList, userID)
	if len(m.playerIdList) >= m.matchNum {
		_ = "enter room"
		// fmt.Println("enter room")
	}
}

func (m *matchMgr) runRankMatch(rank, userID int) {
	rankMatch := m.rankMatchMap[rank]
	rankMatch.ReadData(userID)
}

/////////////////////////////////////
type rankMatch struct {
	req          chan int
	playerIdList []int
	matchNum     int // 策划配置
}

func (r *rankMatch) rankMatchMgrAddPlayer() {
	for {
		select {
		case id := <-r.req:
			r.playerIdList = append(r.playerIdList, id)
			if len(r.playerIdList) >= r.matchNum {
				// fmt.Println("enter room")
				_ = "enter room"
			}
		}
	}
}

func (r *rankMatch) ReadData(userID int) {
	r.req <- userID
}
