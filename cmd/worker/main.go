package main

import (
	"path/filepath"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func NewWorker(cfg *config.Config)  {
	
}

func main() {
	wd := utils.MustGetWorkkingDir()

	logFile := filepath.Join(wd, "internal", "logs", "app.log")

	logger.InitLogger(logger.LoggerConfig{
		Level:      "info",
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     5,
		Compress:   true,
		IsDev:      utils.GetEnv("APP_ENV", "development"),
	})
	err := godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {

		logger.Log.Error().Msgf("Failed to load environment variables: %v", err)
	}

	// Initialize configuration
	cfg := config.NewConfig()

	NewWorker(cfg)



}