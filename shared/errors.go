// Tools for working with [error]

package shared

import (
	"context"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

type UnimplementedError struct{}

func (UnimplementedError) Error() string {
	return "Not Implemented"
}

func LogOrDie(msg string, err error, args ...any) {
	if err != nil {
		args = append([]any{"error", err}, args...)
		slog.Error("DIE: "+msg, args...)
		os.Exit(1)
	}
	slog.Info(msg, args...)
}

type GormLogger struct{}

func (GormLogger) LogMode(logger.LogLevel) logger.Interface {
	// TODO IDK?
	return nil
}

func (GormLogger) Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, args)
}

func (GormLogger) Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, args)
}

func (GormLogger) Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, args)
}

func (GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	// TODO IDK?
}
