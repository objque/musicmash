package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	db "github.com/objque/musicmash/internal/db/gorm"
)

var DbMgr DataMgr

type DataMgr interface {
	ArtistMgr
	UserMgr
	ReleaseMgr
	LastFetchMgr
	ChatMgr
	SubscriptionMgr
	StateMgr
	FeedMgr
	StoreMgr
	Begin() *AppDatabaseMgr
	Commit() *AppDatabaseMgr
	Rollback() *AppDatabaseMgr
	Close() error
	DropAllTables() error
}
type AppDatabaseMgr struct {
	db *gorm.DB
}

func (mgr *AppDatabaseMgr) Close() error {
	return mgr.db.Close()
}
func NewAppDatabaseMgr(db *gorm.DB) *AppDatabaseMgr {
	return &AppDatabaseMgr{
		db: db,
	}
}

func NewMainDatabaseMgr() *AppDatabaseMgr {
	return NewAppDatabaseMgr(db.InitMain(CreateAll))
}
func NewFakeDatabaseMgr() *AppDatabaseMgr {
	return NewAppDatabaseMgr(db.InitFake(CreateTables))
}

func (mgr *AppDatabaseMgr) Begin() *AppDatabaseMgr {
	return &AppDatabaseMgr{
		db: mgr.db.Begin(),
	}
}

func (mgr *AppDatabaseMgr) Commit() *AppDatabaseMgr {
	return &AppDatabaseMgr{
		db: mgr.db.Commit(),
	}
}

func (mgr *AppDatabaseMgr) Rollback() *AppDatabaseMgr {
	return &AppDatabaseMgr{
		db: mgr.db.Rollback(),
	}
}

func (mgr *AppDatabaseMgr) DropAllTables() error {
	return DropAllTables(mgr.db)
}
