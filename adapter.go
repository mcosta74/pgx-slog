package log

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/tracelog"
)

type Logger struct {
	l *slog.Logger
	invalidLevelKey string
}

type Option func(*Logger)

func WithInvalidLevelKey(key string) Option {
	return func(l *Logger) {
		l.invalidLevelKey = key
	}
}

func NewLogger(l *slog.Logger, options ...Option) *Logger {
	logger := &Logger{
		l: l,
		invalidLevelKey: "INVALID_PGX_LOG_LEVEL",
	}

	for _, option := range options {
		option(logger)
	}

	return logger
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	attrs := make([]slog.Attr, 0, len(data))
	for k, v := range data {
		attrs = append(attrs, slog.Any(k, v))
	}

	var lvl slog.Level
	switch level {
	case tracelog.LogLevelTrace:
		lvl = slog.LevelDebug - 1
		attrs = append(attrs, slog.Any("PGX_LOG_LEVEL", level))
	case tracelog.LogLevelDebug:
		lvl = slog.LevelDebug
	case tracelog.LogLevelInfo:
		lvl = slog.LevelInfo
	case tracelog.LogLevelWarn:
		lvl = slog.LevelWarn
	case tracelog.LogLevelError:
		lvl = slog.LevelError
	default:
		lvl = slog.LevelError
		attrs = append(attrs, slog.Any(l.invalidLevelKey, fmt.Errorf("invalid pgx log level: %v", level)))
	}
	l.l.LogAttrs(ctx, lvl, msg, attrs...)
}
