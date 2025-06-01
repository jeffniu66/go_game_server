package main

import (
	"fmt"
	"runtime"
)

type Data struct {
	A int32
	B string
}

var sliceInt []int32
var sliceData []Data
var mapData map[int32]Data
var length int32 = 10
var index int32 = 9
var cpuNum int = 1

func init() {
	// util.RandInt(0, length-1)
	runtime.GOMAXPROCS(cpuNum)

	sliceInt = make([]int32, length)
	sliceData = make([]Data, length)
	mapData = make(map[int32]Data, length)
	for i := int32(0); i < length; i++ {
		data := Data{A: i, B: "username"}
		sliceInt = append(sliceInt, i)
		sliceData = append(sliceData, data)
		mapData[i] = data
	}
	fmt.Println("sliceDatLen:", len(sliceData), index)
}
