package log

import (
	"fmt"
	"runtime"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/sirupsen/logrus"
)

var (
	logger        = logrus.New()
	logFileSearch = "/musicmash/"
)

func SetLogLevel(level logrus.Level) {
	logger.Level = level
}

func SetLogFormatter(formatter logrus.Formatter) {
	logger.Formatter = formatter
}

func GetLogger() *logrus.Logger {
	return logger
}

func formatMessageWithFileInfo(msg string) string {
	res := fmt.Sprintf("[%v] %v", fileInfo(3), msg)
	return res
}

func Debugf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Debugf(format, args...)
}

func Debugln(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Debugln(format)
}

func Infof(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Infof(format, args...)
}

func Infoln(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Infoln(format)
}

func Info(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Infoln(format)
}

func Error(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Error(format)
	raven.CaptureMessage(format, nil)
}

func Errorf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Errorf(format, args...)
	raven.CaptureMessage(fmt.Sprintf(format, args...), nil)
}

func Warn(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Warn(format)
	raven.CaptureMessage(format, nil)
}

func Warnf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Warningf(format, args...)
	raven.CaptureMessage(fmt.Sprintf(format, args...), nil)
}

func Panic(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	raven.CaptureMessageAndWait(format, nil)
	entry.Panic(format)
}

func sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func fileInfo(skip int) string {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		file = "<???>"
		line = 1
	} else {
		slash := strings.LastIndex(file, logFileSearch)
		if slash >= 0 {
			file = file[slash+1:]
		}
	}
	return fmt.Sprintf("%s:%d", file, line)
}
