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
	"github.com/musicmash/musicmash/internal/notifier/telegram"
	tasks "github.com/musicmash/musicmash/internal/tasks/subscriptions"
	"github.com/pkg/errors"
)

func main() {
	configPath := flag.String("config", "/etc/musicmash/musicmash.yaml", "Path to musicmash.yaml config")
	flag.Parse()

	if err := config.InitConfig(*configPath); err != nil {
		panic(err)
	}
	if config.Config.Log.Level == "" {
		config.Config.Log.Level = "info"
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level)

	db.DbMgr = db.NewMainDatabaseMgr()
	tasks.InitWorkerPool()
	telegram.New(config.Config.Notifier.TelegramToken)
	if config.Config.Sentry.Enabled {
		if err := raven.SetDSN(config.Config.Sentry.Key); err != nil {
			panic(errors.Wrap(err, "tried to setup sentry client"))
		}
	}

	log.Info("Running musicmash..")
	go cron.Run(db.ActionFetch, config.Config.Fetching.CountOfSkippedHours, fetcher.Fetch)
	go cron.Run(db.ActionNotify, config.Config.Notifier.CountOfSkippedHours, notifier.Notify)
	if store, ok := config.Config.Stores["deezer"]; ok && store.Fetch {
		go cron.Run(db.ActionReFetch, config.Config.Fetching.RefetchAfterHours, fetcher.ReFetch)
	}
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
