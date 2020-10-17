package tbtype

import (
	"time"
)

type JsonAuctionReview struct {
	Data []*TableAuctionReview `json:"data"`
}

type TableAuctionReview struct {
	Id       uint64 `xorm:"numeric(24,0)" json:"id"`
	MnNotice bool   `xorm:"varchar(8)" json:"mnNotice"`
	Credit   bool   `xorm:"varchar(8)" json:"credit"`
	ItemUrl  string `xorm:"varchar(128)" json:"itemUrl"`
	Status   string `xorm:"varchar(64)" json:"status"`
	Title    string `xorm:"varchar(128)" json:"title"`
	PicUrl   string `xorm:"varchar(128)" json:"picUrl"`
	// price should be string,
	// but for the case of convenient
	InitialPrice float64 `xorm:"numeric(19,4)" json:"initialPrice"`
	CurrentPrice float64 `xorm:"numeric(19,4)" json:"currentPrice"`
	ConsultPrice float64 `xorm:"numeric(19,4)" json:"consultPrice"`
	MarketPrice  float64 `xorm:"numeric(19,4)" json:"marketPrice"`

	SellOff bool `xorm:"varchar(8)" json:"sellOff"`

	Start int64 `xorm:"-" json:"start"`
	End   int64 `xorm:"-" json:"end"`
	// TimeToStart int64 `xorm:"numeric(19,0)" json:"timeToStart"`
	// TimeToEnd   int64 `xorm:"numeric(19,0)" json:"timeToEnd"`

	IStart time.Time `xorm:"timestamp"`
	IEnd   time.Time `xorm:"timestamp"`
	// ITimeToStart time.Time `xorm:"timestamp"`
	// ITimeToEnd   time.Time `xorm:"timestamp"`

	ViewerCount       uint32 `xorm:"integer" json:"viewerCount"`
	BidCount          uint32 `xorm:"integer" json:"bidCount"`
	DelayCount        uint32 `xorm:"integer" json:"delayCount"`
	ApplyCount        uint32 `xorm:"integer" json:"applyCount"`
	CatNames          string `xorm:"varchar(128)" json:"catNames"`
	CollateralCatName string `xorm:"varchar(128)" json:"collateralCatName"`
	XmppVersion       string `xorm:"varchar(128)" json:"xmppVersion"`
	BuyRestrictions   uint32 `xorm:"varchar(128)" json:"buyRestrictions"`
	SupportLoans      uint32 `xorm:"varchar(128)" json:"supportLoans"`
	SupportOrgLoan    uint32 `xorm:"varchar(128)" json:"supportOrgLoan"`
	TrackParams       string `xorm:"varchar(128)" json:"trackParams"`
	Flag              uint32 `xorm:"smallint" `
}

type TableSubjectMatterInfo struct {
	Id                  uint64  `cname:"id" xorm:"numeric(24,0)"`
	InitialPrice        float64 `xorm:"numeric(19,4)"`
	CurrentPrice        float64 `xorm:"numeric(19,4)"`
	ConsultPrice        float64 `xorm:"numeric(19,4)"`
	MarketPrice         float64 `xorm:"numeric(19,4)"`
	PriceRaise          float64 `xorm:"numeric(19,4)"`
	AlarmTimes          uint32  `xorm:"smallint"`
	AuctionTimes        uint32  `cname:"auction_times" xorm:"smallint"`
	GuaranteeDeposit    float64 `xorm:"numeric(19,4)" cname:"deposit" xorm:""`
	Name                string  `xorm:"varchar(128)" cname:"名称" xorm:""`
	PubArea             string  `xorm:"varchar(128)" cname:"公摊" xorm:""`
	PayFirst            string  `xorm:"varchar(128)" cname:"优先购置权" xorm:""`
	Layer               string  `xorm:"varchar(128)" cname:"楼层" xorm:""`
	PriceEvaluatCompany string  `xorm:"varchar(128)" cname:"估价方" xorm:""`
	LandPropertyRight   string  `xorm:"varchar(128)" cname:"土地产权证" xorm:""`
	HousePropertyRight  string  `xorm:"varchar(128)" cname:"房屋产权证" xorm:""`
	RightSource         string  `xorm:"varchar(128)" cname:"权力来源" xorm:""`
	SealUpBy            string  `xorm:"varchar(128)" cname:"查封单位" xorm:""`
	Owner               string  `xorm:"varchar(128)" cname:"所有人" xorm:""`
	HouseUsing          string  `xorm:"varchar(128)" cname:"房屋用途" xorm:""`
	RentStatus          string  `xorm:"varchar(128)" cname:"租赁" xorm:""`
	KeyStatus           string  `xorm:"varchar(128)" cname:"钥匙" xorm:""`
	Facility            string  `xorm:"varchar(128)" cname:"基础设施投资" xorm:""`
	RightLimitStatus    string  `xorm:"varchar(128)" cname:"权力限制" xorm:""`
	RightBookStatus     string  `xorm:"varchar(128)" cname:"权证情况" xorm:""`
	SupportFiles        string  `xorm:"varchar(128)" cname:"文件" xorm:""`
	HouseRightAge       string  `xorm:"varchar(128)" cname:"房产年龄" xorm:""`
	Empty               string  `xorm:"varchar(128)" cname:"已腾空" xorm:""`
	Intro               string  `xorm:"text" cname:"介绍" xorm:""`
	Raw                 string  `xorm:"text" cname:"raw" xorm:""`
}

// func (ts *TableSubjectMatterInfo) SetValue(
// 	from map[string]string) {
//
// 	t := reflect.TypeOf(*ts)
// 	for k, v := range from {
// 		for i := 0; i < t.NumFiled(); i++ {
// 			field := t.Field(i)
// 			tagName := field.Tag.Get("cname")
//
// 			if strings.Countains(tagName){
//
// 			}
// 		}
// 	}
// }
