package notifiers

import (
	"website-monitor/messagequeue"
)

type Notifier interface {
	Name() string
	Notify(name, displayUrl string, result *messagequeue.CrawlResult) error
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
	//case "slack":
		//sl, err := NewSlackNotifier(tmp.Name, tmp.Options)
		//if err != nil {
			//return err
		//}
		//n.Notifier = sl
	case "pushsafer":
		ps, err := NewPushSaferNotifier(tmp.Name, tmp.Options)
		if err != nil {
			return err
		}
		n.Notifier = ps
	case "postgres":
		pg, err := NewPostgresNotifier(tmp.Name, tmp.Options)
		if err != nil {
			return err
		}
		n.Notifier = pg
	}

	return nil
}
