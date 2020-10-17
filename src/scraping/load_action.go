package scraping

import (
	"github.com/gocolly/redisstorage"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/dao"
	"github.com/northfun/house/src/sink"
)

type LoadAuctionManager struct {
	domainName string
	fromUrls   []string

	sp *sink.SinkPool

	clc *ViewerDetailCollector
}

func NewLoadAuctionManager(
	sp *sink.SinkPool,
	storage *redisstorage.Storage) *LoadAuctionManager {

	rawCookies := conf.C().Cookie

	return &LoadAuctionManager{
		domainName: "",
		sp:         sp,
		clc: NewViewerDetailCollector(
			storage, rawCookies),

		fromUrls: []string{},
	}
}

func (am *LoadAuctionManager) Start() error {
	urls, err := dao.LoadUnflagedAuctionItems()
	if err != nil {
		return err
	}
	return am.clc.Start(nil,
		am.LoadDoDetail, urls)
}

func (am *LoadAuctionManager) Stop() {}
