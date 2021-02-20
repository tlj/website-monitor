package content_checkers

import (
	"fmt"
	"github.com/go-rod/rod"
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

func (c *RegexChecker) String() string {
	if c.expectedExisting {
		return fmt.Sprintf("%s - '%s' found", c.name, c.regex)
	} else {
		return fmt.Sprintf("%s - '%s' not found", c.name, c.regex)
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
	if !c.expectedExisting {
		return !exists, nil
	}

	return exists, nil
}

func (c *RegexChecker) CheckRender(p *rod.Page) (bool, error) {
	panic("implement me")
}

func (c *RegexChecker) Type() string {
	return "RegexChecker"
}

func (c *RegexChecker) Equal(y *RegexChecker) bool {
	return c.name == y.name && c.regex == y.regex && c.expectedExisting == y.expectedExisting
}
