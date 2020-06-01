package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/pkg/errors"
)

// TODO (m.kalinin): replace with config values
const (
	maxIdleConns = 10
	maxOpenConns = 100
)

func initDB(dialect, args string) *sqlx.DB {
	db, err := sqlx.Open(dialect, args)
	if err != nil {
		panic(errors.Wrapf(err, "tried to open connection to %s", dialect))
	}

	db.SetMaxIdleConns(maxIdleConns)
	db.SetMaxOpenConns(maxOpenConns)

	return db
}

func InitFake() *sqlx.DB {
	return initDB("sqlite3", ":memory:")
}

func InitMain() *sqlx.DB {
	dialect, args := config.Config.DB.GetConnString()
	return initDB(dialect, args)
}
