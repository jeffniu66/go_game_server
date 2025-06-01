package sdk

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"go_game_server/server/global"
	"go_game_server/server/logger"
	"go_game_server/server/util"
	"io"
	"strings"
)

type OPPOUser struct {
	Uid    string `json:"uid"`    // acct_id
	Result int32  `json:"result"` // 1获胜,2失败,3平局
}
type gameComReq struct {
	PkgName    string `json:"pkgName"`
	AppKey     string `json:"appKey"`
	AppSecret  string `json:"appSecret"`
	TableId    string `json:"tableId"`
	TableToken string `json:"tableToken"`
	TimeStamp  string `json:"timeStamp"`
	Sign       string `json:"sign"`
}
type gameStartReq struct {
	gameComReq
	PlayerList string `json:"playerList"`
}

type gameEndReq struct {
	gameComReq
	Result string `json:"result"`
}

type OPPOSDK struct {
	pkgName      string
	appKey       string
	appID        string
	appSecret    string
	userInfoURL  string
	gameStartURL string
	gameEndURL   string
}
type oppoUsrInfo struct {
	UserID        string `json:"userId"`
	UserName      string `json:"userName"`
	Avatar        string `json:"avatar"`
	Sex           string `json:"sex"`
	Location      string `json:"location"`
	Constellation string `json:"constellation"`
	Age           int32  `json:"age"`
}

type oppoResp struct {
	ErrorCode int32       `json:"errorcode"`
	ErrorMsg  string      `json:"errormsg"`
	UserInfo  oppoUsrInfo `json:"userInfo"`
	Data      string      `json:"data"`
}

func NewOppoSDK() ServerSDK {
	pkgName := global.MyConfig.Read("opposdk", "apkname")
	appid := global.MyConfig.Read("opposdk", "app_id")
	appKey := global.MyConfig.Read("opposdk", "app_key")
	appSecret := global.MyConfig.Read("opposdk", "app_secret")
	userInfoURL := global.MyConfig.Read("opposdk", "verify_url")
	gameStartURL := global.MyConfig.Read("opposdk", "game_start_url")
	gameEndURL := global.MyConfig.Read("opposdk", "game_end_url")
	return &OPPOSDK{
		pkgName:      pkgName,
		appID:        appid,
		appKey:       appKey,
		appSecret:    appSecret,
		userInfoURL:  userInfoURL,
		gameStartURL: gameStartURL,
		gameEndURL:   gameEndURL,
	}
}

func (o *OPPOSDK) getVerifyUrlParam(token string) string {
	paramMap := make(map[string]interface{})
	paramMap["pkgName"] = o.pkgName
	paramMap["appKey"] = o.appKey
	paramMap["appSecret"] = o.appSecret
	paramMap["token"] = token
	paramMap["timeStamp"] = util.MilliTime()
	dictStr := util.MapToDictStr(paramMap, "=", "&")
	signStr := o.md5Sign(dictStr)
	ret := dictStr + "&sign=" + signStr
	logger.Log.Info("oppo sign", ret, dictStr)
	return ret
}

func (o *OPPOSDK) setGameComReq(roomId int32) (map[string]interface{}, *gameComReq) {
	paramMap := make(map[string]interface{}, 0)
	paramMap["pkgName"] = o.pkgName
	paramMap["appKey"] = o.appKey
	paramMap["appSecret"] = o.appSecret
	paramMap["timeStamp"] = util.MilliTime()

	ret := &gameComReq{}
	ret.AppKey = o.appKey
	ret.PkgName = o.pkgName
	ret.AppSecret = o.appSecret
	ret.TableId = util.ToStr(roomId)
	ret.TableToken = o.md5Sign(ret.TableId)
	ret.TimeStamp = util.ToInt64Str(paramMap["timeStamp"].(int64))
	return paramMap, ret
}

func (o *OPPOSDK) getGameStartParam(roomId int32, playerList []int32) *gameStartReq {
	paramMap, gameCom := o.setGameComReq(roomId)
	playerStr := ""
	for i, v := range playerList {
		if i == 0 {
			playerStr += util.ToStr(v)
		} else {
			playerStr += "," + util.ToStr(v)
		}
	}
	player := `{"playerlist":[` + playerStr + "]}"
	paramMap["playerList"] = player

	dictStr := util.MapToDictStr(paramMap, "=", "&")
	signStr := o.md5Sign(dictStr)

	ret := &gameStartReq{}
	ret.gameComReq = *gameCom
	ret.PlayerList = player
	ret.Sign = signStr
	return ret
}

func (o *OPPOSDK) getGameEndParam(roomId int32, playerList []int32) *gameEndReq {
	paramMap, gameCom := o.setGameComReq(roomId)
	playerStr := ""
	for i, v := range playerList {
		if i == 0 {
			playerStr += util.ToStr(v)
		} else {
			playerStr += "," + util.ToStr(v)
		}
	}
	player := `{"playerlist":[` + playerStr + "]}"
	paramMap["playerList"] = player
	paramMap["timeStamp"] = util.MilliTime()

	dictStr := util.MapToDictStr(paramMap, "=", "&")
	signStr := o.md5Sign(dictStr)

	ret := &gameEndReq{}
	ret.gameComReq = *gameCom
	ret.Result = player
	ret.Sign = signStr
	return ret
}

func (o *OPPOSDK) md5Sign(msgStr string) string {
	var buf bytes.Buffer
	buf.WriteString(msgStr)
	md5Obj := md5.New()
	_, _ = io.WriteString(md5Obj, buf.String())
	s := fmt.Sprintf("%x", md5Obj.Sum(nil))
	return strings.ToUpper(s)
}

func (o *OPPOSDK) VerifyAccount(name, token string) (*SDKUserInfo, bool) {
	param := string(o.getVerifyUrlParam(token))
	body, ok := util.HttpGet(o.userInfoURL, param)
	if !ok {
		logger.Log.Warnln("oppo verify account get fail!")
		return nil, false
	}

	info := &oppoResp{}
	err := json.Unmarshal(body, info)
	if err != nil || info.ErrorCode != 200 {
		logger.Log.Warnln("oppo verify account json fail!", o.userInfoURL+"?"+param, string(body), info, name, token, err)
		return nil, false
	}

	return &SDKUserInfo{Name: "oppo_" + info.UserInfo.UserID}, true
}

// TODO
func (o *OPPOSDK) OPPOSDKGameStart(roomId int32, playerList []int32) {
	req := o.getGameStartParam(roomId, playerList)
	body, ok := util.HttpPostJson(o.gameStartURL, req)
	if !ok {
		logger.Log.Info("oppo sdk game start failed")
		return
	}
	info := &oppoResp{}
	err := json.Unmarshal([]byte(body), info)
	if err != nil || info.ErrorCode != 200 {
		logger.Log.Warnln("oppo verify account json fail!", o.gameStartURL, string(body), info, req, err)
		return
	}
	return
}

// TODO
func (o *OPPOSDK) OPPOSDKGameEnd(roomId int32, playerList []int32) {
	req := o.getGameEndParam(roomId, playerList)
	body, ok := util.HttpPostJson(o.gameEndURL, req)
	if !ok {
		logger.Log.Info("oppo sdk game end failed")
		return
	}
	info := &oppoResp{}
	err := json.Unmarshal([]byte(body), info)
	if err != nil || info.ErrorCode != 200 {
		logger.Log.Warnln("oppo verify account json fail!", o.gameEndURL, string(body), info, req, err)
		return
	}
	return
}
