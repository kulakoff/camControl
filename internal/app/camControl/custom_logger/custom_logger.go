package custom_logger

import (
	"log/slog"
	"os"
)

func New(logLevel string) *slog.Logger {
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
	default:
		level = slog.LevelInfo
	}

	handlerOptions := &slog.HandlerOptions{
		Level:     level,
		AddSource: logLevel == "debug", // Добавляем исходный код только для debug
	}
	return slog.New(slog.NewJSONHandler(os.Stdout, handlerOptions))
}
