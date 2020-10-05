package app

import (
	"github.com/northfun/house/common/dao/idb"
	"github.com/northfun/house/common/utils/db"
	"github.com/northfun/house/common/utils/logger"
	"github.com/northfun/house/src/conf"
	"github.com/northfun/house/src/dao"
	"github.com/northfun/house/src/scraping"
	"go.uber.org/zap"
)

type App struct {
	sm *scraping.Manager
}

func NewApp() *App {
	return &App{}
}

func (a *App) Init(cpath string) {
	if err := conf.Init(cpath); err != nil {
		logger.Println("[ERR],init cpath", err)
		return
	}

	config := conf.C()
	dao.WRITE_TO_DB = config.Store

	if err := logger.Init(&config.Log); err != nil {
		logger.Println("[ERR],init log", err)
		return
	}

	if err := db.Init(&config.DB); err != nil {
		logger.Println("[ERR],init db", err)
		return
	}

	if err := idb.InitTables(db.DB()); err != nil {
		logger.Println("[ERR],init tables", err)
		return
	}

	a.sm = &scraping.Manager{}
	a.sm.Init(a)

	logger.Info("[app],init ok", zap.Reflect("c", config))
}

func (a *App) Start() {
	if err := a.sm.Start(); err != nil {
		logger.Warn("[ERR],sm start", zap.Error(err))
	}
}

func (a *App) Stop() {
	a.sm.Stop()

	db.Close()

	logger.Sync()
}
