package main

import (
	"flag"

	"github.com/musicmash/musicmash/internal/api"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/cron"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/feed"
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
	} else if config.Config.Log.Level != "" {
		// Priority to config
		log.ConfigureStdLogger(config.Config.Log.Level)
	}

	tasks.InitWorkerPool()
	db.DbMgr = db.NewMainDatabaseMgr()
	telegram.New(config.Config.Notifier.TelegramToken)
	feed.Formatter = feed.NewFormatter(config.Config.Rss.Title, config.Config.Rss.Link, config.Config.Rss.Description)
}

func main() {
	log.Info("Running musicmash..")
	go cron.Run(db.ActionReFetch, config.Config.Fetching.RefetchAfterHours, fetcher.ReFetch)
	go cron.Run(db.ActionFetch, config.Config.Fetching.CountOfSkippedHours, fetcher.Fetch)
	go cron.Run(db.ActionNotify, config.Config.Notifier.CountOfSkippedHours, notifier.Notify)
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
