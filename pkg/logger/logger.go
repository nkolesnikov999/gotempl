package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Logger represents the application logger
type Logger struct {
	*slog.Logger
}

// NewLogger creates a new structured logger with the specified level
func NewLogger(level string, format string) *Logger {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	case "text":
		handler = slog.NewTextHandler(os.Stdout, opts)
	default:
		// Default to JSON for production
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
	}
}

// NewLoggerWithWriter creates a new logger with a custom writer
func NewLoggerWithWriter(level string, format string, writer io.Writer) *Logger {
	var logLevel slog.Level

	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn", "warning":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var handler slog.Handler
	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug,
	}

	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(writer, opts)
	case "text":
		handler = slog.NewTextHandler(writer, opts)
	default:
		handler = slog.NewJSONHandler(writer, opts)
	}

	logger := slog.New(handler)

	return &Logger{
		Logger: logger,
	}
}

// DefaultLogger creates a default logger with info level and JSON format
func DefaultLogger() *Logger {
	return NewLogger("info", "json")
}
