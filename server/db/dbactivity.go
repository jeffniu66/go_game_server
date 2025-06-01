package db

import (
	"go_game_server/server/util"
)

func InsertExchangeCode(code, platform string, getFlag, giftId, dealTime, userId int32) {
	_, err := DB.Exec("replace into t_exchange_code(code, platform, get_flag, gift_id, deal_time, user_id) "+
		" values(?, ?, ?, ?, ?, ?)",
		code, platform, getFlag, giftId, dealTime, userId)
	util.CheckErr(err)
}

//////////////////////////////////////////////////////////
//type Activity include.Activity
//
//func (activity *Activity) InitData(userId int32) *Activity {
//	rows, err := DB.Query("select * from t_user_activity where user_id = ?", userId)
//	defer rows.Close()
//	util.CheckErr(err)
//
//	var signStr, onlineGiftStr, purchaseLimitStr, monthCardStr, accumulateListStr, fbInviteListStr,
//		blackMarketStr, surveyStr, growthFundStr, rechargeLimitStr, rookieCountStr, happyWeekStr string
//	var spaceRiftByte []byte
//	for rows.Next() {
//		err = rows.Scan(&activity.UserId, &signStr, &activity.SignMonth, &onlineGiftStr, &activity.CommonFlag, &activity.BuyGoodsFlag,
//			&purchaseLimitStr, &monthCardStr, &accumulateListStr, &fbInviteListStr, &blackMarketStr, &activity.PowerAward.IsGet,
//			&spaceRiftByte, &surveyStr, &growthFundStr, &rechargeLimitStr, &activity.DailyGift, &rookieCountStr,
//			&happyWeekStr)
//		util.CheckErr(err)
//
//		_ = json.Unmarshal([]byte(signStr), &activity.SignMap)
//		_ = json.Unmarshal([]byte(onlineGiftStr), activity.OnlineGift)
//		_ = json.Unmarshal([]byte(purchaseLimitStr), &activity.PurchaseLimitMap)
//		_ = json.Unmarshal([]byte(monthCardStr), &activity.MonthCard)
//		_ = json.Unmarshal([]byte(accumulateListStr), &activity.AccumulateList)
//		_ = json.Unmarshal([]byte(fbInviteListStr), &activity.FBInviteMap)
//		_ = json.Unmarshal([]byte(blackMarketStr), &activity.BlackMarket)
//		_ = json.Unmarshal(spaceRiftByte, &activity.SpaceRift)
//		_ = json.Unmarshal([]byte(surveyStr), &activity.SurveyMap)
//		_ = json.Unmarshal([]byte(growthFundStr), &activity.GrowthFundMap)
//		_ = json.Unmarshal([]byte(rechargeLimitStr), &activity.RechargeDailyLimitMap)
//		_ = json.Unmarshal([]byte(rookieCountStr), &activity.RookieCountMap)
//		_ = json.Unmarshal([]byte(happyWeekStr), &activity.HappyWeek)
//	}
//
//	return activity
//}
//
//func (activity *Activity) SaveData(userId int32) interface{} {
//	logger.Log.Infof("----------------------- save activity userId : %d, ", userId)
//	if activity.Update == include.Update {
//		signStr := util.CheckJsonMarshalErr(activity.SignMap)
//		onlineGiftStr := util.CheckJsonMarshalErr(activity.OnlineGift)
//		purchaseLimitStr := util.CheckJsonMarshalErr(activity.PurchaseLimitMap)
//		monthCardStr := util.CheckJsonMarshalErr(activity.MonthCard)
//		accumulateListStr := util.CheckJsonMarshalErr(activity.AccumulateList)
//		fbInviteListStr := util.CheckJsonMarshalErr(activity.FBInviteMap)
//		blackMarketStr := util.CheckJsonMarshalErr(activity.BlackMarket)
//		spaceRiftStr := util.CheckJsonMarshalErr(activity.SpaceRift)
//		surveyStr := util.CheckJsonMarshalErr(activity.SurveyMap)
//		growthFundStr := util.CheckJsonMarshalErr(activity.GrowthFundMap)
//		rechargeLimitStr := util.CheckJsonMarshalErr(activity.RechargeDailyLimitMap)
//		rookieCountStr := util.CheckJsonMarshalErr(activity.RookieCountMap)
//		happyWeekStr := util.CheckJsonMarshalErr(activity.HappyWeek)
//
//		ExecDB(UserDBType, userId, "replace into t_user_activity(user_id, sign_data, sign_month, online_gift, common_flag, buy_goods_flag,"+
//			"purchase_limit, month_card, accumulate_list, fb_invite_list, black_market, power_guarantee_flag, space_rift, survey_list,"+
//			"growth_fund, recharge_daily_limit, daily_gift, rookie_count_limit, happy_week)"+
//			"values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
//			activity.UserId, signStr, activity.SignMonth, onlineGiftStr, activity.CommonFlag, activity.BuyGoodsFlag, purchaseLimitStr,
//			monthCardStr, accumulateListStr, fbInviteListStr, blackMarketStr, activity.PowerAward.IsGet, spaceRiftStr, surveyStr,
//			growthFundStr, rechargeLimitStr, activity.DailyGift, rookieCountStr, happyWeekStr)
//		activity.Update = include.UnUpdate
//	}
//
//	return nil
//}
