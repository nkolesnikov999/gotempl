package main

import (
	"context"
	"log/slog"
	"nkpro/gotempl/config"
	"nkpro/gotempl/internal/pages"
	"nkpro/gotempl/internal/users"
	dbpkg "nkpro/gotempl/pkg/db"
	"nkpro/gotempl/pkg/logger"
	"nkpro/gotempl/pkg/middleware"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/gofiber/storage/postgres/v3"
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

	// Init DB pool
	ctx := context.Background()
	pool, err := dbpkg.NewPool(ctx, dbConf.Url)
	if err != nil {
		appLogger.Error("Failed to init DB", slog.String("error", err.Error()))
		return
	}
	defer pool.Close()

	// Session/Storage configuration
	sessionTTLMinutes := config.GetSessionTTLMinutes()
	sessionGCSeconds := config.GetSessionGCSeconds()

	storage := postgres.New(postgres.Config{
		DB:         pool,
		Table:      "sessions",
		Reset:      false,
		GCInterval: time.Duration(sessionGCSeconds) * time.Second,
	})
	store := session.New(session.Config{
		Storage:    storage,
		Expiration: time.Duration(sessionTTLMinutes) * time.Minute,
	})

	// Init repos/services
	userRepo := users.NewPgxRepository(pool)
	userService := users.NewService(userRepo)

	app := fiber.New()

	// Add slog-fiber middleware for HTTP request logging
	app.Use(slogfiber.New(appLogger.Logger))

	app.Use(middleware.AuthMiddleware(store))

	app.Static("/public", "./public")

	// Initialize page handlers
	pages.NewHandler(app, store, pages.WithUserService(userService))

	appLogger.Info("Server starting",
		slog.String("port", ":3003"),
		slog.Int("session_ttl_minutes", sessionTTLMinutes),
		slog.Int("session_gc_seconds", sessionGCSeconds),
	)

	app.Listen(":3003")
}
