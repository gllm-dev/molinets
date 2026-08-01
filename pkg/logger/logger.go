package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

func Configure(w io.Writer, level string, json bool) {
	opts := &slog.HandlerOptions{Level: parseLevel(level)}
	var h slog.Handler
	if json {
		h = slog.NewJSONHandler(w, opts)
	} else {
		h = slog.NewTextHandler(w, opts)
	}
	slog.SetDefault(slog.New(h))
}

func ConfigureDefault() { Configure(os.Stdout, "info", true) }

func New(component string) *slog.Logger {
	return slog.Default().With(slog.String("component", component))
}

func Error(err error) slog.Attr { return slog.Any("error", err) }

func parseLevel(s string) slog.Level {
	switch strings.ToLower(s) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
