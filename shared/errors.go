// Tools for working with [error]

package shared

type ErrUnimplemented struct{}

func (ErrUnimplemented) Error() string {
	return "not implemented"
}

type ErrMissingState struct{}

func (ErrMissingState) Error() string {
	return "missing state key"
}
