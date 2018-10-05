package main

import (
	"flag"

	"github.com/objque/musicmash/internal/config"
	"github.com/objque/musicmash/internal/cron"
	"github.com/objque/musicmash/internal/db"
	"github.com/objque/musicmash/internal/log"
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

	db.DbMgr = db.NewMainDatabaseMgr()
}

func main() {
	log.Info("Running musicmash..")
	cron.Run()
}
