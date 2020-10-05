package sink

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis"
	"github.com/lib/pq"
	"github.com/northfun/house/common/dao/iredis"
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
		houses: make(chan *rpb.TableHouseDealInfo, 1),

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
				sp.saveAndResetIngSlc(ingSlc)
				return
			case ing := <-sp.ingredients:
				ingSlc[i] = ing
				i++

				if i < cacheNum {
					continue
				}
				nextSave = time.Now().Unix() +
					CACHE_SAVE_INT
				sp.saveAndResetIngSlc(ingSlc)
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
				sp.saveAndResetIngSlc(ingSlc[:i])
				i = 0
			}
		}
	}()
}

func (sp *SinkPool) goCos() {
	go func() {
		sp.wg.Add(1)
		defer sp.wg.Done()

		ticker := time.NewTicker(CACHE_SAVE_DUR)
		cosSlc := make([]*rpb.TableCosmetics, CACHE_NUM)

		var i int
		var nextSave int64
		for {
			select {
			case <-sp.stop:
				close(sp.cosmetcs)

				if i == 0 {
					return
				}
				saveAndResetCosSlc(cosSlc)
				return
			case cos := <-sp.cosmetcs:
				cosSlc[i] = cos
				i++

				if i < CACHE_NUM {
					continue
				}
				nextSave = time.Now().Unix() +
					CACHE_SAVE_INT
				saveAndResetCosSlc(cosSlc)
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
				saveAndResetCosSlc(cosSlc[:i])
				i = 0
			}
		}
	}()
}

func (sp *SinkPool) goCom() {
	go func() {
		sp.wg.Add(1)
		defer sp.wg.Done()

		ticker := time.NewTicker(CACHE_SAVE_DUR)
		comSlc := make([]*rpb.TableComments, CACHE_NUM)

		var i int
		var nextSave int64
		for {
			select {
			case <-sp.stop:
				close(sp.comments)
				if i == 0 {
					return
				}
				saveAndResetComSlc(comSlc)
				return
			case com := <-sp.comments:
				comSlc[i] = com
				i++

				if i < CACHE_NUM {
					continue
				}

				nextSave = time.Now().Unix() +
					CACHE_SAVE_INT
				saveAndResetComSlc(comSlc)
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
				saveAndResetComSlc(comSlc[:i])
				i = 0
			}
		}
	}()
}

func (sp *SinkPool) goUsr() {
	go func() {
		sp.wg.Add(1)
		defer sp.wg.Done()

		ticker := time.NewTicker(CACHE_SAVE_DUR)
		usrSlc := make([]*rpb.TableUsers, CACHE_NUM)
		var nextSave int64

		var i int
		for {
			select {
			case <-sp.stop:
				close(sp.users)
				if i == 0 {
					return
				}
				saveAndResetUsrSlc(usrSlc)
				return
			case usr := <-sp.users:
				usrSlc[i] = usr
				i++

				if i < CACHE_NUM {
					continue
				}

				nextSave = time.Now().Unix() +
					CACHE_SAVE_INT
				saveAndResetUsrSlc(usrSlc)
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
				saveAndResetUsrSlc(usrSlc[:i])
				i = 0
			}
		}
	}()
}

func (sp *SinkPool) saveAndResetIngSlc(ingSlc []*rpb.GoodsIngredients) (err error) {
	var id uint32
	newIngs := make([]*rpb.TableIngredients, 0, len(ingSlc))
	updateGoods := make([]*rpb.GoodsIngredientIds, 0, len(ingSlc))

	for i := range ingSlc {
		gi := rpb.GoodsIngredientIds{
			GoodsId: ingSlc[i].GoodsId,
		}
		var ids pq.Int64Array

		for j := range ingSlc[i].Ingredients {
			if id, err = iredis.GetIngredient(
				sp.rcli,
				ingSlc[i].Ingredients[j].KeyName()); err != nil {
				logger.Warn("[sink],get in redis", zap.Error(err))
				return
			}

			if id > 0 {
				ingSlc[i].Ingredients[j].Id = id
				ids = append(ids, int64(id))
				continue
			}
			// not cached
			id = atomic.AddUint32(
				&sp.ingreIncrId, 1)

			ingSlc[i].Ingredients[j].Id = id

			if err = iredis.AddIngredient(
				sp.rcli,
				ingSlc[i].Ingredients[j].KeyName(),
				id); err != nil {
				logger.Warn("[sink],add in redis", zap.Error(err))
				return
			}

			newIngs = append(newIngs, ingSlc[i].Ingredients[j])
			ids = append(ids, int64(id))
		}

		gi.IngredientIds = ids
		gi.FindName = ingSlc[i].FindName

		updateGoods = append(updateGoods, &gi)
	}

	if err = dao.UpdateGoodsIngredients(updateGoods); err != nil {
		logger.Warn("[dao],updateGoodsIngredients", zap.Error(err))
		return
	}
	if err = dao.SaveIngredients(newIngs); err != nil {
		logger.Warn("[dao],saveIngredients", zap.Error(err))
		return
	}

	for i := range ingSlc {
		ingSlc[i] = nil
	}
	logger.Debug("[sink],insert ing", zap.Int("size", len(ingSlc)))
	return
}

func saveAndResetHouseSlc(hSlc []*tbtype.TableHouseDealInfo) (err error) {
	if err = dao.SaveHouses(cosSlc); err != nil {
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
