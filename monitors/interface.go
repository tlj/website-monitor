package monitors

import "website-monitor/result"

type Monitor interface {
	Check(check Check) (*result.Results, error)
	Type() string
}
