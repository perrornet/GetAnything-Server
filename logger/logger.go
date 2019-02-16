package logger

import "github.com/apsdehal/go-logger"

func NewLogger(name string, level logger.LogLevel) *logger.Logger {
	log, _ := logger.New(name)
	log.SetLogLevel(level)
	return log
}
