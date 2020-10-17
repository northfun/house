package dao

import (
	"github.com/northfun/house/common/dao/idb"
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils/db"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/src/conf"
	"go.uber.org/zap"
	"xorm.io/xorm"
)

func SaveAuctionReview(tc []*tbtype.TableAuctionReview) (err error) {
	if !WRITE_TO_DB {
		return nil
	}

	// 	_, err = db.DB().Table(idb.TableNameCosmetcs).
	// 		Insert(tc)
	var insertNum int64

	if err = db.DoSession(db.DB(), func(ss *xorm.Session) (ierr error) {
		insertNum, ierr = ss.Table(
			idb.TableNameAuctionReview).Insert(tc)
		return
	}); err != nil {
		logger.Warn("[dao],insert auction review", zap.Error(err), zap.Reflect("data", tc))
		return
	}

	logger.Debug("[dao],insert auction review", zap.Int("slc", len(tc)), zap.Int64("islc", insertNum))
	return
}

func SaveSubjectMatterInfo(tc []*tbtype.TableSubjectMatterInfo) (err error) {
	if !WRITE_TO_DB {
		return nil
	}

	idSlc := make([]uint64, len(tc))
	for i := range tc {
		idSlc[i] = tc[i].Id
	}

	// 	_, err = db.DB().Table(idb.TableNameCosmetcs).
	// 		Insert(tc)
	var insertNum, updateNum int64

	if err = db.DoSession(db.DB(), func(ss *xorm.Session) (ierr error) {
		if insertNum, ierr = ss.Table(
			idb.TableNameSubjectMatterInfo).
			Insert(tc); ierr != nil {
			return
		}

		updateNum, ierr = ss.Table(
			idb.TableNameAuctionReview).
			In("id", idSlc).
			Update(map[string]interface{}{
				"flag": 1,
			})
		return
	}); err != nil {
		logger.Warn("[dao],insert auction review", zap.Error(err), zap.Reflect("data", tc))
		return
	}

	logger.Debug("[dao],insert subject matters", zap.Int("slc", len(tc)), zap.Int64("islc", insertNum), zap.Int64("unum", updateNum))
	return
}

func LoadUnflagedAuctionItems() (urls []string, err error) {
	if len(conf.C().TestUrls) > 0 {
		return conf.C().TestUrls, nil
	}
	err = db.DB().Table(idb.TableNameAuctionReview).
		Where("flag=0").
		Cols("item_url").
		Find(&urls)
	return
}
