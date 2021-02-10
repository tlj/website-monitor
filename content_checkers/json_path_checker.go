package content_checkers

import (
	"github.com/antchfx/jsonquery"
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
