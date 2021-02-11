package content_checkers

import (
	"io"
	"io/ioutil"
	"regexp"
)

type RegexChecker struct {
	Name             string
	Regex            string
	ExpectedExisting bool
}

func NewRegexChecker(name, regex string, expectedExisting bool) *RegexChecker {
	return &RegexChecker{
		Name:             name,
		Regex:            regex,
		ExpectedExisting: expectedExisting,
	}
}

func (c *RegexChecker) Check(r io.Reader) (bool, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return false, err
	}

	rx, err := regexp.Compile(c.Regex)
	if err != nil {
		return false, err
	}

	exists := rx.Match(data)
	return exists && c.ExpectedExisting, nil
}
