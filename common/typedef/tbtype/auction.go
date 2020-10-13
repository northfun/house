package tbtype

type JsonAuctionReview struct {
	Data []*TableAuctionReview `json:"data"`
}

type TableAuctionReview struct {
	Id       uint64 `xorm:"integer" json:"id"`
	MnNotice bool   `xorm:"varchar(128)" json:"mnNotice"`
	Credit   bool   `xorm:"boolean" json:"credit"`
	ItemUrl  string `xorm:"varchar(128)" json:"itemUrl"`
	Status   string `xorm:"varchar(64)" json:"status"`
	Title    string `xorm:"varchar(128)" json:"title"`
	PicUrl   string `xorm:"varchar(128)" json:"picUrl"`
	// price should be string,
	// but for the case of convenient
	InitialPrice      float64 `xorm:"-" json:"initialPrice"`
	CurrentPrice      float64 `xorm:"-" json:"currentPrice"`
	ConsultPrice      float64 `xorm:"-" json:"consultPrice"`
	MarketPrice       float64 `xorm:"-" json:"marketPrice"`
	NInitialPrice     string  `xorm:"numeric(19,4)" json:"-"`
	NCurrentPrice     string  `xorm:"numeric(19,4)" json:"-"`
	NConsultPrice     string  `xorm:"numeric(19,4)" json:"-"`
	NMarketPrice      string  `xorm:"numeric(19,4)" json:"-"`
	SellOff           bool    `xorm:"boolean" json:"sellOff"`
	Start             int64   `xorm:"timestamp" json:"start"`
	End               int64   `xorm:"timestamp" json:"end"`
	TimeToStart       int64   `xorm:"timestamp" json:"timeToStart"`
	TimeToEnd         int64   `xorm:"timestamp" json:"timeToEnd"`
	ViewerCount       uint32  `xorm:"integer" json:"viewerCount"`
	BidCount          uint32  `xorm:"integer" json:"bidCount"`
	DelayCount        uint32  `xorm:"integer" json:"delayCount"`
	ApplyCount        uint32  `xorm:"integer" json:"applyCount"`
	CatNames          string  `xorm:"varchar(128)" json:"catNames"`
	CollateralCatName string  `xorm:"varchar(128)" json:"collateralCatName"`
	XmppVersion       string  `xorm:"varchar(128)" json:"xmppVersion"`
	BuyRestrictions   uint32  `xorm:"varchar(128)" json:"buyRestrictions"`
	SupportLoans      uint32  `xorm:"varchar(128)" json:"supportLoans"`
	SupportOrgLoan    uint32  `xorm:"varchar(128)" json:"supportOrgLoan"`
	TrackParams       string  `xorm:"varchar(128)" json:"trackParams"`
}
