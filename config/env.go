package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file")
		return
	}
	log.Println(".env file loaded")
}

func getString(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	i, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return i
}

func getBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

type DatabaseConfig struct {
	Url string
}

func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Url: getString("DATABASE_URL", ""),
	}
}

// GetLogLevel returns the log level from environment variable LOG_LEVEL
// Default is "info" if not set or invalid
func GetLogLevel() string {
	level := getString("LOG_LEVEL", "info")
	validLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
	}

	if !validLevels[level] {
		return "info"
	}
	return level
}

// GetLogFormat returns the log format from environment variable LOG_FORMAT
// Default is "json" if not set or invalid
func GetLogFormat() string {
	format := getString("LOG_FORMAT", "json")
	validFormats := map[string]bool{
		"json": true,
		"text": true,
	}

	if !validFormats[format] {
		return "json"
	}
	return format
}
