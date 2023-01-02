package logging

import (
	"log"

	"go.uber.org/zap"
)

var zapLogger *zap.Logger

func init() {
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}
	zapLogger = logger
}

func GetLogger() *zap.Logger {
	return zapLogger
}
