package sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
)

type U8SDKData struct {
	UserID   int32  `json:"userID"`
	Username string `json:"username"`
}

type U8SDKRespInfo struct {
	Msg   string    `json:"msg"`
	State int32     `json:"state"`
	Data  U8SDKData `json:"data"`
}

type U8VerifyParam struct {
	UserID string `json:"userID"`
	Token  string `json:"token"`
	Sign   string `json:"sign"`
}

type U8SDK struct {
	appID     int32
	appKey    string
	verifyUrl string
	orderUrl  string
}

func NewU8SDK() ServerSDK {
	appID := global.MyConfig.ReadInt32("u8sdk", "app_id")
	appKey := global.MyConfig.Read("u8sdk", "app_key")
	verifyUrl := global.MyConfig.Read("u8sdk", "verify_url")
	orderUrl := global.MyConfig.Read("u8sdk", "order_url")

	return &U8SDK{
		appID:     appID,
		appKey:    appKey,
		verifyUrl: verifyUrl,
		orderUrl:  orderUrl}
}

func (s *U8SDK) VerifyAccount(name, token string) (*SDKUserInfo, bool) {
	body, ok := util.HttpGet(s.verifyUrl, string(s.getVerifyUrlencodedParam(name, token)))
	if !ok {
		logger.Log.Warnln("u8 verify account get fail!")
		return nil, false
	}

	info := &U8SDKRespInfo{}
	err := json.Unmarshal(body, info)
	if err != nil || 1 != info.State {
		logger.Log.Warnln("u8 verify account json fail!", string(body), err, info, name, token)
		return nil, false
	}

	return &SDKUserInfo{Name: "u8_" + util.ToStr(info.Data.UserID)}, true
}

func (s *U8SDK) getVerifyJsonParam(userID, token string) []byte {
	v := &U8VerifyParam{
		UserID: userID,
		Token:  token,
		Sign:   s.md5Sign(userID, token),
	}

	param, _ := json.Marshal(v)

	return param
}

func (s *U8SDK) getVerifyUrlencodedParam(userID, token string) []byte {
	v := url.Values{}
	v.Set("userID", userID)
	v.Set("token", token)
	v.Set("sign", s.md5Sign(userID, token))

	return []byte(v.Encode())
}

func (s *U8SDK) md5Sign(userID, token string) string {
	// md5("userID="+userID+"token="+token+appKey)
	var buf bytes.Buffer
	buf.WriteString("userID=")
	buf.WriteString(userID)
	buf.WriteString("token=")
	buf.WriteString(token)
	buf.WriteString(s.appKey)

	md5Obj := md5.New()
	_, _ = io.WriteString(md5Obj, buf.String())

	return fmt.Sprintf("%x", md5Obj.Sum(nil))
}
