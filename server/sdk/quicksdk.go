package sdk

import (
	"encoding/json"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

type QuickSDK struct {
	verifyUrl string
	orderUrl  string
}

type QuickSDKRespInfo struct {
	Msg    string       `json:"message"`
	Status bool         `json:"status"`
	Data   QuickSDKData `json:"data"`
}

type QuickSDKData struct {
	UserData QuickSDKUserData `json:"userData"`
}

type QuickSDKUserData struct {
	BindInfo QuickSDKBindInfo `json:"bindInfo"`
}

type QuickSDKBindInfo struct {
	BindFB         QuickSDKBindFB         `json:"bindFB"`
	BindGameCenter QuickSDKBindGameCenter `json:"bindGameCenter"`
	BindGoogle     QuickSDKBindGoogle     `json:"bindGoogle"`
	BindEmail      QuickSDKBindEmail      `json:"bindEmail"`
	BindNaver      QuickSDKBindNaver      `json:"bindNaver"`
	BindTwitter    QuickSDKBindTwitter    `json:"bindTwitter"`
	BindLine       QuickSDKBindLine       `json:"bindLine"`
	BindVK         QuickSDKBindVK         `json:"bindVK"`
}

type QuickSDKBindFB struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindGameCenter struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindGoogle struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindEmail struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindNaver struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindTwitter struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindLine struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

type QuickSDKBindVK struct {
	IsBind           int32  `json:"isBind"`
	OtherAccountName string `json:"otherAccountName"`
}

func NewQuickSDK() ServerSDK {
	verifyUrl := global.MyConfig.Read("quicksdk", "verify_url")
	orderUrl := global.MyConfig.Read("quicksdk", "order_url")

	return &QuickSDK{
		verifyUrl: verifyUrl,
		orderUrl:  orderUrl}
}

func (q *QuickSDK) VerifyAccount(name, token string) (*SDKUserInfo, bool) {
	body, ok := util.HttpGet(q.verifyUrl, "token="+token+"&uid="+name)
	if !ok {
		logger.Log.Warnln("quick sdk verify account get fail!")
		return nil, false
	}

	info := &QuickSDKRespInfo{}
	err := json.Unmarshal(body, info)
	if err != nil || !info.Status {
		logger.Log.Warnln("quick sdk verify account json fail!", string(body), err, info, name, token)
		return nil, false
	}

	//bindInfo := info.Data.UserData.BindInfo
	userInfo := &SDKUserInfo{Name: "quick_" + name, BindInfos: make(map[int32]*SDKBindInfo)}
	////userInfo.BindInfos[BindFaceBook] = &SDKBindInfo{IsBind: bindInfo.BindFB.IsBind, AccountName: bindInfo.BindFB.OtherAccountName}
	////userInfo.BindInfos[BindGameCenter] = &SDKBindInfo{IsBind: bindInfo.BindGameCenter.IsBind, AccountName: bindInfo.BindGameCenter.OtherAccountName}
	////userInfo.BindInfos[BindGoogle] = &SDKBindInfo{IsBind: bindInfo.BindGoogle.IsBind, AccountName: bindInfo.BindGoogle.OtherAccountName}
	////userInfo.BindInfos[BindEmail] = &SDKBindInfo{IsBind: bindInfo.BindEmail.IsBind, AccountName: bindInfo.BindEmail.OtherAccountName}
	////userInfo.BindInfos[BindNaver] = &SDKBindInfo{IsBind: bindInfo.BindNaver.IsBind, AccountName: bindInfo.BindNaver.OtherAccountName}
	////userInfo.BindInfos[BindTwitter] = &SDKBindInfo{IsBind: bindInfo.BindTwitter.IsBind, AccountName: bindInfo.BindTwitter.OtherAccountName}
	////userInfo.BindInfos[BindLine] = &SDKBindInfo{IsBind: bindInfo.BindLine.IsBind, AccountName: bindInfo.BindLine.OtherAccountName}
	////userInfo.BindInfos[BindVK] = &SDKBindInfo{IsBind: bindInfo.BindVK.IsBind, AccountName: bindInfo.BindVK.OtherAccountName}

	return userInfo, true
}
