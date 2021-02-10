package monitors

import (
	"fmt"
	"log"
	"website-monitor/content_checkers"
	"website-monitor/notifiers"
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
	ExpectedStatusCode  int                               `yaml:"expected_status_code"`
	ContentChecks       []content_checkers.ContentChecker `yaml:"-"`
	LastSeenState       bool                              `yaml:"last_seen_state"`
	ContentChecksConfig []map[string]string               `yaml:"content_checks"`
	Notifiers           []notifiers.Notifier              `yaml:"-"`
	NotifiersConfig     []map[string]string               `yaml:"notifiers"`
	Interval            int                               `yaml:"interval"`
}

func (c *Check) ParseConfig() error {
	for _, cc := range c.ContentChecksConfig {
		var contentCheck content_checkers.ContentChecker
		switch cc["type"] {
		case "JsonPath":
			contentCheck = content_checkers.NewJsonPathChecker(cc["name"], cc["path"], cc["expected"], cc["expected_to_exists"] == "true")
		case "Regex":
			contentCheck = content_checkers.NewRegexChecker(cc["name"], cc["regex"], cc["expected_to_exist"] == "true")
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
	default:
		return fmt.Errorf("invalid monitortype '%s'", c.Type)
	}
	if c.Type == HttpMonitorType {
		jm = &HttpMonitor{}
	}
	result, err := jm.Check(*c)
	if err != nil {
		log.Println(err)
	}
	log.Printf("%s %s: %t", c.Name, c.Url, result)

	if result != c.LastSeenState {
		c.LastSeenState = result
		for _, n := range c.Notifiers {
			log.Println("Sending notification...")
			err := n.Notify(fmt.Sprintf("%s new state: %t (%s)", c.Name, result, c.DisplayUrl))
			if err != nil {
				log.Println(err)
			}
		}
	}

	return nil
}
