package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
	"website-monitor/app"
	"website-monitor/monitors"
	"website-monitor/prometheus"
)

func main() {
	prometheus.Init()

	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.SetLevel(log.InfoLevel)

	config := &app.Config{}
	err := config.LoadConfigFromFile("config/config.yaml")
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
	for _, m := range config.Monitors {
		prometheus.LastSeenState.WithLabelValues(m.Name).Set(0)
		prometheus.MonitorsIndividualProcessed.WithLabelValues(m.Name).Add(0)
		prometheus.MonitorsIndividualErrored.WithLabelValues(m.Name).Add(0)
		prometheus.MonitorsNextCheckInfo.WithLabelValues(m.Name).Set(float64(m.NextCheckAt().Unix()))
	}

	log.Infof("Starting %d checks...", len(checks))

	queue := make(chan *monitors.Monitor, 10)

	go func() {
		t := time.NewTimer(1 * time.Second)
		for {
			select {
			case <-t.C:
				log.Debug("Looking for job...")
				for _, c := range checks {
					if c.ShouldUpdate() {
						log.Infof("Queuing %s...", c.Name)
						c.CheckPending = true
						queue <- c
						prometheus.JobQueueGauge.Inc()
					}
				}
				t.Reset(1 * time.Second)
			}
		}
	}()

	go func() {
		for {
			select {
			case c := <-queue:
				prometheus.JobQueueGauge.Dec()
				go func(ch *monitors.Monitor) {
					prometheus.MonitorsProcessedTotal.Inc()
					prometheus.MonitorsIndividualProcessed.WithLabelValues(ch.Name).Inc()
					if err := ch.Run(); err != nil {
						prometheus.MonitorsErroredTotal.Inc()
						log.Errorf("Error in %s: %v", ch.Name, err)
					}
					if ch.LastSeenState {
						prometheus.LastSeenState.WithLabelValues(ch.Name).Set(1)
					} else {
						prometheus.LastSeenState.WithLabelValues(ch.Name).Set(0)
					}
					prometheus.MonitorsNextCheckInfo.WithLabelValues(ch.Name).Set(float64(ch.NextCheckAt().Unix()))
				}(c)
			}
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(":2112", nil))

}
