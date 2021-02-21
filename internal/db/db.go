package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"

	// load dialects
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
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

func NewMainDatabaseMgr(args string) *AppDatabaseMgr {
	db := InitMain(args)
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

func (mgr *AppDatabaseMgr) TruncateAllTables() error {
	const query = `TRUNCATE
		artist_details,
		artist_headers,
		artist_gallery,
		artist_external_links,
		artist_associations,
		stores,
		releases,
		subscriptions,
		artists,
		last_actions
	RESTART IDENTITY CASCADE`

	_, err := mgr.newdb.Exec(query)

	return err
}

func (mgr *AppDatabaseMgr) Close() error {
	return mgr.parent.Close()
}

func (mgr *AppDatabaseMgr) Ping() error {
	return mgr.parent.Ping()
}

func (mgr *AppDatabaseMgr) ApplyMigrations(pathToMigrations string) error {
	databaseInstance, err := postgres.WithInstance(mgr.parent, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("can't create migrate postgres instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrations, "postgres", databaseInstance)
	if err != nil {
		return fmt.Errorf("can't create migrate file driver: %w", err)
	}

	err = m.Up()
	if err != nil && errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("can't apply up migrations: %w", err)
	}

	return nil
}

func (mgr *AppDatabaseMgr) DropAllTablesViaMigrations(pathToMigrations string) error {
	databaseInstance, err := postgres.WithInstance(mgr.parent, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("can't create migrate postgres instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(pathToMigrations, "postgres", databaseInstance)
	if err != nil {
		return fmt.Errorf("can't create migrate file driver: %w", err)
	}

	if err = m.Down(); err != nil {
		return fmt.Errorf("can't apply down migrations: %w", err)
	}

	return nil
}
