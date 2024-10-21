package config

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

var logs *zap.SugaredLogger

func SetLogger(l *zap.Logger) {
	logs = l.Sugar()
}

func LoadEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		logs.Errorf("Error loading .env file: %v", err)

	}
}
