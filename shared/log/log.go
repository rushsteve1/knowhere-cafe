// This package implements a wrapper around log/slog providing more features

package log

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
)

var ShowCaller bool = false

func caller(args ...any) []any {
	if ShowCaller {
		// Have to be careful with this wrappers to keep a consistent depth
		_, file, line, ok := runtime.Caller(2)
		if ok {
			args = append(args, "caller", fmt.Sprintf("%s:%d", file, line))
		}
	}
	return args
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	slog.DebugContext(ctx, msg, caller(args...)...)
}

func Debug(msg string, args ...any) {
	slog.Debug(msg, caller(args...)...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, msg, caller(args...)...)
}

func Info(msg string, args ...any) {
	slog.Info(msg, caller(args...)...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, msg, caller(args...)...)
}

func Warn(msg string, args ...any) {
	slog.Warn(msg, caller(args...)...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, caller(args)...)
}

func Error(msg string, args ...any) {
	slog.Error(msg, caller(args)...)
}

func FatalContext(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, msg, caller(args...)...)
	panic(msg)
}

func Fatal(msg string, args ...any) {
	slog.Error(msg, caller(args...)...)
	panic(msg)
}

func Inspect[T any](v T, args ...any) T {
	msg := "output"
	if len(args) > 0 {
		msg = args[0].(string)
		args = args[1:]
	}

	slog.Debug(msg, args...)
	return v
}

func Must[T any](t T, err error, args ...any) T {
	msg := "must not be an error"
	if len(args) > 0 {
		msg = args[0].(string)
		args = args[1:]
	}

	slog.Debug(msg, caller(args...)...)

	if err != nil {
		slog.Error(msg, caller([]any{"error", err}...)...)
		panic(msg)
	}
	return t
}

func Check(err error, args ...any) {
	msg := "checking for an error"
	if len(args) > 0 {
		msg = args[0].(string)
		args = args[1:]
	}

	if err != nil {
		args = append(args, "error", err)
		slog.Error(msg, caller(args...)...)
		panic(msg)
	}

	slog.Debug(msg, caller(args...)...)
}
