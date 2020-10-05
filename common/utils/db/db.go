package db

import (
	_ "github.com/lib/pq"
	"xorm.io/core"
	"xorm.io/xorm"
)

// import _ "github.com/jinzhu/gorm/dialects/mysql"

var (
	_db  *xorm.Engine
	_cdb *core.DB
)

func Init(conf *Config) (err error) {
	_db, err = xorm.NewEngine("postgres", conf.ConnectInfo())
	if err != nil {
		return
	}

	_cdb, err = core.Open("postgres", conf.ConnectInfo())
	return
}

func DB() *xorm.Engine {
	return _db
}

func CDB() *core.DB {
	return _cdb
}

func Close() {
	if _db == nil {
		return
	}

	_db.Close()
}

func DoSession(
	conn *xorm.Engine,
	doS func(*xorm.Session) error) (err error) {
	ss := conn.NewSession()
	defer ss.Close()

	if err = ss.Begin(); err != nil {
		return
	}

	if err = doS(ss); err != nil {
		return
	}

	if err = ss.Commit(); err != nil {
		return
	}
	return
}
