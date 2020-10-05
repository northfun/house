package dao

import (
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils/db"
	"github.com/northfun/house/common/utils/logger"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

func SaveHouses(tc []*tbtype.TableHouseDealInfo) (err error) {
	if !WRITE_TO_DB {
		return nil
	}

	// 	_, err = db.DB().Table(idb.TableNameCosmetcs).
	// 		Insert(tc)

	err = db.DoSession(db.DB(), func(ss *xorm.Session) (ierr error) {
		for i := range tc {
			_, ierr = ss.Exec("insert into table_cosmetics(brand,name,e_name,volume,price,type,grassed,purchased,like_rate,effect,goods_id,ingredient_ids) values($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)",
				tc[i].Brand,
				tc[i].Name,
				tc[i].EName, tc[i].Volume,
				tc[i].Price, tc[i].Type,
				tc[i].Grassed, tc[i].Purchased,
				tc[i].LikeRate, tc[i].Effect,
				tc[i].GoodsId, tc[i].IngredientIds)
			if ierr != nil {
				return
			}

		}
		return
	})
	if err != nil {
		logger.Warn("[dao],insert cosmetics", zap.Error(err), zap.Reflect("data", tc))
	}
	return
}
