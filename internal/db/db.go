package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
	"github.com/musicmash/musicmash/internal/log"
	migrate "github.com/rubenv/sql-migrate"

	// load dialects
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DbMgr DataMgr

type DataMgr interface {
	ArtistMgr
	AssociationMgr
	StoreMgr
	ReleaseMgr
	InternalReleaseMgr
	LastActionMgr
	SubscriptionMgr
	NotificationMgr
	NotificationServiceMgr
	NotificationSettingsMgr
	InternalNotificationMgr
	Begin() *AppDatabaseMgr
	Commit() *AppDatabaseMgr
	Rollback() *AppDatabaseMgr
	Close() error
	Ping() error
	GetDialectName() string
	ApplyMigrations(pathToMigrations string) error
}

type AppDatabaseMgr struct {
	db *gorm.DB
}

func NewAppDatabaseMgr(db *gorm.DB) *AppDatabaseMgr {
	return &AppDatabaseMgr{db: db}
}

func NewMainDatabaseMgr() *AppDatabaseMgr {
	db := InitMain()
	return NewAppDatabaseMgr(db)
}

func NewFakeDatabaseMgr() *AppDatabaseMgr {
	db := InitFake()
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

func (mgr *AppDatabaseMgr) Ping() error {
	return mgr.db.DB().Ping()
}

func (mgr *AppDatabaseMgr) GetDialectName() string {
	return mgr.db.Dialect().GetName()
}

func (mgr *AppDatabaseMgr) ApplyMigrations(pathToMigrations string) error {
	migrations := &migrate.FileMigrationSource{Dir: pathToMigrations}
	n, err := migrate.Exec(mgr.db.DB(), mgr.GetDialectName(), migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Infoln(fmt.Sprintf("Applied %d migrations!", n))
	return nil
}
