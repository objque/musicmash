package db

import (
	"fmt"
	"os"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/pkg/errors"
)

//nolint:godox
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
	conf := config.DBConfig{
		Host:  os.Getenv("TEST_DB_HOST"),
		Name:  os.Getenv("TEST_DB_NAME"),
		Login: os.Getenv("TEST_DB_USER"),
		Pass:  os.Getenv("TEST_DB_PASSWORD"),
	}

	if port := os.Getenv("TEST_DB_PORT"); port != "" {
		value, err := strconv.Atoi(port)
		if err != nil {
			panic(fmt.Sprintf("can't parse port value '%v' as int", value))
		}

		conf.Port = value
	}

	return InitMain(conf.GetConnString())
}

func InitMain(args string) *sqlx.DB {
	return initDB("postgres", args)
}
