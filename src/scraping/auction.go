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

	doViewer(viewerQ, detailQ)
	doDetail()

	for i := range initUrls {
		viewerQ.AddURL(initUrls[i])
		vd.ViewerC.SetCookies(initUrls[i], vd.Ck)
	}

	if err := viewerQ.Run(vd.ViewerC); err != nil {
		logger.Error("[vd],run viewer", zap.Error(err))
		return err
	}
	if err := detailQ.Run(vd.DetailC); err != nil {
		logger.Error("[vd],run detail", zap.Error(err))
		return err
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

	rawCookies := ""

	return &AuctionManager{
		domainName: "",
		sp:         sp,
		clc: NewViewerDetailCollector(
			storage, rawCookies),

		fromUrls: []string{
			// zhongyuan
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410102.htm?spm=a219w.7474998.filter.57.35743c54koWIz8&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// erqi
			"ihttps://zc-paimai.taobao.com/list/0_______56950002_6_410103.htm?spm=a219w.7474998.filter.58.19c33c546W8Wgk&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// guancheng
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410104.htm?spm=a219w.7474998.filter.59.3f943c54GxfItm&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// jinshui
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410105.htm?spm=a219w.7474998.filter.60.4ce13c542FWsY8&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// shangjie
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410106.htm?spm=a219w.7474998.filter.61.664c3c54mq24MX&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// huiji
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410108.htm?spm=a219w.7474998.filter.62.2e043c5480GlQi&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// zhengdong
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410186.htm?spm=a219w.7474998.filter.69.49d33c54SFgdlZ&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
			// gaoxin
			"https://zc-paimai.taobao.com/list/0_______56950002_6_410187.htm?spm=a219w.7474998.filter.70.40e73c54O21KvA&auction_source=0&sorder=2&st_param=-1&auction_start_seg=-1",
		},
	}
}

func (am *AuctionManager) Start() error {
	return am.clc.Start(am.DoViewer,
		am.DoDetail, am.fromUrls)
}

func (am *AuctionManager) Stop() {}
