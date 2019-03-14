package db

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// load dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DbMgr DataMgr

type DataMgr interface {
	StoreMgr
	ArtistMgr
	ArtistStoreInfoMgr
	ArtistDetailsMgr
	UserMgr
	SubscriptionMgr
	ReleaseMgr
	LastActionMgr
	ChatMgr
	FeedMgr
	NotificationMgr
	Begin() *AppDatabaseMgr
	Commit() *AppDatabaseMgr
	Rollback() *AppDatabaseMgr
	Close() error
	DropAllTables() error
	Ping() error
}

type AppDatabaseMgr struct {
	db *gorm.DB
}

func NewAppDatabaseMgr(db *gorm.DB) *AppDatabaseMgr {
	return &AppDatabaseMgr{db: db}
}

func NewMainDatabaseMgr() *AppDatabaseMgr {
	db := InitMain()
	if err := CreateAll(db); err != nil {
		panic(errors.Wrap(err, "tried to create all"))
	}
	return NewAppDatabaseMgr(db)
}

func NewFakeDatabaseMgr() *AppDatabaseMgr {
	db := InitFake()
	if err := CreateTables(db); err != nil {
		panic(errors.Wrap(err, "tried to create tables"))
	}
	return NewAppDatabaseMgr(db)
}

func (mgr *AppDatabaseMgr) Begin() *AppDatabaseMgr {
	return &AppDatabaseMgr{db: mgr.db.Begin()}
}

func (mgr *AppDatabaseMgr) Commit() *AppDatabaseMgr {
	return &AppDatabaseMgr{db: mgr.db.Commit()}
}

func (mgr *AppDatabaseMgr) Rollback() *AppDatabaseMgr {
	return &AppDatabaseMgr{db: mgr.db.Rollback()}
}

func (mgr *AppDatabaseMgr) Close() error {
	return mgr.db.Close()
}

func (mgr *AppDatabaseMgr) DropAllTables() error {
	return DropAllTables(mgr.db)
}

func (mgr *AppDatabaseMgr) Ping() error {
	return mgr.db.Exec("select 1").Error
}
