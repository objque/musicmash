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
	"github.com/musicmash/musicmash/internal/version"
	"github.com/pkg/errors"
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
			exitWithError(err)
		}

		// set not provided flags as config values
		config.Config.FlagReload()
		// override config values with provided flags
		flag.Parse()
	}
	if config.Config.Log.Level == "" {
		config.Config.Log.Level = "info"
	}
	if config.Config.HTTP.Port < 0 || config.Config.HTTP.Port > 65535 {
		exitWithError(errors.New("Invalid port: value should be in range: 0 < value < 65535"))
	}

	log.SetLogFormatter(&log.DefaultFormatter)
	log.ConfigureStdLogger(config.Config.Log.Level, config.Config.Log.File)
	log.Debugf("CLI Args: %v", os.Args[1:])
	log.Debugf("Application configuration: \n%s", config.Config.Dump())

	db.Mgr = db.NewMainDatabaseMgr(config.Config.DB.GetConnString())
	if config.Config.DB.AutoMigrate {
		if config.Config.DB.MigrationsDir == "" {
			exitWithError(errors.New("Auto-migrations is enabled, but path to folder is empty"))
		}

		if err := db.Mgr.ApplyMigrations(config.Config.DB.MigrationsDir); err != nil {
			exitWithError(err)
		}
	}

	if config.Config.Sentry.Enabled {
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              config.Config.Sentry.Key,
			AttachStacktrace: true,
			Environment:      config.Config.Sentry.Environment,
		})
		if err != nil {
			exitWithError(errors.Wrap(err, "Sentry initialization failed"))
		}
	}

	log.Info("Running musicmash..")
	if config.Config.Fetcher.Enabled {
		if config.Config.Fetcher.Delay <= 0 {
			exitWithError(errors.New("Invalid fetcher delay: value should be greater than zero"))
		}
		go cron.Run(db.ActionFetch, config.Config.Fetcher.Delay, fetcher.Fetch)
	}
	if config.Config.Notifier.Enabled {
		if config.Config.Notifier.Delay <= 0 {
			exitWithError(errors.New("Invalid notifier delay: value should be greater than zero"))
		}

		notify, err := notifier.New(config.Config.Notifier.URL)
		if err != nil {
			exitWithError(fmt.Errorf("bad url for notifier: %v", err))
		}
		go cron.Run(notifier.Action, config.Config.Notifier.Delay, notify.Notify)
	}
	log.Panic(api.ListenAndServe(config.Config.HTTP.IP, config.Config.HTTP.Port))
}

func isArgProvided(argName string) bool {
	for _, arg := range os.Args {
		if strings.Contains(arg, argName) {
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

func exitWithError(err error) {
	log.Error(err)
	os.Exit(2)
}
