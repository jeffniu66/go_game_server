package main

import (
	"testing"
)

// slice 先指数增加1024，随后线性增加

func insertSlice(data Data) {
	sliceData = append(sliceData, data)
	return
}

func insertMap(data Data) {
	mapData[data.A] = data
	return
}

func BenchmarkArrayInt(b *testing.B) {
	data := Data{A: 10000001, B: "username"}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		sliceData[3] = data
	}
}

func BenchmarkSliceInt(b *testing.B) {
	data := Data{A: 10000001, B: "username"}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		insertSlice(data)
	}
}

func BenchmarkMapInt(b *testing.B) {
	data := Data{A: 10000001, B: "username"}
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		insertMap(data)
	}
}
