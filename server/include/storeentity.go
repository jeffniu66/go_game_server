package include

type Store struct {
	UserId                 int32
	BoxUseAdNum            int32 // 高级宝箱使用广告次数（每天重置）
	BoxLastUseAdTime       int32 // 高级宝箱上次使用广告时间
	MysSkinBuyNum          int32 // 神秘皮肤购买次数（每天6 18点重置）
	MysSkinLastRefreshTime int32 // 上次神秘皮肤刷新时间
	MysSkinChipId          int32 // 待出售神秘皮肤碎片id 当天不重复
	SkinUseAdNum           int32 // 皮肤使用广告次数
	SkinLastRefreshTime    int32 // 皮肤上次刷新时间
	GoldViewAdNum          int32 // 免费金币观看广告数
	GoldLastViewAdTime     int32 // 免费金币观看广告时间
	LastViewAdTime         int32 // 上次看广告时间（秒）
	Update                 int32 // 更新标志，1为需要回写数据库
}
