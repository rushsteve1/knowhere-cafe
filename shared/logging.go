package shared

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"gorm.io/gorm/logger"
)

func Dbg[T any](v T, msgs ...string) T {
	_, fn, ln, _ := runtime.Caller(1)
	msg := "output"
	if len(msgs) > 0 {
		msg = msgs[0]
	}
	slog.Debug(msg, "value", v, "file", fn, "line", ln)
	return v
}

// LogOrPanic checks the passed error and causes the program to panic.
// If there was no error then an info log is printed.
//
// This function is sometimes called "die", which is fitting.
func LogOrPanic(msg string, err error, args ...any) {
	if err != nil {
		args = append([]any{"error", err.Error()}, args...)
		slog.Error("FATAL "+msg, args...)
		panic(msg)
	}
	slog.Info(msg, args...)
}

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
