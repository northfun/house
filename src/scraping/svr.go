package scraping

import (
	"time"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/common/utils/redis"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/sink"

	oredis "github.com/go-redis/redis"
	"github.com/gocolly/redisstorage"
)

const (
	ChengjiaoPage = "https://zz.lianjia.com/chengjiao"
)

var (
	PAGES     = []string{ChengjiaoPage}
	UserAgent = []string{"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/72.0.3626.119 Safari/537.36"}
)

type AppItfc interface{}

type Manager struct {
	app AppItfc

	sp *sink.SinkPool

	mainC   *colly.Collector
	detailC *colly.Collector
	sameC   *colly.Collector

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

	detailQ, _ := queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	sameQ, _ := queue.New(
		2,
		&queue.InMemoryQueueStorage{MaxSize: 10000},
	)

	m.mainC = colly.NewCollector(
		colly.UserAgent(UserAgent[0]),
	)
	m.detailC = colly.NewCollector(
		colly.UserAgent(UserAgent[0]),
	)
	m.sameC = colly.NewCollector(
		colly.UserAgent(UserAgent[0]),
	)

	// add storage to the collector
	if conf.C().Store {

		var err error
		err = m.mainC.SetStorage(storage)
		if err != nil {
			return err
		}
		err = m.detailC.SetStorage(storage)
		if err != nil {
			return err
		}
		err = m.sameC.SetStorage(storage)
		if err != nil {
			return err
		}

		limit := &colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 2,
			RandomDelay: 10 * time.Second,
		}
		m.mainC.Limit(limit)
		m.detailC.Limit(limit)
		m.sameC.Limit(limit)
	}

	m.MainPage(PAGES, mainQ, detailQ)
	mainQ.Run(m.mainC)
	detailQ.Run(m.detailC)
	sameQ.Run(m.sameC)

	logger.Info("[scraping],start,goods")
	return nil
}

func (m *Manager) Stop() {
	m.sp.Stop()
}
