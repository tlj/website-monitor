package content_checkers

import (
	"io"
)

type ContentChecker interface {
	Check(r io.Reader ) (bool, error)
}
