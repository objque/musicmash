package db

import (
	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/pkg/errors"
)

func initDB(dialect, args string, logging bool) *gorm.DB {
	db, err := gorm.Open(dialect, args)
	if err != nil {
		panic(errors.Wrapf(err, "tried to open connection to %s", dialect))
	}

	if logging {
		db = db.LogMode(true)
		db.SetLogger(gorm.Logger{LogWriter: log.GetLogger()})
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	if err = db.Error; err != nil {
		panic(errors.Wrap(err, "Error configure database"))
	}
	return db
}

func InitFake() *gorm.DB {
	return initDB("sqlite3", ":memory:", false)
}

func InitMain() *gorm.DB {
	dialect, args := config.Config.DB.GetConnString()
	return initDB(dialect, args, config.Config.DB.Log)
}
