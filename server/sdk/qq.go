package sdk

import (
	"go_game_server/server/global"
	"go_game_server/server/logger"
)

type QqSDK struct {
	appID     string
	appSecret string
	verifyUrl string
}

func (q *QqSDK) VerifyAccount(name, token string) (*SDKUserInfo, bool) {
	wxLoginResp := AppletWeChatLogin(token, "qq")
	if wxLoginResp == nil {
		logger.Log.Warnln("wx VerifyAccount AppletWeChatLogin  account get fail!")
		return nil, false
	}
	ret := &SDKUserInfo{}
	ret.Name = wxLoginResp.OpenId
	ret.SessionKey = wxLoginResp.SessionKey
	return ret, true
}

func NewQqSDK2() ServerSDK {
	appID := global.MyConfig.Read("qq", "app_id")
	appSecret := global.MyConfig.Read("qq", "app_secret")
	verifyUrl := global.MyConfig.Read("qq", "verify_url")

	return &QqSDK{
		appID:     appID,
		appSecret: appSecret,
		verifyUrl: verifyUrl}
}

var qqSdk *QqSDK

func NewQqSDK() *QqSDK {
	appID := global.MyConfig.Read("qq", "app_id")
	appSecret := global.MyConfig.Read("qq", "app_secret")
	verifyUrl := global.MyConfig.Read("qq", "verify_url")

	return &QqSDK{
		appID:     appID,
		appSecret: appSecret,
		verifyUrl: verifyUrl}
}
