package content_checkers

import (
	"fmt"
	"github.com/go-rod/rod"
	"io"
)

type HtmlRenderSelectorChecker struct {
	name          string
	path          string
	expected      string
	expectedEqual bool
}

func NewHtmlRenderSelectorChecker(name, path, expected string, expectedEqual bool) *HtmlRenderSelectorChecker {
	return &HtmlRenderSelectorChecker{
		name:          name,
		path:          path,
		expected:      expected,
		expectedEqual: expectedEqual,
	}
}

func (h *HtmlRenderSelectorChecker) Check(r io.Reader) (bool, error) {
	panic("implement me")
}

func (h *HtmlRenderSelectorChecker) CheckRender(p *rod.Page) (bool, error) {
	el, err := p.Sleeper(rod.NotFoundSleeper).Element(h.path)
	if err != nil {
		if !h.expectedEqual {
			return true, err
		} else {
			return false, err
		}
	}
	txt, err := el.Text()
	if err != nil {
		if !h.expectedEqual {
			return true, err
		} else {
			return false, err
		}
	}

	if !h.expectedEqual {
		if txt == h.expected {
			return false, nil
		}
	} else {
		if txt != h.expected {
			return false, nil
		}
	}

	return true, nil
}

func (h *HtmlRenderSelectorChecker) String() string {
	if h.expectedEqual {
		return fmt.Sprintf("%s - '%s' is '%s'", h.name, h.path, h.expected)
	} else {
		return fmt.Sprintf("%s - '%s' is not '%s'", h.name, h.path, h.expected)
	}
}

func (h *HtmlRenderSelectorChecker) Type() string {
	return "HtmlRenderSelectorChecker"
}
