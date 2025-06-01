package sdk

import (
	"go_game_server/server/global"
	"go_game_server/server/logger"
)

type SDKChannel string

const (
	SDKOPPO  SDKChannel = "oppo"
	SDKU8    SDKChannel = "U8"
	SDKQUICK SDKChannel = "QUICK"
	SDKWX    SDKChannel = "wx"
	SDKQQ    SDKChannel = "qq"
)

type ServerSDK interface {
	VerifyAccount(name, token string) (*SDKUserInfo, bool)
}

type SDKUtil struct {
	verify bool
	sdkMap map[SDKChannel]ServerSDK
}

type SDKUserInfo struct {
	Name       string                 // openid 其他渠道统一id wx-openID oppo-userID
	BindInfos  map[int32]*SDKBindInfo // 暂时未使用，预留功能
	SessionKey string                 // 渠道返回，各渠道可整合复用
	ErrCode    int32                  // 登录错误码，渠道返回，各渠道可整合复用
}

type SDKBindInfo struct {
	IsBind      int32
	AccountName string
}

var SdkUtil SDKUtil

func NewSDKUtil() SDKUtil {
	verify := global.MyConfig.ReadBool("sdk", "verify")
	sdkKind := global.MyConfig.Read("sdk", "sdk_kind")
	//logger.Log.Warnf(">>>>>>>>>>>> sdk util config! verify:%v, sdk kind:%v", verify, sdkKind)
	SdkUtil = SDKUtil{verify: verify}
	SdkUtil.sdkMap = make(map[SDKChannel]ServerSDK)
	// login and pay
	switch sdkKind {
	case "u8":
		SdkUtil.sdkMap[SDKU8] = NewU8SDK()
	case "quicksdk":
		SdkUtil.sdkMap[SDKQUICK] = NewQuickSDK()
	case "opposdk":
	default:
		logger.Log.Warnln("SdkUtil sdk kind type err:", sdkKind)
	}
	wxSdk = NewWxSDK()
	qqSdk = NewQqSDK()
	SdkUtil.sdkMap[SDKOPPO] = NewOppoSDK()
	SdkUtil.sdkMap[SDKWX] = NewWxSDK2()
	SdkUtil.sdkMap[SDKQQ] = NewQqSDK2()
	return SdkUtil
}

func (s *SDKUtil) GetSDKUserInfo(sdk SDKChannel, name, token string) (*SDKUserInfo, bool) {
	if !s.verify {
		return &SDKUserInfo{Name: name}, true
	}

	// u3d才能以gm身份登录
	if global.OrderConfig.CanUseGm(name) && token == "2^8f$" {
		logger.Log.Warnln("gm login success: ", name)
		return &SDKUserInfo{Name: name}, true
	}

	if s.sdkMap[sdk] == nil {
		return nil, false
	}

	return s.sdkMap[sdk].VerifyAccount(name, token)
}
