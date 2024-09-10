// Tools for working with [error]

package shared

type UnimplementedError struct{}

func (UnimplementedError) Error() string {
	return "Not Implemented"
}

type NilValue struct{}

func (NilValue) Error() string {
	return "nil value"
}

type Maybe[T any] struct {
	v *T
}

func (m Maybe[T]) Unwrap() (T, error) {
	if m.v != nil {
		return *m.v, nil
	}
	var v T
	return v, NilValue{}
}

func Some[T any](v T) Maybe[T] {
	return Maybe[T]{&v}
}

func None[T any]() Maybe[T] {
	return Maybe[T]{nil}
}
