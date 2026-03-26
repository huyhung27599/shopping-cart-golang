package utils

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"user-management-api/pkg/logger"

	"github.com/rs/zerolog"
)


func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}

func GetIntEnv(key string, defaultValue int) int {
value := os.Getenv(key)
if value == "" {
	return defaultValue
}
intValue, err := strconv.Atoi(value)
if err != nil {
	return defaultValue
}
return intValue

	
}



func NewLoggerWithPath(path string, level string) *zerolog.Logger {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	path = filepath.Join(wd, "internal", "logs", path)

	config := logger.LoggerConfig{
		Level: level,
		Filename: path,
		MaxSize: 1,
		MaxBackups: 5,
		MaxAge: 5,
		Compress: true,
		IsDev: GetEnv("APP_ENV", "development"),
	}
	return logger.NewLogger(config)
}