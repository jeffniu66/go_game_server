package util

import (
	"bytes"
	"container/list"
	"crypto/md5"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"go_game_server/server/constant"
	"io"
	"math/rand"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func ToStr(key int32) string {
	id := int(key)
	return strconv.Itoa(id)
}

func ToInt64Str(key int64) string {
	id := int(key)
	return strconv.Itoa(id)
}

func ToInt(str string) int32 {
	if str == "" {
		return 0
	}
	id, err := strconv.Atoi(str)
	CheckErr(err)
	return int32(id)
}

func ToFloat64(str string, bitSize int) float64 {
	if str == "" {
		return 0.0
	}

	val, err := strconv.ParseFloat(str, bitSize)
	CheckErr(err)
	return val
}

func ToUInt(str string) uint32 {
	if str == "" {
		return 0
	}
	id, err := strconv.Atoi(str)
	CheckErr(err)
	return uint32(id)
}

func ToBool(str string) bool {
	b, err := strconv.ParseBool(str)
	CheckErr(err)
	return b
}

func BoolToInt32(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func ToInt64(str string) int64 {
	if str == "" {
		return 0
	}
	id, err := strconv.Atoi(str)
	CheckErr(err)
	return int64(id)
}
func ToUint32(str string) uint32 {
	id, err := strconv.Atoi(str)
	CheckErr(err)
	return uint32(id)
}

// 时间戳（秒）
func UnixTime() int32 {
	return int32(time.Now().Unix())
}

// 时间戳（毫秒）
func MilliTime() int64 {
	return time.Now().UnixNano() / 1e6
}

func DateTimeString(t int64) string {
	return time.Unix(t, 0).String()
}

func GetHour() int32 {
	return int32(time.Now().Hour())
}

func GetMinute() int32 {
	return int32(time.Now().Minute())
}

func GetSecond() int32 {
	return int32(time.Now().Second())
}

func GetWeek() time.Weekday {
	return time.Now().Weekday()
}

// 两个时间是否在同一周
func IsSameWeek(a, b int32) bool {
	tma := time.Unix(int64(a), 0)
	_, wa := tma.ISOWeek()

	tmb := time.Unix(int64(b), 0)
	_, wb := tmb.ISOWeek()
	return wa == wb
}

// 现在到下一个整点的秒数
func GetNextIntPoint() int32 {
	return 60*(60-GetMinute()-1) + (60 - GetSecond())
}

// 某一时刻到下一个整点的时间戳
func GetNextTimeStamp(ts int32, t int32) (nts int32) {
	h := int32(time.Unix(int64(ts), 0).Hour())
	m := int32(time.Unix(int64(ts), 0).Minute())
	s := int32(time.Unix(int64(ts), 0).Second())
	if h < t {
		nts = ts + int32((t-h-1)*constant.HourSecond+(60-m-1)*constant.MinSecond+(60-s))
	} else {
		nts = ts + int32((24-h+t-1)*constant.HourSecond+(60-m-1)*constant.MinSecond+(60-s))
	}
	return
}

// 两个时间戳是否包含某一时刻
func IsTwoTimeStampHasTime(timeStamp1, timestamp2, time int32) bool {
	nextTimeStamp := GetNextTimeStamp(timeStamp1, time)
	if nextTimeStamp < timestamp2 {
		return true
	}
	return false
}

// 两个时间戳是否包含某一时时间戳
func IsTwoTimeHasTime(timeStamp1, timestamp2, time int32) bool {
	return timeStamp1 <= time && time <= timestamp2
}

// 返回当前时间距离指定时间过了几天
func GetOverDay(preStamp int32, hour, minutes, second int) (day int32) {
	day = GetOverDay2(preStamp, UnixTime(), hour, minutes, second)
	return
}

func GetOverDay2(preStamp, nowStamp int32, hour, minutes, second int) (day int32) {
	PrevTm := time.Unix(int64(preStamp), 0)
	d := nowStamp - preStamp

	seconds := 0
	if PrevTm.Hour() < hour {
		seconds = (hour - PrevTm.Hour()) * constant.HourSecond
	} else {
		seconds = ((constant.DayHour - PrevTm.Hour()) + hour) * constant.HourSecond
	}

	seconds = (seconds + minutes*constant.MinSecond + second) - (PrevTm.Minute()*constant.MinSecond + PrevTm.Second())
	if d >= int32(seconds) {
		day += 1
		day += (d - int32(seconds)) / constant.DaySecond
	}

	return
}

//返回某个时间戳到下个几时几分的时间戳
func NextTimeStamp(preStamp int32, nextMin, nextHour int) (timeStamp int32) {
	prevTm := time.Unix(int64(preStamp), 0)

	timeStamp = int32((nextHour-prevTm.Hour())*constant.HourSecond +
		(nextMin-prevTm.Minute())*constant.MinSecond - prevTm.Second())
	if timeStamp <= 0 {
		timeStamp += constant.DaySecond + preStamp
	} else {
		timeStamp += preStamp
	}

	return
}

//返回某个时间戳到下个周几几时几分的时间戳
func NextWeekDayTimeStamp(preStamp int32, nextMin, nextHour int, weekday time.Weekday) (timeStamp int32) {
	prevTm := time.Unix(int64(preStamp), 0)

	timeStamp = int32(int((weekday-prevTm.Weekday())*constant.DaySecond) +
		(nextHour-prevTm.Hour())*constant.HourSecond +
		(nextMin-prevTm.Minute())*constant.MinSecond - prevTm.Second())
	if timeStamp <= 0 {
		timeStamp += constant.WeekSecond + preStamp
	} else {
		timeStamp += preStamp
	}

	return
}

// 获取今天N点N时N分的时间戳
func GetDayTimeStamp(h, m, s int) (timeStamp int32) {
	currentTime := time.Now()
	timeStamp = int32(time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), h, m, s, 0, currentTime.Location()).Unix())
	return
}

