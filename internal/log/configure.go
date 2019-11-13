package log

import (
	"os"

	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02 15:04:05"
)

var DefaultFormatter = logrus.TextFormatter{FullTimestamp: true, TimestampFormat: timeFormat}

func ConfigureStdLogger(logLevel, logPath string) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		Errorf("Cannot parse level %v: %v, setting default loglevel INFO.", logLevel, err)
		lvl = logrus.InfoLevel
	}

	logrus.SetLevel(lvl)
	logger.Level = lvl
	Debugf("Logging level set as %s", logLevel)
	logger.Out = os.Stdout

	if logPath != "" {
		configureFileLogger(logPath)
	}
}

func configureFileLogger(path string) {
	if path == "" {
		return
	}
	hook := lfshook.NewHook(lfshook.PathMap{
		logrus.DebugLevel: path,
		logrus.InfoLevel:  path,
		logrus.ErrorLevel: path,
		logrus.FatalLevel: path,
		logrus.WarnLevel:  path,
	}, nil)
	logger.Hooks.Add(hook)
	Infof("Configured logging to %v", path)
}
