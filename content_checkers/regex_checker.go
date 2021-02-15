package content_checkers

import (
	"fmt"
	"github.com/go-rod/rod"
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

func (c *RegexChecker) String() string {
	if c.ExpectedExisting {
		return fmt.Sprintf("%s - '%s' found", c.Name, c.Regex)
	} else {
		return fmt.Sprintf("%s - '%s' not found", c.Name, c.Regex)
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
	if !c.ExpectedExisting {
		return !exists, nil
	}

	return exists, nil
}

func (c *RegexChecker) CheckRender(p *rod.Page) (bool, error) {
	panic("implement me")
}
