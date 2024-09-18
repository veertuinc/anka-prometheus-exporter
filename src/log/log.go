package log

import (
	"log/slog"
	"os"
)

var Logger = func() *slog.Logger {
	var logLevel string
	if value, exists := os.LookupEnv("LOG_LEVEL"); exists {
		logLevel = value
	} else {
		logLevel = "INFO"
	}
	var level slog.Level
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	case "fatal":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}()

func Info(message string) {
	Logger.Info(message)
}

func Warn(message string) {
	Logger.Warn(message)
}

func Error(message string) {
	Logger.Error(message)
}

func Fatal(message string) {
	Logger.Error(message)
	os.Exit(1)
}

func Debug(message string) {
	Logger.Debug(message)
}
