package log

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/tracelog"
	"golang.org/x/exp/slog"
)

type Logger struct {
	l *slog.Logger
}

func NewLogger(l *slog.Logger) *Logger {
	return &Logger{l: l}
}

func (l *Logger) Log(ctx context.Context, level tracelog.LogLevel, msg string, data map[string]interface{}) {
	logger := l.l
	for k, v := range data {
		logger = logger.With(k, v)
	}

	switch level {
	case pgx.LogLevelTrace:
		logger.Log(context.Background(), slog.LevelDebug-1, msg, "PGX_LOG_LEVEL", level)
	case pgx.LogLevelDebug:
		logger.Debug(msg)
	case pgx.LogLevelInfo:
		logger.Info(msg)
	case pgx.LogLevelWarn:
		logger.Warn(msg)
	case pgx.LogLevelError:
		logger.Error(msg)
	default:
		logger.Error(msg, "INVALID_PGX_LOG_LEVEL", level)
	}
}
