package dao

import (
	"github.com/northfun/house/common/dao/idb"
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
	var insertNum int64

	if err = db.DoSession(db.DB(), func(ss *xorm.Session) (ierr error) {
		insertNum, ierr = ss.Table(idb.TableNameHouseDealInfo).Insert(tc)
		return
	}); err != nil {
		logger.Warn("[dao],insert cosmetics", zap.Error(err), zap.Reflect("data", tc))
	}

	logger.Debug("[dao],insert house", zap.Int("slc", len(tc)), zap.Int64("islc", insertNum))
	return
}
