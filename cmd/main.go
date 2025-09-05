package main

import (
	"log/slog"
	"nkpro/gotempl/config"
	"nkpro/gotempl/internal/pages"
	"nkpro/gotempl/pkg/logger"

	"github.com/gofiber/fiber/v2"
	slogfiber "github.com/samber/slog-fiber"
)

func main() {
	config.Init()

	// Get log configuration from environment
	logLevel := config.GetLogLevel()
	logFormat := config.GetLogFormat()

	// Initialize custom logger with environment configuration
	appLogger := logger.NewLogger(logLevel, logFormat)

	// Set as default slog logger for the application
	slog.SetDefault(appLogger.Logger)

	// Log application startup
	appLogger.Info("Starting application",
		slog.String("service", "go-templ"),
		slog.String("version", "1.0.0"),
		slog.String("log_level", logLevel),
		slog.String("log_format", logFormat),
	)

	dbConf := config.NewDatabaseConfig()

	// Log database configuration (be careful with sensitive data)
	appLogger.Info("Database configuration loaded",
		slog.String("url", dbConf.Url),
	)

	app := fiber.New()

	// Add slog-fiber middleware for HTTP request logging
	app.Use(slogfiber.New(appLogger.Logger))

	// Initialize page handlers
	pages.NewHandler(app)

	appLogger.Info("Server starting",
		slog.String("port", ":3003"),
	)

	app.Listen(":3003")
}
