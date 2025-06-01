package util

import "testing"

func TestRandNoRepeatIntN(t *testing.T) {
	arr := RandNoRepeatIntN(0, 2, 1)
	t.Log(arr)
}

func TestMapToDictStr(t *testing.T) {
	paramMap := make(map[string]interface{})
	paramMap["pkgName"] = "22"
	paramMap["appKey"] = "33"
	paramMap["appSecret"] = "32"
	paramMap["token"] = "11"
	paramMap["timeStamp"] = MilliTime()
	dictStr := MapToDictStr(paramMap, "=", "&")
	t.Log(dictStr)
}
