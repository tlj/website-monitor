package app

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/url"
	"website-monitor/monitors"
)

type Config struct {
	LogLevel string `yaml:"loglevel"`
	Global   struct {
		Headers            map[string]string   `yaml:"headers"`
		ExpectedStatusCode int                 `yaml:"expected_status_code"`
		Interval           int                 `yaml:"interval"`
		NotifiersConfig    []map[string]string `yaml:"notifiers"`
	} `yaml:"global"`
	Monitors []monitors.Check `yaml:"monitors"`
}

func LoadConfig(filename string) (*Config, error) {
	config := &Config{}

	configData, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	err = config.Parse(configData)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Parse(data []byte) error {
	err := yaml.Unmarshal(data, &c)
	if err != nil {
		return err
	}

	c.ExpandMonitors()

	return nil
}

func (c *Config) ExpandMonitors() {
	for k, monitor := range c.Monitors {
		if monitor.Headers == nil {
			c.Monitors[k].Headers = make(map[string]string)
		}
		if monitor.Interval == 0 {
			c.Monitors[k].Interval = c.Global.Interval
		}
		if monitor.ExpectedStatusCode == 0 {
			c.Monitors[k].ExpectedStatusCode = c.Global.ExpectedStatusCode
		}
		for gk, gv := range c.Global.Headers {
			if _, ok := monitor.Headers[gk]; !ok {
				c.Monitors[k].Headers[gk] = gv
			}
		}
		for _, gm := range c.Global.NotifiersConfig {
			c.Monitors[k].NotifiersConfig = append(c.Monitors[k].NotifiersConfig, gm)
		}
		if monitor.DisplayUrl == "" {
			c.Monitors[k].DisplayUrl = monitor.Url
		}
		if _, ok := monitor.Headers["Referer"]; !ok {
			if monitor.DisplayUrl != "" && monitor.Url != "" {
				c.Monitors[k].Headers["Referer"] = monitor.DisplayUrl
			} else {
				u, _ := url.Parse(monitor.Url)
				c.Monitors[k].Headers["Referer"] = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
			}
		}
		_ = c.Monitors[k].ParseConfig()
	}
}