// 两个时间戳是否同一天
func IsSameDay(timeStampA, timeStampB int32) bool {
	tmA := time.Unix(int64(timeStampA), 0)
	tmB := time.Unix(int64(timeStampB), 0)
	return tmA.Year() == tmB.Year() && tmA.YearDay() == tmB.YearDay()
}

// 返回某个时刻那天的24点的时间戳
func Get24TimeStamp(timeStamp int32) (result int32) {
	tm := time.Unix(int64(timeStamp), 0)
	result = int32(time.Date(tm.Year(), tm.Month(), tm.Day(), 23, 59, 59, 0, tm.Location()).Unix())
	return
}

// 字符符串转时间戳
func ToTimeStamp(timeStr string) int32 {
	timeLayout := "2006.01.02 15:04:05"
	loc, _ := time.LoadLocation("Local")                           //重要：获取时区
	theTime, err := time.ParseInLocation(timeLayout, timeStr, loc) //使用模板在对应时区转化为time.time类型
	if err != nil {
		return 0
	}
	return int32(theTime.Unix())
}

func GetDateStr(timestamp int64) string {
	//格式化为字符串,tm为Time类型
	tm := time.Unix(timestamp, 0)
	return tm.Format("2006-01-02")
}

func Min(x int32, y int32) int32 {
	if x < y {
		return x
	}
	return y
}

func Max(x int32, y int32) int32 {
	if x > y {
		return x
	}
	return y
}

func Clamp(x int32, min int32, max int32) int32 {
	if x <= min {
		return min
	} else if x >= max {
		return max
	} else {
		return x
	}
}

func HasElem(s interface{}, elem interface{}) bool {
	arrV := reflect.ValueOf(s)
	if arrV.Kind() == reflect.Slice {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}

	return false
}

func testHasElem() {
	fmt.Println("Hello, playground")
	foo := []int{23, 12, 891}
	fmt.Println("foo has 23", HasElem(foo, 23))
	bar := []interface{}{43, "heyy", 12.1, false}
	fmt.Println("bar has 23", HasElem(bar, 23))
	fmt.Println("bar has heyy", HasElem(bar, "heyy"))
}

//结构体转为map
func Struct2Map(obj interface{}) map[string]interface{} {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	var data = make(map[string]interface{})
	for i := 0; i < t.NumField(); i++ {
		data[t.Field(i).Name] = v.Field(i).Interface()
	}
	return data
}

func Contains(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}
	return false, errors.New("not in")
}

