package db

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/musicmash/musicmash/internal/log"
	migrate "github.com/rubenv/sql-migrate"

	// load dialects
	_ "github.com/mattn/go-sqlite3"
)

var Mgr *AppDatabaseMgr

type AppDatabaseMgr struct {
	newdb  SQLCommon
	parent *sql.DB
}

type SQLCommon interface {
	QueryRow(query string, args ...interface{}) *sql.Row
	Get(dest interface{}, query string, args ...interface{}) error
	Select(dest interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	NamedExec(string, interface{}) (sql.Result, error)
}

type sqlDb interface {
	Beginx() (*sqlx.Tx, error)
}

type sqlTx interface {
	Commit() error
	Rollback() error
}

func NewAppDatabaseMgr(db *sqlx.DB) *AppDatabaseMgr {
	return &AppDatabaseMgr{newdb: db, parent: db.DB}
}

func NewMainDatabaseMgr() *AppDatabaseMgr {
	db := InitMain()
	return NewAppDatabaseMgr(db)
}

func NewFakeDatabaseMgr() *AppDatabaseMgr {
	db := InitFake()
	return NewAppDatabaseMgr(db)
}

func (mgr *AppDatabaseMgr) Begin() (*AppDatabaseMgr, error) {
	if conn, ok := mgr.newdb.(sqlDb); ok && conn != nil {
		tx, err := conn.Beginx()
		if err != nil {
			return nil, err
		}

		return &AppDatabaseMgr{newdb: interface{}(tx).(SQLCommon), parent: mgr.parent}, nil
	}

	return nil, fmt.Errorf("error begin: %w", ErrAlreadyInTx)
}

func (mgr *AppDatabaseMgr) Commit() error {
	if tx, ok := mgr.newdb.(sqlTx); ok && tx != nil {
		return tx.Commit()
	}

	return fmt.Errorf("error commit: %w", ErrNotInTx)
}

func (mgr *AppDatabaseMgr) Rollback() error {
	if tx, ok := mgr.newdb.(sqlTx); ok && tx != nil {
		return tx.Rollback()
	}

	return fmt.Errorf("error rollback: %w", ErrNotInTx)
}

func (mgr *AppDatabaseMgr) Close() error {
	return mgr.parent.Close()
}

func (mgr *AppDatabaseMgr) Ping() error {
	return mgr.parent.Ping()
}

func (mgr *AppDatabaseMgr) GetDialectName() string {
	return "sqlite3"
}

func (mgr *AppDatabaseMgr) ApplyMigrations(pathToMigrations string) error {
	migrations := &migrate.FileMigrationSource{Dir: pathToMigrations}
	n, err := migrate.Exec(mgr.parent, mgr.GetDialectName(), migrations, migrate.Up)
	if err != nil {
		return err
	}

	log.Infoln(fmt.Sprintf("Applied %d migrations!", n))
	return nil
}
