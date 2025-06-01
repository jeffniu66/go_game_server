package include

type Goods struct {
	GoodsId    int32
	Num        int32
	Type       int32
	SubType    int32
	Location   int32
	Update     int32 // 更新标志，1为需要回写数据库
	CreateTime int32
}
