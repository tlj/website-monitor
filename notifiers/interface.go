package notifiers

import (
	"website-monitor/result"
)

type Notifier interface {
	Notify(name, displayUrl string, result *result.Results) error
}

type NotifierHolder struct {
	Notifier Notifier
}

func (n *NotifierHolder) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Name    string            `yaml:"name"`
		Type    string            `yaml:"type"`
		Options map[string]string `yaml:"options"`
	}

	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	switch tmp.Type {
	case "slack":
		sl, err := NewSlackNotifier(tmp.Options)
		if err != nil {
			return err
		}
		n.Notifier = sl
	case "pushsafer":
		ps, err := NewPushSaferNotifier(tmp.Options)
		if err != nil {
			return err
		}
		n.Notifier = ps
	}

	return nil
}
