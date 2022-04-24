package utils

import (
	"go.uber.org/zap"
	"os"
	"strings"
)

func NewLogger() *zap.SugaredLogger {
	env := os.Getenv("ENVIRONMENT")
	if strings.Contains("prod", strings.ToLower(env)) {
		logger, _ := zap.NewProduction()
		return logger.Sugar()
	}
	logger, _ := zap.NewDevelopment()
	return logger.Sugar()
}
