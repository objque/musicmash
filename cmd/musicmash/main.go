package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	sentry "github.com/getsentry/sentry-go"
	"github.com/musicmash/musicmash/internal/api"
	"github.com/musicmash/musicmash/internal/config"
	"github.com/musicmash/musicmash/internal/cron"
	"github.com/musicmash/musicmash/internal/db"
	"github.com/musicmash/musicmash/internal/fetcher"
	"github.com/musicmash/musicmash/internal/log"
	"github.com/musicmash/musicmash/internal/notifier"
)

func main() {
	configPath := flag.String("config", "", "Path to musicmash.yaml config")
	if helpRequired() {
		flag.PrintDefaults()
		os.Exit(0)
	}

	config.Config = config.New()
	config.Config.FlagSet()
	flag.Parse()
	if *configPath != "" {
		if err := config.Config.LoadFromFile(*configPath); err != nil {
			panic(err)
		}
	}
	// override config values
	flag.Parse()
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

func helpRequired() bool {
	for _, flag := range os.Args {
		if strings.Contains(flag, "-help") {
			return true
		}
	}
	return false
}
