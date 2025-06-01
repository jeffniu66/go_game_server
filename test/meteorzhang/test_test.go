package main

import (
	"fmt"
	"testing"
	"time"
)

var testTimeSlice = []string{"aa", "bb", "cc", "dd", "ee", "jj", "kk", "zz"}

var testTimeMap = map[string]bool{"aa": true, "bb": true, "cc": true, "dd": true, "ee": true, "ff": true, "zz": true}
var testIntMap = map[int]bool{1: true, 2: true, 3: true, 4: true, 5: true, 6: true, 7: true}

//以上为第一组查询测试数据

var testTimeSlice2 = []string{"aa", "bb", "cc", "dd", "ee", "ff", "gg", "hh", "ii", "jj", "kk", "ll", "mm", "nn", "oo", "pp", "qq", "rr", "ss", "tt", "uu", "vv", "ww", "xx", "yy", "zz", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "zz"}

var testTimeMap2 = map[string]bool{"aa": true, "bb": true, "cc": true, "dd": true, "ee": true, "ff": true, "qq": true, "ww": true, "rr": true, "tt": true, "zz": true, "uu": true, "ii": true, "oo": true, "pp": true, "lk": true, "kl": true, "jk": true, "kj": true, "hl": true, "lh": true, "fg": true, "gfdd": true, "df": true, "fd": true,
	"i": true, "j": true, "l": true, "m": true, "n": true, "o": true, "p": true, "q": true, "k": true, "x": true, "y": true, "z": true,
	"1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true, "10": true}

//以上为第二组查询测试数据
func testSlice(a []string) {
	now := time.Now()

	for j := 0; j < 100000; j++ {
		for _, v := range a {
			if v == "zz" {
				break
			}
		}
	}
	finish1 := time.Since(now)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!testSlice!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(finish1)
}

func testMap(a map[string]bool) {
	now := time.Now()
	for j := 0; j < 100000; j++ {
		if _, ok := a["zz"]; ok {
			continue
		}
	}
	finish2 := time.Since(now)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!testMap!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(finish2)
}

func testIntMapf(a map[int]bool) {
	now := time.Now()
	for j := 0; j < 100000; j++ {
		if _, ok := a[6]; ok {
			continue
		}
	}
	finish2 := time.Since(now)
	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!!!!testIntMapf!!!!!!!!!!!!!!!!!!!!!!!!!!!")
	fmt.Println(finish2)
}

func TestT1(t *testing.T) {
	testSlice(testTimeSlice)
	testMap(testTimeMap)
	testIntMapf(testIntMap)
}

func TestT2(t *testing.T) {
	testSlice(testTimeSlice2)
	testMap(testTimeMap2)
}
