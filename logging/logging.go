package logging

import (
	"log"

	"go.uber.org/zap"
)

var zapLogger *zap.Logger

// **InitZapLogger** Initialises and returns the pointer to a production Zap logger instance which
// can be passed to other functions via dependency injection as recommended by the Uber team. For example:
//
//	utils.InitLogger()
//	logger := utils.GetLogger()
//	defer logger.Sync()
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
