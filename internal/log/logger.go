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

// Debugf logs a message at level Debug on the standard logger.
func Debugf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Debugf(format, args...)
}

// Debugln logs a message with fields at level Debug on the standard logger.
func Debugln(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Debugln(format)
}

// Infof logs a message at level Info on the standard logger.
func Infof(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Infof(format, args...)
}

// Infoln logs a message with fields at level Debug on the standard logger.
func Infoln(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Infoln(format)
}

// Info logs a message with fields at level Debug on the standard logger.
func Info(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Infoln(format)
}

// Error logs a message with fields at level Debug on the standard logger.
func Error(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Error(format)
	raven.CaptureMessage(format, nil)
}

// Errorf logs a message at level Error on the standard logger.
func Errorf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Errorf(format, args...)
	raven.CaptureMessage(fmt.Sprintf(format, args...), nil)
}

// Warning logs a message with fields at level Debug on the standard logger.
func Warn(args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format := formatMessageWithFileInfo(sprintlnn(args...))
	entry.Warn(format)
	raven.CaptureMessage(format, nil)
}

// Warningf logs a message at level Warning on the standard logger.
func Warnf(format string, args ...interface{}) {
	entry := logger.WithFields(logrus.Fields{})
	format = formatMessageWithFileInfo(format)
	entry.Warningf(format, args...)
	raven.CaptureMessage(fmt.Sprintf(format, args...), nil)
}

// Panic logs a message at level Panic on the standard logger.
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
