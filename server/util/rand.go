package util

import (
	"bytes"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandGroup(p ...uint32) int {
	if p == nil {
		panic("args not found")
	}

	r := make([]uint32, len(p))
	for i := 0; i < len(p); i++ {
		if i == 0 {
			r[0] = p[0]
		} else {
			r[i] = r[i-1] + p[i]
		}
	}

	rl := r[len(r)-1]
	if rl == 0 {
		return 0
	}

	rn := uint32(rand.Int63n(int64(rl)))
	for i := 0; i < len(r); i++ {
		if rn < r[i] {
			return i
		}
	}

	panic("bug")
}

// [b1,b2]
func RandInt(b1, b2 int32) int32 {
	if b1 == b2 {
		return b1
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	return int32(rand.Int63n(max-min+1) + min)
}

func RandIntN(b1, b2 int32, n uint32) []int32 {
	if b1 == b2 {
		return []int32{b1}
	}

	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min + 1
	if int64(n) > l {
		n = uint32(l)
	}

	r := make([]int32, n)
	m := make(map[int32]int32)
	for i := uint32(0); i < n; i++ {
		v := int32(rand.Int63n(l) + min)

		if mv, ok := m[v]; ok {
			r[i] = mv
		} else {
			r[i] = v
		}

		lv := int32(l - 1 + min)
		if v != lv {
			if mv, ok := m[lv]; ok {
				m[v] = mv
			} else {
				m[v] = lv
			}
		}

		l--
	}

	return r
}

func RandSliceInt(arr []int32) int32 {
	if len(arr) == 0 {
		return 0
	}

	r := RandInt(0, int32(len(arr))-1)
	return arr[r]
}

func RandNoRepeatIntN(b1, b2 int, n int) []int {
	min, max := int64(b1), int64(b2)
	if min > max {
		min, max = max, min
	}
	l := max - min
	if int64(n) > l {
		n = int(l)
	}

	r := make([]int, n)
	m := make(map[int]int32)
	for i := 0; i < n; i++ {
		v := int(rand.Int63n(l) + min)
		if _, ok := m[v]; ok {
			i--
			continue
		} else {
			r[i] = v
			m[v] = 1
		}
	}

	return r
}

// 随机字符
func RandStr(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := int32(b.Len())
	for i := 0; i < len; i++ {
		num := RandInt(0, length-1)
		container += string(str[num])
	}
	return container
}
