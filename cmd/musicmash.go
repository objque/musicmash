package main

import (
	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
)

func main() {
	// TODO (m.kalinin): replace it with a consul or external cfg
	config.Config = &config.AppConfig{
		DB: config.DBConfig{
			DBType:  "mysql",
			DBHost:  "mariadb",
			DBLogin: "musicmash",
			DBPass:  "musicmash",
			DBName:  "musicmash",
			Log:     false,
		},
		Log: config.LogConfig{
			File:          "musicmash.log",
			Level:         "INFO",
			SyslogEnabled: false,
		},
		Fetching: config.Fetching{
			CountOfSkippedHoursToFetch: 8,
		},
		StoreURL: "https://itunes.apple.com",
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level)
	db.DbMgr = db.NewMainDatabaseMgr()

	log.Info("Hello, from musicmash")
}
