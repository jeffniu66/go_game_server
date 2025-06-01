package main

import (
	"testing"
)

func DeleteSliceInt(index int32) {
	for i := 0; i < len(sliceInt); i++ {
		if sliceInt[i] == index {
			// if int(index) == len(sliceData)-1 {
			// 	sliceInt = sliceInt[0 : index-1]
			// } else {
			// 	sliceInt = append(sliceInt[:i], sliceInt[i+1:]...)
			// }
			sliceInt = append(sliceInt[:i], sliceInt[i+1:]...)
			break
		}
	}
	// j := 0
	// for i := 0; i < len(sliceInt); i++ {
	// 	if sliceInt[i] != index {
	// 		sliceInt[j] = sliceInt[i]
	// 		j++
	// 	}
	// }
	return
}

func DeleteSlice(index int32) {
	for i := 0; i < len(sliceData); i++ {
		if sliceData[i].A == index {
			sliceData = append(sliceData[:i], sliceData[i+1:]...)
			break
		}
	}

	// j := 0
	// for i := 0; i < len(sliceData); i++ {
	// 	if sliceData[i].A != index {
	// 		sliceData[j] = sliceData[i]
	// 		j++
	// 	}
	// }
	return
}

func DeleteMap(index int32) {
	delete(mapData, index)
	return
}

func BenchmarkSliceIntDel(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DeleteSliceInt(index)
	}
}

func BenchmarkSliceDel(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DeleteSlice(index)
	}
}

func BenchmarkMapDel(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		DeleteMap(index)
	}
}
