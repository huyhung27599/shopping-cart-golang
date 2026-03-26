package main

import (
	"log"
	"os"
	"path/filepath"
	"user-management-api/internal/app"
	"user-management-api/internal/config"
	"user-management-api/internal/utils"
	"user-management-api/pkg/logger"

	"github.com/joho/godotenv"
)

func main() {
	wd := mustGetWorkkingDir()

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
	loadEnv(filepath.Join(wd, ".env"))
	// Initialize configuration
	cfg := config.NewConfig()

	// Initialize application
	application := app.NewApplication(cfg)

	// Start server
	if err := application.Run(); err != nil {
		panic(err)
	}
}

func mustGetWorkkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	return wd
}

func loadEnv(path string) {
	err := godotenv.Load(path)
	if err != nil {

		logger.Log.Error().Msgf("Failed to load environment variables: %v", err)
	}
}
