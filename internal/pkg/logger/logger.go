package logger

import (
	"go.uber.org/zap"
)

var log *zap.Logger

func Init() error {
	var err error
	log, err = zap.NewProduction()
	if err != nil {
		return err
	}

	zap.ReplaceGlobals(log)
	return nil
}

func L() *zap.Logger {
	if log == nil {
		panic("logger not initialized — call logger.Init() first")
	}
	return log
}

// Sync flushes buffered logs.
func Sync() {
	if log != nil {
		_ = log.Sync()
	}
}
