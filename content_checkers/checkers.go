package content_checkers

import (
	"fmt"
	"io"
)

type ContentChecker interface {
	Check(r io.Reader ) (bool, error)
	fmt.Stringer
}
