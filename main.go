package main

import (
	log "github.com/sirupsen/logrus"
	"time"
	"website-monitor/app"
	"website-monitor/monitors"
)

func schedule(what func() error, delay time.Duration) chan bool {
	stop := make(chan bool)

	go func() {
		for {
			err := what()
			if err != nil {
				log.Error(err)
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
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)

	config, err := app.LoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error while loading config/config.yaml: %s", err)
	}

	switch config.LogLevel {
	case "debug":
		log.SetLevel(log.DebugLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "error":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	checks := config.Monitors
	log.Infof("Starting %d checks...", len(checks))

	for _, c := range checks {
		log.Infof("Starting %s interval %ds", c.Name, c.Interval)
		go func(che monitors.Check) {
			schedule(che.Run, time.Duration(che.Interval)*time.Second)
		}(c)
	}

	for { select {} }
}
