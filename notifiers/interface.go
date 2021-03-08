package notifiers

import (
	"fmt"
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
		if webhook, ok := tmp.Options["webhook"]; ok {
			n.Notifier = NewSlackNotifier(webhook)
		} else {
			return fmt.Errorf("slack option 'webhook' is required")
		}
	case "pushsafer":
		if privateKey, ok := tmp.Options["private_key"]; ok {
			n.Notifier = NewPushSaferNotifier(privateKey)
		} else {
			return fmt.Errorf("pushsafer option 'private_key' is required")
		}
	}

	return nil
}
