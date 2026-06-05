package utilities

import (
	"log/slog"
	"os"
)

var Logger *slog.Logger

func NewLogger(_debug, _verbose bool) {
	var logLevel slog.Level

	if _debug {
		logLevel = slog.LevelDebug
	} else if _verbose {
		logLevel = slog.LevelInfo
	} else {
		logLevel = slog.LevelError + 1 //dubway of saying off
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	Logger = slog.New(handler)
	slog.SetDefault(Logger)
}
