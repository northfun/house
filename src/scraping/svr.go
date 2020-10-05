package scraping

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/utils/icolly"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/common/utils/redis"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/sink"
	"go.uber.org/zap"

	oredis "github.com/go-redis/redis"
	"github.com/gocolly/redisstorage"
)

const (
	ChengjiaoPage = "https://zz.lianjia.com/chengjiao"
)

var (
	DOMAIN_NAME = "https://zz.lianjia.com"
	PAGES       = []string{ChengjiaoPage}
	UserAgent   = []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"}
)

type AppItfc interface{}

type Manager struct {
	app AppItfc

	sp *sink.SinkPool

	mainC   *colly.Collector
	innerC  *colly.Collector
	moreC   *colly.Collector
	detailC *colly.Collector

	storage *redisstorage.Storage
}

func (m *Manager) Init(_app AppItfc) {
	m.app = _app
}

func (m *Manager) Start() error {
	redisC := conf.C().Redis

	var rds *oredis.Client
	var err error
	if len(redisC.Addr) > 0 {
		rds, err = redis.NewClient(&redisC)
		if err != nil {
			return err
		}
	}

	m.sp = sink.NewSinkPool(rds)

	m.sp.Start()
	// create the redis storage
	storage := &redisstorage.Storage{
		Address:  redisC.Addr,
		Password: redisC.Password,
		DB:       redisC.DB,
		Prefix:   "colley storage",
	}

	if conf.C().StartModule.House {
		if err := m.ScrapingHouse(storage); err != nil {
			return err
		}
	}

	logger.Info("[scraping],start")
	return nil
}

func (m *Manager) ScrapingHouse(storage *redisstorage.Storage) error {
	mainQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	innerQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	moreQ, _ := queue.New(
		2, // Number of consumer threads
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	detailQ, _ := queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	limit := &colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 10 * time.Second,
	}

	// add storage to the collector
	if err := icolly.BatchInitCollector(storage,
		UserAgent[0], limit,
		&m.mainC, &m.innerC, &m.moreC, &m.detailC); err != nil {
		return err
	}

	m.MainPage(mainQ, innerQ)
	m.InnerPage(moreQ)
	m.MorePage(moreQ, detailQ)
	m.DetailPage()

	for i := range PAGES {
		mainQ.AddURL(PAGES[i])
	}

	if err := mainQ.Run(m.mainC); err != nil {
		logger.Error("[scraping],run main", zap.Error(err))
		return err
	}
	if err := innerQ.Run(m.innerC); err != nil {
		logger.Error("[scraping],run inner", zap.Error(err))
		return err
	}
	if err := moreQ.Run(m.moreC); err != nil {
		logger.Error("[scraping],run more", zap.Error(err))
		return err
	}
	if err := detailQ.Run(m.detailC); err != nil {
		logger.Error("[scraping],run detail", zap.Error(err))
		return err
	}

	logger.Info("[scraping],start,goods")
	return nil
}

func (m *Manager) Stop() {
	m.sp.Stop()
}
