package config

import (
	logger "JourneyPlanner/pkg/log"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func LoadEnv() {
	logger := logger.GetLogger()
	err := godotenv.Load("../.env")
	if err != nil {
		logger.Fatal("Error loading .env file", zap.Error(err))
	}
}
