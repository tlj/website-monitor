package main

import (
	log "github.com/sirupsen/logrus"
	"time"
	"website-monitor/app"
	"website-monitor/monitors"
)

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

	queue := make(chan *monitors.Check, 10)

	go func() {
		t := time.NewTimer(1 * time.Second)
		for {
			select {
				case <- t.C:
					log.Println("Looking for job...")
					for _, c := range checks {
						if c.ShouldUpdate() {
							log.Printf("Should check %s...", c.Name)
							c.CheckPending = true
							queue <- c
						}
					}
					t.Reset(1 * time.Second)
			}
		}
	}()

	for {
		select {
		case c := <- queue:
			go func(ch *monitors.Check) {
				if err := ch.Run(); err != nil {
					log.Println(err)
				}
			}(c)
		}
	}
}
