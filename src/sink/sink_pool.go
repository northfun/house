package sink

import (
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/dao"
	"go.uber.org/zap"
)

const (
	CACHE_NUM      = 20
	CACHE_SAVE_DUR = 1 * time.Second
	CACHE_SAVE_INT = 1
)

type SinkPool struct {
	houses chan *tbtype.TableHouseDealInfo

	rcli *redis.Client

	wg   sync.WaitGroup
	stop chan struct{}
}

func NewSinkPool(cli *redis.Client) *SinkPool {
	return &SinkPool{
		houses: make(chan *tbtype.TableHouseDealInfo, 1),

		rcli: cli,

		stop: make(chan struct{}),
	}
}

func (sp *SinkPool) Start() (err error) {
	if conf.C().StartModule.House {
		sp.StartHouse()
	}

	logger.Info("[sinkpool],start")
	return
}

func (sp *SinkPool) StartHouse() {
	sp.goHouse()
}

func (sp *SinkPool) Stop() {
	close(sp.stop)

	sp.wg.Wait()
	logger.Info("[sinkpool],stop")
}

func (sp *SinkPool) goHouse() {
	go func() {
		sp.wg.Add(1)
		defer sp.wg.Done()

		ticker := time.NewTicker(CACHE_SAVE_DUR)
		cacheNum := CACHE_NUM / 2
		ingSlc := make([]*tbtype.TableHouseDealInfo, cacheNum)

		var i int
		var nextSave int64
		for {
			select {
			case <-sp.stop:
				close(sp.houses)

				if i == 0 {
					return
				}
				sp.saveAndResetHouseSlc(ingSlc)
				return
			case ing := <-sp.houses:
				ingSlc[i] = ing
				i++

				if i < cacheNum {
					continue
				}
				nextSave = time.Now().Unix() +
					CACHE_SAVE_INT
				sp.saveAndResetHouseSlc(ingSlc)
				i = 0
			case <-ticker.C:
				if i == 0 {
					continue
				}

				if time.Now().Unix() < nextSave {
					continue
				}

				nextSave = time.Now().Unix() +
					CACHE_SAVE_INT
				sp.saveAndResetHouseSlc(ingSlc[:i])
				i = 0
			}
		}
	}()
}

func (sp *SinkPool) saveAndResetHouseSlc(hSlc []*tbtype.TableHouseDealInfo) (err error) {
	if err = dao.SaveHouses(hSlc); err != nil {
		logger.Warn("[dao],saveHouses", zap.Error(err))
		return
	}

	for i := range hSlc {
		hSlc[i] = nil
	}
	return
}

func (sp *SinkPool) InsertHouse(tc *tbtype.TableHouseDealInfo) {
	sp.houses <- tc
	return
}
