package utilities

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func NewLogger(_debug bool) {
	var logLevel slog.Level
	if _debug {
		logLevel = slog.LevelDebug
	} else {
		logLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}
