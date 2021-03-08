package content_checkers

import (
	"fmt"
	"github.com/go-rod/rod"
	"io"
)

type CheckType string

const (
	RegexCheckType CheckType = "regex"
	HtmlXpathType  CheckType = "html_xpath"
	JsonPathType   CheckType = "json_path"
	HtmlRenderType CheckType = "html_render"
)

type ContentChecker interface {
	Check(r io.Reader) (bool, error)
	CheckRender(p *rod.Page) (bool, error)
	Type() string
	fmt.Stringer
}

type ContentCheckerHolder struct {
	ContentChecker ContentChecker
}

func (cch *ContentCheckerHolder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Name       string    `yaml:"name"`
		CheckType  CheckType `yaml:"type"`
		Path       string    `yaml:"path"`
		Value      string    `yaml:"value"`
		IsExpected bool      `yaml:"is_expected"`
	}

	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	switch tmp.CheckType {
	case RegexCheckType:
		cch.ContentChecker = NewRegexChecker(tmp.Name, tmp.Value, tmp.IsExpected)
	case HtmlXpathType:
		cch.ContentChecker = NewHtmlXPathChecker(tmp.Name, tmp.Path, tmp.Value, tmp.IsExpected)
	case JsonPathType:
		cch.ContentChecker = NewJsonPathChecker(tmp.Name, tmp.Path, tmp.Value, tmp.IsExpected)
	case HtmlRenderType:
		cch.ContentChecker = NewHtmlRenderSelectorChecker(tmp.Name, tmp.Path, tmp.Value, tmp.IsExpected)
	default:
		return fmt.Errorf("unsupported contentCheck config: '%s'", tmp.CheckType)

	}

	return nil
}
