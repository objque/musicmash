package main

import (
	"flag"

	"github.com/musicmash/musicmash/internal/api"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/cron"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/fetcher"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier"
	"github.com/musicmash/musicmash/internal/notifier/telegram"
	tasks "github.com/musicmash/musicmash/internal/tasks/subscriptions"
)

func init() {
	log.SetLogFormatter(&log.DefaultFormatter)
	configPath := flag.String("config", "/etc/musicmash/musicmash.yaml", "Path to musicmash.yaml config")
	logLevel := flag.String("log-level", "info", "log level {debug,info,warning,error}")
	flag.Parse()

	if err := config.InitConfig(*configPath); err != nil {
		panic(err)
	}

	if *logLevel != "info" || config.Config.Log.Level == "" {
		// Priority to command-line
		log.ConfigureStdLogger(*logLevel)
	} else {
		// Priority to config
		if config.Config.Log.Level != "" {
			log.ConfigureStdLogger(config.Config.Log.Level)
		}
	}

	tasks.InitWorkerPool()
	db.DbMgr = db.NewMainDatabaseMgr()
	telegram.New(config.Config.Notifier.TelegramToken)
}

func main() {
	log.Info("Running musicmash..")
	go cron.Run(db.ActionFetch, config.Config.Fetching.CountOfSkippedHoursToFetch, fetcher.Fetch)
	go cron.Run(db.ActionNotify, config.Config.Notifier.CountOfSkippedHours, notifier.Notify)
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
