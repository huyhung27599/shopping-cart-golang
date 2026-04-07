package main

import (
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	wd := utils.MustGetWorkkingDir()

	logFile := filepath.Join(wd, "internal", "logs", "app.log")

	 logger.InitLogger(logger.LoggerConfig{
		Level: "info",
		Filename: logFile,
		MaxSize: 1,
		MaxBackups: 5,
		MaxAge: 5,
		Compress: true,
		IsDev: utils.GetEnv("APP_ENV", "development"),
	 })
	 err := godotenv.Load(filepath.Join(wd, ".env"))
	if err != nil {

		logger.Log.Error().Msgf("Failed to load environment variables: %v", err)
	}

	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application, err := app.NewApplication(cfg)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to initialize application")
		return
	}

	// Start server
	if err := application.Run(); err != nil {
		logger.Log.Fatal().Err(err).Msg("Failed to run application")
		return
	}
}


