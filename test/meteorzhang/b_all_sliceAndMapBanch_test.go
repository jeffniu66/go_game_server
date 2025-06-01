package main

import (
	"testing"
)

func AllSlice() {
	for i := 0; i < len(sliceData); i++ {
		a := sliceData[i]
		_ = a
	}
	return
}

func AllMap() {
	for k, v := range mapData {
		_, _ = k, v

	}
	return
}

func BenchmarkSliceAll(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AllSlice()
	}
}

func BenchmarkMapAll(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		AllMap()
	}
}
