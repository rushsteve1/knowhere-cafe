// Functions that make my life easier

package easy

import (
	"fmt"
	"log/slog"
	"runtime"
)

func Ternary[T any](t bool, a, b T) T {
	if t {
		return a
	}
	return b
}

func Inspect[T any](v T, args ...any) T {
	msg, args := PopOr(args, "output")

	slog.Debug(msg.(string), args...)
	return v
}

func Must[T any](t T, err error, args ...any) T {
	msg, args := PopOr(args, "must not be an error")

	_, file, line, _ := runtime.Caller(1)
	args = append(args, "caller", fmt.Sprintf("%s:%d", file, line))

	if err != nil {
		args = append(args, []any{"error", err})
		slog.Error(msg.(string), args...)
		panic(msg)
	}

	slog.Debug(msg.(string), args...)
	return t
}

func Expect(err error, args ...any) {
	msg, args := PopOr(args, "checking for an error")

	_, file, line, _ := runtime.Caller(1)
	args = append(args, "caller", fmt.Sprintf("%s:%d", file, line))

	if err != nil {
		args = append(args, "error", err)
		slog.Error(msg.(string), args...)
		panic(err)
	}

	slog.Debug(msg.(string), args...)
}

func Assert(t bool, args ...any) {
	if t {
		return
	}

	msg, args := PopOr(args, "assertion failure")

	_, file, line, _ := runtime.Caller(1)
	args = append(args, "caller", fmt.Sprintf("%s:%d", file, line))

	slog.Error(msg.(string), args...)
	panic(msg)
}

func AssertEq[T comparable](a, b T, args ...any) {
	if a != b {
		return
	}

	msg, args := PopOr(args, "assertion failure")

	args = append(args, "left", a, "right", b)

	_, file, line, _ := runtime.Caller(1)
	args = append(args, "caller", fmt.Sprintf("%s:%d", file, line))

	slog.Error(msg.(string), args...)
	panic(msg)
}
