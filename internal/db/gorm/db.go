package gorm

import (
	"fmt"
	"sync/atomic"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/log"
)

var mainDB *gorm.DB
var panicInit int32

type postCreateFunc func(*gorm.DB) error

func InitFake(postCreate postCreateFunc) *gorm.DB {
	mainDB, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		panic(fmt.Sprintf("Can't open connection to %s %v", "sqlite3", err))
	}

	mainDB = mainDB.LogMode(false)
	mainDB.SetLogger(gorm.Logger{mainDBLogger{}})
	mainDB.DB().SetMaxIdleConns(10)
	mainDB.DB().SetMaxOpenConns(100)
	if err = mainDB.Error; err != nil {
		panic("Error configure database: " + err.Error())
	}
	if err := postCreate(mainDB); err != nil {
		panic(err)
	}
	return mainDB
}

func InitMain(postCreate postCreateFunc) *gorm.DB {
	if !atomic.CompareAndSwapInt32(&panicInit, 0, 1) {
		panic("Attempt to init main database twice!")
	}

	dbtype, connstring := config.Config.DB.GetConnString()
	db, err := gorm.Open(dbtype, connstring)
	if err != nil {
		panic("Can't open database connection: " + err.Error())
	}
	mainDB = db

	if config.Config.DB.Log {
		mainDB = mainDB.LogMode(true)
		mainDB.SetLogger(gorm.Logger{mainDBLogger{}})
	}
	if err := postCreate(mainDB); err != nil {
		panic(err)
	}
	return mainDB
}

type mainDBLogger struct {
}

func (l mainDBLogger) Println(v ...interface{}) {
	log.Info(v...)
}
