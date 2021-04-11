package content_checkers

import (
	"fmt"
	"github.com/antchfx/jsonquery"
	"github.com/go-rod/rod"
	"io"
)

type JsonPathChecker struct {
	name          string
	path          string
	expected      string
	expectedEqual bool
}

func NewJsonPathChecker(name, path, expected string, expectedEqual bool) *JsonPathChecker {
	return &JsonPathChecker{
		name:          name,
		path:          path,
		expected:      expected,
		expectedEqual: expectedEqual,
	}
}

func (j *JsonPathChecker) String() string {
	if j.expectedEqual {
		return fmt.Sprintf("%s - '%s' is '%s'", j.name, j.path, j.expected)
	} else {
		return fmt.Sprintf("%s - '%s' is not '%s'", j.name, j.path, j.expected)
	}
}

func (j *JsonPathChecker) Check(r io.Reader) (bool, error) {
	doc, err := jsonquery.Parse(r)
	if err != nil {
		return false, err
	}

	stockMeta := jsonquery.FindOne(doc, j.path)
	if stockMeta == nil {
		if j.expectedEqual {
			return false, nil
		} else {
			return true, nil
		}
	}

	if !j.expectedEqual {
		if stockMeta.InnerText() == j.expected {
			return false, nil
		}
	} else {
		if stockMeta.InnerText() != j.expected {
			return false, nil
		}
	}

	return true, nil
}

func (j *JsonPathChecker) CheckRender(p *rod.Page) (bool, error) {
	panic("implement me")
}

func (j *JsonPathChecker) Type() string {
	return "JsonPathChecker"
}

func (j *JsonPathChecker) Equal(y *JsonPathChecker) bool {
	return j.name == y.name && j.path == y.path && j.expected == y.expected && j.expectedEqual == y.expectedEqual
}

func (j *JsonPathChecker) Name() string {
	return j.name
}