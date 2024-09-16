// Tools for working with [error]

package shared

import "log/slog"

type ErrUnimplemented struct{}

func (ErrUnimplemented) Error() string {
	return "not implemented"
}

type ErrMissingState struct{}

func (ErrMissingState) Error() string {
	return "missing state key"
}

func Must[T any](t T, err error) T {
	if err != nil {
		slog.Error("must not be an error", "error", err)
		panic(err)
	}
	return t
}
