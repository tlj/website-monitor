package content_checkers

import (
	"io"
	"io/ioutil"
	"regexp"
)

type RegexChecker struct {
	name             string
	regex            string
	expectedExisting bool
}

func NewRegexChecker(name, regex string, expectedExisting bool) *RegexChecker {
	return &RegexChecker{
		name:             name,
		regex:            regex,
		expectedExisting: expectedExisting,
	}
}

func (c *RegexChecker) Check(r io.Reader) (bool, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return false, err
	}

	rx, err := regexp.Compile(c.regex)
	if err != nil {
		return false, err
	}

	exists := rx.Match(data)
	return exists && c.expectedExisting, nil
}
