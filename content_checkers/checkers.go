package content_checkers

import (
	"fmt"
	"github.com/go-rod/rod"
	"io"
)

type ContentChecker interface {
	Check(r io.Reader) (bool, error)
	CheckRender(p *rod.Page) (bool, error)
	fmt.Stringer
}
