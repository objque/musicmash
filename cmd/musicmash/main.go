package main

import (
	"flag"

	raven "github.com/getsentry/raven-go"
	"github.com/musicmash/musicmash/internal/api"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/cron"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/fetcher"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier"
	"github.com/pkg/errors"
)

func main() {
	configPath := flag.String("config", "/etc/musicmash/musicmash.yaml", "Path to musicmash.yaml config")
	flag.Parse()

	config.Config = config.New()
	if *configPath != "" {
		if err := config.Config.LoadFromFile(*configPath); err != nil {
			panic(err)
		}
	}
	if config.Config.Log.Level == "" {
		config.Config.Log.Level = "info"
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level)

	db.DbMgr = db.NewMainDatabaseMgr()
	if config.Config.Sentry.Enabled {
		if err := raven.SetDSN(config.Config.Sentry.Key); err != nil {
			panic(errors.Wrap(err, "tried to setup sentry client"))
		}
		raven.SetEnvironment(config.Config.Sentry.Environment)
	}

	log.Info("Running musicmash..")
	go cron.Run(db.ActionFetch, config.Config.Fetching.CountOfSkippedHours, fetcher.Fetch)
	go cron.Run(db.ActionNotify, config.Config.Notifier.CountOfSkippedHours, notifier.Notify)
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
