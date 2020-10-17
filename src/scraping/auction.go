package scraping

import (
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
	"github.com/northfun/house/common/utils"
	"github.com/northfun/house/common/utils/icolly"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/sink"
	"go.uber.org/zap"
)

type ViewerDetailCollector struct {
	ViewerC *colly.Collector
	DetailC *colly.Collector

	Ck []*http.Cookie
}

func NewViewerDetailCollector(
	storage *redisstorage.Storage,
	rawCookies string) *ViewerDetailCollector {
	vd := &ViewerDetailCollector{
		Ck: utils.GenCookies(rawCookies),
	}

	limit := &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 10 * time.Second,
	}

	icolly.BatchInitCollector(storage,
		UserAgent[0], limit, icolly.OnRequestAliPai,
		&vd.ViewerC, &vd.DetailC)

	return vd
}

func (vd *ViewerDetailCollector) Start(
	doViewer func(vQ, dQ *queue.Queue),
	doDetail func(),
	initUrls []string) error {

	viewerQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	detailQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	if doViewer != nil {
		doViewer(viewerQ, detailQ)

		for i := range initUrls {
			viewerQ.AddURL(initUrls[i])
			vd.ViewerC.SetCookies(initUrls[i], vd.Ck)
		}
	}

	if doDetail != nil {
		doDetail()

		if doViewer == nil {
			for i := range initUrls {
				detailQ.AddURL(initUrls[i])
				vd.DetailC.SetCookies(initUrls[i], vd.Ck)
			}
		}
	}

	if doViewer != nil {
		if err := viewerQ.Run(vd.ViewerC); err != nil {
			logger.Error("[vd],run viewer", zap.Error(err))
			return err
		}
	}

	if doDetail != nil {
		if err := detailQ.Run(vd.DetailC); err != nil {
			logger.Error("[vd],run detail", zap.Error(err))
			return err
		}
	}

	logger.Info("[vd],start")
	return nil
}

func (vd *ViewerDetailCollector) Stop() {}

type AuctionManager struct {
	domainName string
	fromUrls   []string

	sp *sink.SinkPool

	clc *ViewerDetailCollector
}

func NewAuctionManager(
	sp *sink.SinkPool,
	storage *redisstorage.Storage) *AuctionManager {

	rawCookies := conf.C().Cookie

	return &AuctionManager{
		domainName: "",
		sp:         sp,
		clc: NewViewerDetailCollector(
			storage, rawCookies),

		fromUrls: []string{},
	}
}

func (am *AuctionManager) Start() error {
	return am.clc.Start(am.DoViewer,
		am.DoDetail, conf.C().AuctionUrls)
}

func (am *AuctionManager) Stop() {}
