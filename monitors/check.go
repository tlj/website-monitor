package monitors

import (
	"fmt"
	"website-monitor/content_checkers"
	"website-monitor/notifiers"

	log "github.com/sirupsen/logrus"
)

type monitorType string

const (
	HttpMonitorType monitorType = "http"
)

type Check struct {
	Name                string                            `yaml:"name"`
	Url                 string                            `yaml:"url"`
	DisplayUrl          string                            `yaml:"display_url"`
	Type                monitorType                       `yaml:"type"`
	Headers             map[string]string                 `yaml:"headers"`
	RegexNotExpected    string                            `yaml:"regex_not_expected"`
	RegexExpected       string                            `yaml:"regex_expected"`
	ExpectedStatusCode  int                               `yaml:"expected_status_code"`
	ContentChecks       []content_checkers.ContentChecker `yaml:"-"`
	LastSeenState       bool                              `yaml:"last_seen_state"`
	ContentChecksConfig []map[string]string               `yaml:"content_checks"`
	Notifiers           []notifiers.Notifier              `yaml:"-"`
	NotifiersConfig     []map[string]string               `yaml:"notifiers"`
	Interval            int                               `yaml:"interval"`
}


func (c *Check) ParseConfig() error {
	if c.RegexExpected != "" {
		c.ContentChecks = append(c.ContentChecks, content_checkers.NewRegexChecker(c.RegexExpected, c.RegexExpected, true))
	}

	if c.RegexNotExpected != "" {
		c.ContentChecks = append(c.ContentChecks, content_checkers.NewRegexChecker(c.RegexNotExpected, c.RegexNotExpected, false))
	}

	for _, cc := range c.ContentChecksConfig {
		var expected string
		var expectedToExist bool
		if v, ok := cc["expected"]; ok {
			expected = v
			expectedToExist = true
		}
		if v, ok := cc["not_expected"]; ok {
			expected = v
			expectedToExist = false
		}

		var contentCheck content_checkers.ContentChecker
		switch cc["type"] {
		case "JsonPath":
			contentCheck = content_checkers.NewJsonPathChecker(cc["name"], cc["path"], expected, expectedToExist)
		case "Regex":
			contentCheck = content_checkers.NewRegexChecker(cc["name"], expected, expectedToExist)
		default:
			return fmt.Errorf("unsupported contentCheck config: %s", cc["type"])
		}

		c.ContentChecks = append(c.ContentChecks, contentCheck)
	}

	for _, n := range c.NotifiersConfig {
		var notifier notifiers.Notifier
		switch n["type"] {
		case "slack":
			notifier = notifiers.NewSlackNotifier(n["webhook"])
		default:
			return fmt.Errorf("unsupported notifiers config: %s", n["type"])
		}
		c.Notifiers = append(c.Notifiers, notifier)
	}

	return nil
}

func (c *Check) Run() error {
	var jm Monitor
	switch c.Type {
	case HttpMonitorType:
		jm = &HttpMonitor{}
	case "":
		jm = &HttpMonitor{}
	default:
		return fmt.Errorf("invalid monitortype '%s'", c.Type)
	}
	if c.Type == HttpMonitorType {
		jm = &HttpMonitor{}
	}
	result, err := jm.Check(*c)
	if err != nil {
		return err
	}
	if result == nil {
		return fmt.Errorf("empty results from Check")
	}

	for _, result := range result.Results {
		log.Debugf("%s: %t (err: %v)", result.ContentChecker, result.Result, result.Err)
	}

	if result.AllTrue() != c.LastSeenState {
		log.Debugf("%s %s: %t", c.Name, c.Url, result.AllTrue())
		log.Infof("State change for %s: %t", c.Name, result)
		c.LastSeenState = result.AllTrue()
		for _, n := range c.Notifiers {
			log.Debugf("Sending notification...")
			err := n.Notify(c.Name, c.DisplayUrl, result)
			if err != nil {
				log.Warn(err)
			}
		}
	}

	return nil
}
