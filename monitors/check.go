package monitors

import (
	"fmt"
	"time"
	"website-monitor/content_checkers"
	"website-monitor/notifiers"
	"website-monitor/scheduler"

	log "github.com/sirupsen/logrus"
)

type MonitorType string

const (
	HttpMonitorType       MonitorType = "http"
	HttpRenderMonitorType MonitorType = "http_render"
)

type Check struct {
	Name               string
	Url                string
	DisplayUrl         string
	RenderServerURN    string
	Type               MonitorType
	Headers            map[string]string
	ExpectedStatusCode int
	ContentChecks      []content_checkers.ContentChecker
	Notifiers          []notifiers.Notifier
	Scheduler          *scheduler.Scheduler
	LastSeenState      bool
	lastCheckedAt      time.Time
	nextCheckAt        time.Time
	CheckPending       bool
}

/*
func (c *Check) ParseConfig() error {


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
		case "HtmlXPath":
			contentCheck = content_checkers.NewHtmlXPathChecker(cc["name"], cc["path"], expected, expectedToExist)
		case "HtmlRenderSelector":
			contentCheck = content_checkers.NewHtmlRenderSelectorChecker(cc["name"], cc["path"], expected, expectedToExist)
		default:
			return fmt.Errorf("unsupported contentCheck config: %s", cc["type"])
		}

		c.ContentChecks = append(c.ContentChecks, contentCheck)
	}


	return nil
}

 */


func (c *Check) updateTimestamps() {
	c.lastCheckedAt = time.Now().UTC()
	c.nextCheckAt = c.Scheduler.CalculateNextFrom(c.lastCheckedAt)
	log.Debugf("%s next run: %s (in %ds)", c.Name, c.nextCheckAt.String(), int(c.nextCheckAt.Sub(time.Now()).Seconds()))

	c.CheckPending = false
}

func (c *Check) ShouldUpdate() bool {
	if !c.CheckPending && c.nextCheckAt.Sub(time.Now()) <= 0 {
		return true
	}

	return false
}

func (c *Check) Run() error {
	defer c.updateTimestamps()

	var jm Monitor
	switch c.Type {
	case HttpMonitorType:
		jm = &HttpMonitor{}
	case HttpRenderMonitorType:
		if c.RenderServerURN == "" {
			log.Fatal("Config key 'render_server_urn' is missing or empty, required for http_render type monitors.")
		}
		jm = NewHttpRenderMonitor(c.RenderServerURN)
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