func ContainsList(l *list.List, target interface{}) (bool, *list.Element) {
	for e := l.Front(); e != nil; e = e.Next() {
		if e.Value == target {
			return true, e
		}
	}
	return false, nil
}

// slice元素去重
// 通过两重循环过滤重复元素(时间换空间)
func RemoveRepByLoop(slc []int32) []int32 {
	result := []int32{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

// 通过map主键唯一的特性过滤重复元素(空间换时间)
func RemoveRepByMap(slc []int32) []int32 {
	result := []int32{}
	tempMap := map[int32]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// 元素去重 效率第一，如果节省计算时间，则可以采用如下方式
func RemoveRep(slc []int32) []int32 {
	if len(slc) < 1024 {
		// 切片长度小于1024的时候，循环来过滤
		return RemoveRepByLoop(slc)
	} else {
		// 大于的时候，通过map来过滤
		return RemoveRepByMap(slc)
	}
}
func RemoveRepStrByLoop(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

func Abs(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}

// 深拷贝数据，dst和src里面的成员变量必须大写
func DeepCopy(dst, src interface{}) error {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(src); err != nil {
		return err
	}
	return gob.NewDecoder(bytes.NewBuffer(buf.Bytes())).Decode(dst)
}

func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func MergeMapII32(m1, m2 map[int32]int32) (m map[int32]int32) {
	m = make(map[int32]int32)
	for k, v := range m1 {
		m[k] = v
	}
	for k, v := range m2 {
		if value, ok := m[k]; ok {
			m[k] = value + v
		} else {
			m[k] = v
		}
	}
	return
}

// map相减
func SubMapII32(m1, m2 map[int32]int32) (m map[int32]int32) {
	m = m1
	for k, v := range m2 {
		if value, ok := m[k]; ok {
			m[k] = value - v
		} else {
			m[k] = -v
		}
	}
	return
}

func MergeMap(m1, m2 map[int32]int32) (m map[int32]int32) {
	m = m1
	for k, v := range m2 {
		m[k] = v
	}
	return
}

func SliceInt64Reverse(s []int64) {
	if s == nil {
		return
	}
	length := len(s)
	if length == 0 {
		return
	}

	mid := length / 2
	for i := 0; i < mid; i++ {
		s[i], s[length-i-1] = s[length-i-1], s[i]
	}
}

// 获取随机名字
func GetRandName() string {
	str := "abcdefghijkl"
	s := []rune(str)
	for _, v := range s {
		next := rand.Intn(12)
		v, s[next] = s[next], v
	}
	return string(s)
}

// 获取随机名字
func GetRandID() int64 {
	str := "123456789"
	s := []rune(str)
	for _, v := range s {
		next := rand.Intn(len(str))
		v, s[next] = s[next], v
	}
	return ToInt64(string(s))
}

const (
	TileWidth      int32 = 4
	TileHeight     int32 = 2
	HalfTileWidth  int32 = 2
	HalfTileHeight int32 = 1
)

// 逻辑坐标转世界坐标
func TileToWorldPos(x, y int32) (wx, wy int32) {
	wx = x*TileWidth + (y&1)*HalfTileWidth + HalfTileWidth
	wy = y*HalfTileHeight + HalfTileHeight
	fmt.Printf("逻辑转世界:(%d,%d)=>(%d,%d)\n", x, y, wx, wy)
	return
}

// 世界坐标转逻辑坐标
func WorldToTilePos(x, y float64) (tx, ty int32) {
	tx = int32(x / float64(TileWidth))
	ty = int32(y / float64(HalfTileHeight))
	for m := tx - 1; m <= tx+1; m++ {
		for n := ty - 1; n <= ty+1; n++ {
			ox, oy := TileToWorldPos(m, n)
			if isPointInDiamond(ox, oy, int32(x), int32(y)) {
				tx = m
				ty = n
				fmt.Printf("世界转逻辑:(%f,%f)=>(%d,%d)\n", x, y, tx, ty)
				return
			}
		}
	}
	ty = ty - 1
	fmt.Printf("世界转逻辑:(%f,%f)=>(%d,%d)\n", x, y, tx, ty)
	return
}

// 点是否在菱形区域内,菱形的宽高就是 tile的宽高
func isPointInDiamond(ox, oy, px, py int32) bool {
	return Abs(ox-px)*TileHeight+Abs(oy-py)*TileWidth <= TileHeight*TileHeight/2
}

func GenMd5Sign(param ...string) string {
	signKey := ""
	for k, v := range param {
		if k == 0 {
			signKey = v
		} else {
			signKey += "&" + v
		}
	}

	var buf bytes.Buffer
	buf.WriteString(signKey)

	md5Hash := md5.New()
	_, _ = io.WriteString(md5Hash, buf.String())

	return fmt.Sprintf("%x", md5Hash.Sum(nil))
}

func Shuffle(slice []int32) {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	for len(slice) > 0 {
		n := len(slice)
		randIndex := r.Intn(n)
		slice[n-1], slice[randIndex] = slice[randIndex], slice[n-1]
		slice = slice[:n-1]
	}
}

// 地图区域的坐标，返回x、y的最小最大值
// 4 | 3
// 1 | 2
// xMin, xMax, yMin, yMax
func MapAreaCoordinate(area, width, height int32) (int32, int32, int32, int32) {
	if area == 1 {
		return 0, width / 2, 0, height / 2
	} else if area == 2 {
		return width / 2, width, 0, height / 2
	} else if area == 3 {
		return width / 2, width, height / 2, height
	} else {
		return 0, width / 2, height / 2, height
	}
}

// 参数
//  curNum：当前数量
//  maxNum：最大数量
//  stamp：之前开始恢复的时间戳
//  duration：多久恢复一次
//  restoreNum：多久恢复一次的数值
// 返回值
//  newNum：恢复后的数值
//  newStamp：恢复后的时间戳
func CalcValRestore(curNum, maxNum, stamp, duration, restoreNum int32) (newNum, newStamp int32) {
	nowTime := UnixTime()
	newNum, newStamp = curNum, stamp
	if stamp == 0 {
		if curNum < maxNum {
			newStamp = nowTime
		}
	} else {
		if curNum < maxNum {
			perFinishTime := (nowTime - newStamp) / duration
			if perFinishTime > 0 {
				newNum += restoreNum * perFinishTime
				if newNum > maxNum {
					newNum = maxNum
					newStamp = 0
				} else {
					perLeftTime := (nowTime - newStamp) % duration
					newStamp = nowTime - perLeftTime
				}
			}
		} else {
			newStamp = 0
		}
	}

	return
}

// 统计字符数
func CountStrNum(str string) int32 {
	var han, ch int32 = 0, 0
	for _, c := range str {
		if c == ' ' {
			continue
		}
		if unicode.Is(unicode.Han, c) {
			han++
		} else {
			ch++
		}
	}
	return han + ch
}

func StToMap(st interface{}) map[string]interface{} {
	stJson, err := json.Marshal(st)
	if err != nil {
		CheckErr(err)
		return nil
	}
	ret := make(map[string]interface{}, 0)
	err = json.Unmarshal(stJson, &ret)
	if err != nil {
		CheckErr(err)
		return nil
	}
	return ret
}

// mp-map data, e-值等符号, split-分隔符
func MapToDictStr(mp map[string]interface{}, e, sep string) string {
	dictKeyList := make([]string, 0)
	for k := range mp {
		if len(dictKeyList) == 0 {
			dictKeyList = append(dictKeyList, k)
		} else {
			for i := range dictKeyList {
				if strings.Compare(k, dictKeyList[i]) < 1 {
					if i == 0 {
						tmp := []string{k}
						dictKeyList = append(tmp, dictKeyList...)
						fmt.Println("in 1:", k, dictKeyList)
					} else {
						tmp := []string{}
						tmp = append(tmp, dictKeyList[i:]...)
						dictKeyList = append(dictKeyList[:i], k)
						dictKeyList = append(dictKeyList, tmp...)
					}
					break
				}
				if i == len(dictKeyList)-1 {
					dictKeyList = append(dictKeyList, k)
					fmt.Println("in 3:", k, dictKeyList)
					break
				}
			}
		}
	}
	retList := make([]string, 0)
	for _, key := range dictKeyList {
		v := mp[key]
		str := key + e + fmt.Sprintf("%v", v)
		retList = append(retList, str)
	}
	if len(retList) == 0 {
		return ""
	}
	ret := strings.Join(retList, sep)
	return ret
}
