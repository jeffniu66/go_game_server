package game

//////////////////////////兑换码////////////////////////////////

func CreateCode() {
	//if util.Exists("../script/code.xlsx") {
	//	return
	//}
	//
	//file := xlsx.NewFile()
	//sheet, err := file.AddSheet("Sheet1")
	//f := func(giftId int32, platform string) {
	//	str := "0123456789abcdef"
	//	s := []rune(str)
	//	for _, v := range s {
	//		next := rand.Intn(16)
	//		v, s[next] = s[next], v
	//	}
	//	code := string(s)
	//	db.InsertExchangeCode(code, platform, 0, giftId, 0, 0)
	//	row := sheet.AddRow()
	//	cell := row.AddCell()
	//	cell.Value = code
	//	cell1 := row.AddCell()
	//	cell1.Value = util.ToStr(giftId)
	//	cell2 := row.AddCell()
	//	cell2.Value = platform
	//}
	//
	//length := len(dataConfig.MapCfgExchangeCode)
	//for x := 1; x < length+1; x++ {
	//	num := dataConfig.GetCfgExchangeCode(int32(x)).Num
	//	platform := dataConfig.GetCfgExchangeCode(int32(x)).Sdk
	//	for i := 0; i < int(num); i++ {
	//		f(int32(x), platform)
	//	}
	//}
	//err = file.Save("code.xlsx")
	//if err != nil {
	//	panic(err)
	//}
}
