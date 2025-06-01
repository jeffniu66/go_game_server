package main

import (
	"fmt"
	"testing"
	"time"
)

func andSlice(data Data, index int32) {
	sliceData = append(sliceData, data)

	for i := 0; i < len(sliceData); i++ {
		if sliceData[i].A == index {
			sliceData = append(sliceData[:i], sliceData[i+1:]...)
			break
		}
	}
}

func andMap(data Data, index int32) {
	mapData[data.A] = data
	delete(mapData, index)
}

func TestAndSlice(t *testing.T) {
	data := Data{A: 10000001, B: "username"}
	num := 1
	now := time.Now().UnixNano()
	for i := 0; i < num; i++ {
		andSlice(data, index)
	}
	finish1 := time.Now().UnixNano()
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!testSlice!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(finish1, now, finish1-now, (finish1-now)/1000000)

	now = time.Now().UnixNano()
	for i := 0; i < num; i++ {
		andMap(data, index)
	}
	finish1 = time.Now().UnixNano()
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!testMap!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(finish1, now, finish1-now, (finish1-now)/1000000)
}

func BenchmarkSliceAnd(b *testing.B) {
	data := Data{A: 10000001, B: "username"}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andSlice(data, index)
	}
}

func BenchmarkMapAnd(b *testing.B) {
	data := Data{A: 10000001, B: "username"}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		andMap(data, index)
	}
}
