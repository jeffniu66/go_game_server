package util

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"go_game_server/server/logger"
)

const (
	ContentTypeJson       = "application/json;charset=utf-8"
	ContentTypeUrlencoded = "application/x-www-form-urlencoded;charset=utf-8"
)

func HttpPostJson(url string, msg interface{}) (string, bool) {
	if req, err := json.Marshal(msg); err == nil {
		resp, err := http.Post(url, ContentTypeJson, bytes.NewReader(req))
		if err == nil {
			defer resp.Body.Close()

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				logger.Log.Warnf("HttpPostJson ReadAll err: %v, msg: %v", err, string(body))
			} else {
				return string(body), true
			}
		} else {
			logger.Log.Warnf("HttpPostJson http.Post err: %v, url: %v, msg: %v", err, url, msg)
		}
	} else {
		logger.Log.Warnf("HttpPostJson json.Marshal err: %v, msg: %v", err, msg)
	}

	return "", false
}

func HttpPostUrlencoded(url string, msg string) ([]byte, bool) {
	resp, err := http.Post(url, ContentTypeUrlencoded, bytes.NewReader([]byte(msg)))
	if err == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Warnf("HttpPostUrlencoded ReadAll err: %v, msg: %v", err, string(body))
		} else {
			return body, true
		}
	} else {
		logger.Log.Warnf("HttpPostUrlencoded Post err: %v, msg: %v", err, resp)
	}

	return nil, false
}

func HttpGet(url string, param string) ([]byte, bool) {
	resp, err := http.Get(url + "?" + param)
	if err == nil {
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Log.Warnf("HttpGet ReadAll err: %v, msg: %v", err, string(body))
		} else {
			return body, true
		}
	}

	return nil, false
}

//func postForm() {
//	bodyType := "application/x-www-form-urlencoded"
//	params := url.Values{}
//
//	params.Set("ename", "registe")
//	params.Set("app_id", "1")
//	params.Set("channel_id", "111")
//	params.Set("userid", "222")
//	params.Set("roleid", "222")
//	params.Set("ip", "1")
//	params.Set("logtime", "1234557")
//	params.Set("server_id", "1")
//	params.Set("phone", "123456")
//	params.Set("device_id", "123")
//	params.Set("device_type", "effd")
//	params.Set("app_id", "1")
//	params.Set("imei", "werew")
//	params.Set("username", "gggg")
//
//	http.Post(urlStr, bodyType, strings.NewReader(params.Encode()))
//
//	resp, err := http.PostForm("http://192.168.254.251:8001/api/reports/", params)
//	body, err1 := ioutil.ReadAll(resp.Body)
//	logger.Log.Infof("err: %v, err1: %v, body: %v, msg: %v", err, err1, body, string(body))
//}
//
//func jsonForm()  {
//	bodyType := "application/json"
//	registeMsg := &include.RegisteMsg{Ename:"registe", AppID:10, ChannelID:1001, UserID:123, RoleID:123,
//		Ip:"127.0.0.1", LogTime:UnixTime(), ServerID:"1", Phone:"", DeviceID:"", DeviceType:"", Imei:"", UserName:"aaa"}
//
//	if req, err := json.Marshal(registeMsg); err == nil {
//
//		fmt.Println("aaaaa: ", req, string(req))
//
//		resp, err := http.Post(urlStr, bodyType, bytes.NewReader(req))
//
//
//		body, err1 := ioutil.ReadAll(resp.Body)
//		logger.Log.Infof("err: %v, err1: %v, body: %v, msg: %v", err, err1, body, string(body))
//
//	}
//}
