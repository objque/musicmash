package main

import (
	"flag"
	"fmt"
	"os"

	sentry "github.com/getsentry/sentry-go"
	"github.com/musicmash/musicmash/internal/api"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/cron"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/fetcher"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier"
)

const (
	configParamName        = "config"
	configParamValue       = ""
	configParamDescription = "Path to musicmash.yaml config"
)

func main() {
	showHelp := flag.Bool("help", false, "Show usage and exit")
	config.Config = config.New()
	config.Config.FlagSet()

	configPath := flag.String(configParamName, configParamValue, configParamDescription)
	flag.Parse()
	if *showHelp {
		flag.PrintDefaults()
		os.Exit(0)
	}
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
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.Config.Sentry.Key,
			AttachStacktrace: true,
			Environment:      config.Config.Sentry.Environment,
		})
		if err != nil {
			fmt.Printf("Sentry initialization failed: %v\n", err)
		}
	}

	log.Info("Running musicmash..")
	go cron.Run(db.ActionFetch, config.Config.Fetching.CountOfSkippedHours, fetcher.Fetch)
	go cron.Run(db.ActionNotify, config.Config.Notifier.CountOfSkippedHours, notifier.Notify)
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}
