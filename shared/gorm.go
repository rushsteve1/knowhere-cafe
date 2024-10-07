package shared

import (
	"context"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

type GormLogger struct {
	logger *slog.Logger
}

func (GormLogger) LogMode(lvl logger.LogLevel) logger.Interface {
	level := slog.LevelWarn
	switch lvl {
	case logger.Silent:
		level = slog.LevelDebug
	case logger.Error:
		level = slog.LevelError
	case logger.Warn:
		level = slog.LevelWarn
	case logger.Info:
		level = slog.LevelInfo
	}

	return GormLogger{
		logger: slog.New(
			slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}),
		),
	}
}

func (GormLogger) Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args...)
}

func (GormLogger) Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args...)
}

func (GormLogger) Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args...)
}

func (GormLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (sql string, rowsAffected int64),
	err error,
) {
	// TODO IDK?
}
