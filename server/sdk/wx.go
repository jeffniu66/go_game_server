package sdk

import (
	"encoding/json"
	"fmt"
	"go_game_server/server/global"
	"go_game_server/server/include"
	"go_game_server/server/logger"
	"net/http"
)

var TokenMap = make(map[int32]string) // key: userId	value: openid

type WxSDK struct {
	appID     string
	appSecret string
	verifyUrl string
}

func NewWxSDK2() ServerSDK {
	appID := global.MyConfig.Read("wx", "app_id")
	appSecret := global.MyConfig.Read("wx", "app_secret")
	verifyUrl := global.MyConfig.Read("wx", "verify_url")

	return &WxSDK{
		appID:     appID,
		appSecret: appSecret,
		verifyUrl: verifyUrl}
}

func (w *WxSDK) VerifyAccount(name, token string) (*SDKUserInfo, bool) {
	wxLoginResp := AppletWeChatLogin(token, "wx")
	if wxLoginResp == nil {
		logger.Log.Warnln("wx VerifyAccount AppletWeChatLogin  account get fail!")
		return nil, false
	}
	ret := &SDKUserInfo{}
	ret.Name = wxLoginResp.OpenId
	ret.SessionKey = wxLoginResp.SessionKey
	ret.ErrCode = int32(wxLoginResp.ErrCode)
	return ret, true
}

var wxSdk *WxSDK

func NewWxSDK() *WxSDK {
	appID := global.MyConfig.Read("wx", "app_id")
	appSecret := global.MyConfig.Read("wx", "app_secret")
	verifyUrl := global.MyConfig.Read("wx", "verify_url")

	return &WxSDK{
		appID:     appID,
		appSecret: appSecret,
		verifyUrl: verifyUrl}
}

// 这个函数以 code 作为输入, 返回调用微信接口得到的对象指针和异常情况
func WXLogin(code string, channel string) *include.WXLoginResp {
	// 合成url, 这里的appId和secret是在微信公众平台上获取的
	var url string
	if channel == "wx" {
		url = fmt.Sprintf(wxSdk.verifyUrl, wxSdk.appID, wxSdk.appSecret, code)
	} else if channel == "qq" {
		url = fmt.Sprintf(qqSdk.verifyUrl, qqSdk.appID, qqSdk.appSecret, code)
	}

	// 创建http get请求
	resp, err := http.Get(url)
	if err != nil {
		logger.Log.Errorf("WXLogin httpGet is error: %v", err.Error())
		return nil
	}
	defer resp.Body.Close()

	// 解析http请求中body 数据到我们定义的结构体中
	wxResp := include.WXLoginResp{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&wxResp); err != nil {
		logger.Log.Errorf("WXLogin decode is error: %v", err.Error())
		return nil
	}
	return &wxResp
}

// 微信小程序登录
func AppletWeChatLogin(code string, channel string) *include.WXLoginResp {
	// 根据code获取 openID 和 session_key
	loginResp := WXLogin(code, channel)
	if loginResp == nil {
		logger.Log.Errorln("AppletWeChatLogin loginResp is nil!")
		return nil
	}
	return loginResp
}
