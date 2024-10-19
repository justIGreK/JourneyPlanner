package logger

import "go.uber.org/zap"

var logger *zap.Logger

func init() {
	var err error
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.StacktraceKey = "" 
    logger, err = loggerConfig.Build()
	if err != nil {
		panic(err)
	}
}

func GetLogger() *zap.Logger {
	return logger
}
