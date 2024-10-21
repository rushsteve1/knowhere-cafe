// Tools for working with [error]

package shared

import (
	"errors"
	"fmt"
)

var ErrUnimplemented = errors.New("not implemented")
var ErrMissingState = errors.New("missing state key")
var ErrNotAuth = errors.New("not authorized over tailscale")

type ErrUnknownTemplate struct{ Name string }

func (e ErrUnknownTemplate) Error() string {
	return fmt.Sprintf("unknown template '%s'", e.Name)
}
