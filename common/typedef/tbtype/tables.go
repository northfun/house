package tbtype

type TableHouseDealInfo struct {
	Id                 uint32  `xorm:"autoincr"`
	CommunityName      string  `xorm:"varchar(128)"`
	Price              string  `xorm:"varchar(32)"`
	ListedPrice        string  `xorm:"varchar(32) comment('挂牌价格')"`
	AdjustPriceTimes   uint16  `xorm:"smallint"`
	VisitTimes         uint16  `xorm:"smallint"`
	RoomType           string  `xorm:"varchar(32) comment('房屋户型')"`
	Floor              string  `xorm:"varchar(32) comment('所在楼层')"`
	Area               string  `xorm:"varchar(32) comment('建筑面积')"`
	AllArea            string  `xorm:"varchar(32) comment('套内面积')"`
	HouseStructure     string  `xorm:"varchar(32) comment('户型结构')"`
	HouseType          string  `xorm:"varchar(32) comment('')"`
	BuildingType       string  `xorm:"varchar(32) comment('建筑类型')"`
	RoomForward        string  `xorm:"varchar(32) comment('房屋朝向')"`
	BuildYearN         string  `xorm:"varchar(32) comment('建成年代')"`
	DecorationType     string  `xorm:"varchar(32) comment('装修情况')"`
	BuildingStructureN string  `xorm:"varchar(32) comment('建筑结构')"`
	HeatingMode        string  `xorm:"varchar(32) comment('供暖方式')"`
	StairResident      string  `xorm:"varchar(32) comment('梯户比例')"`
	Elevator           string  `xorm:"varchar(16) comment('配备电梯')"`
	TxTypeN            string  `xorm:"varchar(32) comment('交易权属')"`
	SaleTime           string  `xorm:"varchar(32) comment('挂牌时间')"`
	SoldTime           string  `xorm:"varchar(32) comment('成交时间')"`
	HouseLimitYearN    string  `xorm:"varchar(32) comment('房屋年限')"`
	HouseBelong        string  `xorm:"varchar(32) comment('房权所属')"`
	HouseUsingN        string  `xorm:"varchar(32) comment('房屋用途')"`
	LianjianId         string  `xorm:"varchar(16) comment('链家编号')"`
	Pic                []uint8 `xorm:"blob comment('')"`
	PicUrl             string  `xorm:"varchar(256) comment('')"`
	Extra              string  `xorm:"varchar(256) comment('')"`
	Saler              string  `xorm:"varchar(32) comment('交易员')"`
	SalerPhone         string  `xorm:"varchar(64) comment('交易员电话')"`
	Sold               string  `xorm:"varchar(16) comment('')"`
}

type TableBuidingLoc struct {
	Id   uint32 `xorm:"autoincr"`
	Name string `xorm:"varchar(128)"`
}
