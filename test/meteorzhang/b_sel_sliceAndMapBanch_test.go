package main

import (
	"testing"
)

func SelectSlice(index int32) *Data {
	for i := 0; i < int(length); i++ {
		if sliceData[i].A == index {
			return &sliceData[i]
		}
	}
	return nil
}

func SelectMap(index int32) *Data {
	d, _ := mapData[index]
	return &d
}

func BenchmarkSliceSel(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SelectSlice(index)
	}
}

func BenchmarkMapSel(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		SelectMap(index)
	}
}
