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

type Monitor struct {
	tableName struct{} `pg:"checks,alias:check"`

	// MonitorInterface
	ID                 int               `yaml:"-"`
	Name               string            `yaml:"name"`
	Url                string            `yaml:"url"`
	DisplayUrl         string            `yaml:"display_url"`
	Type               MonitorType       `yaml:"type"`
	Headers            map[string]string `yaml:"headers" pg:"-"`
	ExpectedStatusCode int               `yaml:"expected_status_code"`

	// Schedule
	Scheduler *scheduler.Scheduler `yaml:"schedule" pg:"-"`

	// Status
	lastCheckedAt time.Time `pg:"-" yaml:"-"`
	nextCheckAt   time.Time `pg:"-" yaml:"-"`
	CheckPending  bool      `pg:"-" yaml:"-"`
	LastSeenState bool      `pg:"-" yaml:"-"`

	// Config
	RenderServerURN string                                  `yaml:"render_server_urn" pg:"-"`
	ContentChecks   []content_checkers.ContentCheckerHolder `yaml:"checks" pg:"-"`
	RequireSome     bool                                    `yaml:"require_some" pg:"-"`

	// Notifiers
	Notifiers []notifiers.NotifierHolder `yaml:"notifiers" pg:"-"`
}

func (c *Monitor) updateTimestamps() {
	c.lastCheckedAt = time.Now().UTC()
	c.nextCheckAt = c.Scheduler.CalculateNextFrom(c.lastCheckedAt)
	log.Debugf("%s next run: %s (in %ds)", c.Name, c.nextCheckAt.String(), int(c.nextCheckAt.Sub(time.Now()).Seconds()))

	c.CheckPending = false
}

func (c *Monitor) NextCheckAt() time.Time {
	return c.nextCheckAt
}

func (c *Monitor) ShouldUpdate() bool {
	// This is probably most correct, but leads to no immediate checks. Should
	// we do immediate checks?
	//if c.nextCheckAt.IsZero() {
	//	c.updateTimestamps()
	//}

	if !c.CheckPending && c.nextCheckAt.Sub(time.Now()) <= 0 {
		return true
	}

	return false
}

func (c *Monitor) Run() error {
	defer c.updateTimestamps()

	var jm MonitorInterface
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
		return fmt.Errorf("empty results from Monitor")
	}

	for _, result := range result.Results {
		log.Debugf("%s: %t (err: %v)", result.ContentChecker, result.Result, result.Err)
	}

	var endResult bool
	switch c.RequireSome {
	case true:
		endResult = result.SomeTrue()
	case false:
		endResult = result.AllTrue()
	}

	if endResult != c.LastSeenState {
		log.Debugf("%s %s: %t", c.Name, c.Url, endResult)
		log.Infof("State change for %s: %t", c.Name, result)
		c.LastSeenState = endResult
		for _, n := range c.Notifiers {
			log.Debugf("Sending notification to '%s'...", n.Notifier.Name())
			err := n.Notifier.Notify(c.Name, c.DisplayUrl, result)
			if err != nil {
				log.Warn(err)
			}
		}
	}

	return nil
}
