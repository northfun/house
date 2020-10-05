package scraping

import (
	"testing"
	"time"

	"github.com/northfun/house/common/dao/idb"
	"github.com/northfun/house/common/utils/db"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/dao"
	"github.com/stretchr/testify/require"
)

type TestApp struct{}

func TestPageNum(t *testing.T) {
	// dao.WRITE_TO_DB = false

	// conf.C().StartModule.Goods = true
	conf.C().StartModule.Ingredients = true

	c := db.DefaultConf()
	err := db.Init(&c)
	defer db.Close()
	require.Empty(t, err)

	if dao.WRITE_TO_DB &&
		conf.C().StartModule.Goods {
		err = idb.TruncateAllTables(db.DB())
		require.Empty(t, err)

		// err = idb.InitTables(db.DB())
		// require.Empty(t, err)
	}

	var m Manager
	m.Init(&TestApp{})
	err = m.Start()
	require.Empty(t, err)

	time.Sleep(5 * time.Second)
}
