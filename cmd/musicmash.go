package main

import (
	"os"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/fetcher"
	"github.com/objque/musicmash/internal/log"
	"github.com/objque/musicmash/internal/notify"
	"github.com/objque/musicmash/internal/notify/services"
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
			Level:         "DEBUG",
			SyslogEnabled: false,
		},
		Fetching: config.Fetching{
			CountOfSkippedHoursToFetch: 8,
		},
		Store: config.Store{
			URL:    "https://itunes.apple.com",
			Region: "us",
		},
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level)
	db.DbMgr = db.NewMainDatabaseMgr()
	notify.Service = services.New(os.Getenv("TG_TOKEN"))

	log.Info("Running fetching...")
	fetcher.Run()
}
