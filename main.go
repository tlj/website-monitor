package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/url"
	"time"
	"website-monitor/monitors"
)

type Config struct {
	Global struct {
		Headers            map[string]string   `yaml:"headers"`
		ExpectedStatusCode int                 `yaml:"expected_status_code"`
		Interval           int                 `yaml:"interval"`
		NotifiersConfig    []map[string]string `yaml:"notifiers"`
	} `yaml:"global"`
	Monitors []monitors.Check `yaml:"monitors"`
}

func schedule(what func() error, delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			err := what()
			if err != nil {
				log.Println(err)
			}
			select {
			case <-time.After(delay):
			case <-stop:
				return
			}
		}
	}()

	return stop
}

func main() {
	var config Config
	configData, err := ioutil.ReadFile("config/config.yaml")
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(configData, &config)
	if err != nil {
		panic(err)
	}

	for k, monitor := range config.Monitors {
		if monitor.Headers == nil {
			config.Monitors[k].Headers = make(map[string]string)
		}
		if monitor.Interval == 0 {
			config.Monitors[k].Interval = config.Global.Interval
		}
		if monitor.ExpectedStatusCode == 0 {
			config.Monitors[k].ExpectedStatusCode = config.Global.ExpectedStatusCode
		}
		for gk, gv := range config.Global.Headers {
			if _, ok := monitor.Headers[gk]; !ok {
				config.Monitors[k].Headers[gk] = gv
			}
		}
		for _, gm := range config.Global.NotifiersConfig {
			config.Monitors[k].NotifiersConfig = append(config.Monitors[k].NotifiersConfig, gm)
		}
		if monitor.DisplayUrl == "" {
			config.Monitors[k].DisplayUrl = monitor.Url
		}
		if _, ok := monitor.Headers["Referer"]; !ok {
			u, _ := url.Parse(monitor.Url)
			config.Monitors[k].Headers["Referer"] = fmt.Sprintf("%s://%s/", u.Scheme, u.Host)
		}
		_ = config.Monitors[k].ParseConfig()
	}

	checks := config.Monitors
	log.Printf("Starting %d checks...", len(checks))

	for _, c := range checks {
		log.Printf("Starting %s interval %ds", c.Name, c.Interval)
		go func(che monitors.Check) {
			schedule(che.Run, time.Duration(che.Interval)*time.Second)
		}(c)
	}

	for {}
}
