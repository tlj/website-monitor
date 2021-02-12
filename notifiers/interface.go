package notifiers

import (
	"website-monitor/result"
)

type Notifier interface {
	Notify(name, displayUrl string, result *result.Results) error
}
