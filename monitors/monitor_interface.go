package monitors

import "website-monitor/result"

type MonitorInterface interface {
	Check(check Monitor) (*result.Results, error)
	Type() string
}
