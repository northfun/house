package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// import _ "github.com/jinzhu/gorm/dialects/mysql"

var (
	_db *gorm.DB
)

func Init(conf *Config) error {
	var err error
	_db, err = gorm.Open("postgres", conf.ConnectInfo())
	if err != nil {
		return err
	}
	return nil
}

func DB() *gorm.DB {
	return _db
}

func Close() {
	if _db == nil {
		return
	}

	_db.Close()
}
