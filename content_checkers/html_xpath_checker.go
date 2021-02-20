package content_checkers

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/go-rod/rod"
	"io"
)

type HtmlXPathChecker struct {
	name          string
	path          string
	expected      string
	expectedEqual bool
}

func NewHtmlXPathChecker(name, path, expected string, expectedEqual bool) *HtmlXPathChecker {
	return &HtmlXPathChecker{
		name:          name,
		path:          path,
		expected:      expected,
		expectedEqual: expectedEqual,
	}
}

func (j *HtmlXPathChecker) String() string {
	if j.expectedEqual {
		return fmt.Sprintf("%s - '%s' is '%s'", j.name, j.path, j.expected)
	} else {
		return fmt.Sprintf("%s - '%s' is not '%s'", j.name, j.path, j.expected)
	}
}

func (j *HtmlXPathChecker) Check(r io.Reader) (bool, error) {
	doc, err := htmlquery.Parse(r)
	if err != nil {
		return false, err
	}

	stockMeta := htmlquery.FindOne(doc, j.path)
	if stockMeta == nil {
		if j.expectedEqual {
			return false, nil
		} else {
			return true, nil
		}
	}

	if !j.expectedEqual {
		if stockMeta.Data == j.expected {
			return false, nil
		}
	} else {
		if stockMeta.Data != j.expected {
			return false, nil
		}
	}

	return true, nil
}

func (j *HtmlXPathChecker) CheckRender(p *rod.Page) (bool, error) {
	panic("implement me")
}

func (j *HtmlXPathChecker) Type() string {
	return "HtmlXPathChecker"
}

func (j *HtmlXPathChecker) Equal(y *HtmlXPathChecker) bool {
	return j.name == y.name && j.path == y.path && j.expected == y.expected && j.expectedEqual == y.expectedEqual
}