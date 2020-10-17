package idb

import (
	"github.com/northfun/house/common/typedef/tbtype"
	"github.com/northfun/house/common/utils/logger"
	"xorm.io/xorm"
)

const (
	TableNameHouseDealInfo = "table_house_deal_info"
	TableNameAuctionReview = "table_auction_review"
)

var (
	tables = []interface{}{
		&tbtype.TableHouseDealInfo{},
		&tbtype.TableAuctionReview{},
		// &rpb.TableSkinType{},
	}

	tableNames = []string{
		TableNameHouseDealInfo,
	}
)

func InitTables(conn *xorm.Engine) error {
	for i := range tables {
		has, err := conn.IsTableExist(tables[i])
		if err != nil {
			logger.Println("[idb],table,%d,exist", i)
			return err
		}
		if has {
			continue
		}
		if err = conn.CreateTables(tables[i]); err != nil {
			logger.Println("[idb],table,%d,create", i)
			return err
		}

	}
	return nil
}

func TruncateAllTables(conn *xorm.Engine) error {
	for i := range tableNames {
		if _, err := conn.Exec("truncate table " + tableNames[i]); err != nil {
			return err
		}
	}
	return nil
}

func DropAllTables(conn *xorm.Engine) error {
	for i := range tables {
		if err := conn.DropTables(tables[i]); err != nil {
			return err
		}
	}
	return nil
}
