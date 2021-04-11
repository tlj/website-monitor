package monitors

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"time"
	"website-monitor/content_checkers"
	"website-monitor/result"
	"website-monitor/schedule"
)

type MonitorType string

const (
	HttpMonitorType       MonitorType = "http"
	HttpRenderMonitorType MonitorType = "http_render"
)

type Monitor struct {
	tableName struct{} `pg:"checks,alias:check"`

	// MonitorInterface
	Id     int64 `yaml:"-"`
	UserId int64 `yaml:"-"`

	Name               string            `yaml:"name"`
	Url                string            `yaml:"url"`
	DisplayUrl         string            `yaml:"display_url"`
	Type               MonitorType       `yaml:"type"`
	Headers            map[string]string `yaml:"headers" pg:"-"`
	ExpectedStatusCode int               `yaml:"expected_status_code"`

	// Schedule
	Scheduler *schedule.Schedule `yaml:"schedule" pg:"-"`

	// Status
	lastCheckedAt time.Time `pg:"-" yaml:"-"`
	nextCheckAt   time.Time `pg:"-" yaml:"-"`
	CheckPending  bool      `pg:"-" yaml:"-"`
	LastSeenState bool      `pg:"-" yaml:"-"`

	// Config
	RenderServerURN string                            `yaml:"render_server_urn" pg:"-"`
	ContentChecks   []content_checkers.ContentChecker `yaml:"checks" pg:"-"`
	RequireSome     bool                              `yaml:"require_some" pg:"-"`
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
		// when doing a separate schedule we don't know about when the update was done, so
		// let's just update this now
		c.updateTimestamps()

		return true
	}

	return false
}

func (c *Monitor) Run() (*result.Results, error) {
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
		return nil, fmt.Errorf("invalid monitortype '%s'", c.Type)
	}
	if c.Type == HttpMonitorType {
		jm = &HttpMonitor{}
	}
	result, err := jm.Check(*c)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("empty results from Monitor")
	}

	for _, res := range result.Results {
		log.Debugf("%s: %t (err: %v)", res.ContentChecker, res.Result, res.Err)
	}

	return result, nil
}

func (c *Monitor) GetId() int64 {
	return c.Id
}

func (c *Monitor) GetName() string {
	return c.Name
}
