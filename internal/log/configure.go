package log

import (
	"log/syslog"
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	logrus_syslog "github.com/sirupsen/logrus/hooks/syslog"
	"github.com/objque/musicmash/internal/config"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var DefaultFormatter = logrus.TextFormatter{FullTimestamp: true, TimestampFormat: timeFormat}

func ConfigureStdLogger(logLevel string) {
	Infof("Applying loglevel %v...", logLevel)

	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Warnf("Cannot parse loglevel %v: %v, setting default loglevel INFO.", logLevel, err)
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	logger.Level = lvl
	logger.Out = os.Stdout

	path := config.Config.Log.File

	if path != "" {
		configureFileLogger(path)
	}

	if config.Config.Log.SyslogEnabled {
		hook, err := logrus_syslog.NewSyslogHook("", "", syslog.LOG_LOCAL0, "musicmash")

		if err == nil {
			logger.Hooks.Add(hook)
		} else {
			Warnf("Failed to configure syslog hook: %v", err)
		}
	}
}

func configureFileLogger(path string) {
	if path == "" {
		return
	}
	hook := lfshook.NewHook(
		lfshook.PathMap{
			logrus.DebugLevel: path,
			logrus.InfoLevel:  path,
			logrus.ErrorLevel: path,
			logrus.FatalLevel: path,
			logrus.WarnLevel:  path,
		}, nil,
	)
	logger.Hooks.Add(hook)
	Infof("Configured logging to %v.", path)
}
