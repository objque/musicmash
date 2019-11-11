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
	"github.com/musicmash/musicmash/internal/notifier/telegram"
	"github.com/musicmash/musicmash/internal/version"
)

func main() {
	_ = flag.Bool("version", false, "Show build info and exit")
	if versionRequired() {
		_, _ = fmt.Fprintln(os.Stdout, version.FullInfo)
		os.Exit(0)
	}

	config.Config = config.New()
	config.Config.FlagSet()
	configPath := flag.String("config", "", "Path to musicmash.yaml config")
	_ = flag.Bool("help", false, "Show this message and exit")
	if helpRequired() {
		flag.PrintDefaults()
		os.Exit(0)
	}

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
	if config.Config.HTTP.Port < 0 || config.Config.HTTP.Port > 65535 {
		log.Error("Invalid port: value should be in range: 0 < value < 65535")
		os.Exit(2)
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
	if config.Config.Fetching.Enabled {
		go cron.Run(db.ActionFetch, config.Config.Fetching.Delay, fetcher.Fetch)
	}
	if config.Config.Notifier.Enabled {
		telegram.New(config.Config.Notifier.TelegramToken)
		if err := db.DbMgr.EnsureNotificationServiceExists("telegram"); err != nil {
			log.Error(err)
			os.Exit(2)
		}
		go cron.Run(db.ActionNotify, config.Config.Notifier.Delay, notifier.Notify)
	}
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}

func isArgProvided(argName string) bool {
	for _, flag := range os.Args {
		if strings.Contains(flag, argName) {
			return true
		}
	}
	return false
}

func helpRequired() bool {
	return isArgProvided("-help")
}

func versionRequired() bool {
	return isArgProvided("-version")
}
