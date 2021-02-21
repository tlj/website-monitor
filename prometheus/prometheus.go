package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	MonitorsProcessedTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "monitors_processed_total",
		Help: "The total number of processed monitors",
	})
	MonitorsErroredTotal = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "monitors_errored_total",
		Help: "The total number of errored monitors",
	})
	JobQueueGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "monitors_queued",
		Help: "Number of monitors in queue",
	})
	LastSeenState = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "monitors_state",
		Help: "The current state of the monitor.",
	},
		[]string{"monitor"})
	MonitorsIndividualProcessed = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "monitors_individual_processed_total",
		Help: "The total number of individually processed monitors.",
	},
		[]string{"monitor"})
	MonitorsIndividualErrored = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "monitors_individual_errored_total",
		Help: "The total number of individually errored monitors.",
	},
		[]string{"monitor"})
	MonitorsNextCheckInfo = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "monitors_next_check_at",
		Help: "Unix timestamp for next check.",
	},
		[]string{"monitor"})
)

func Init() {
	prometheus.MustRegister(
		MonitorsErroredTotal,
		MonitorsProcessedTotal,
		JobQueueGauge,
		LastSeenState,
		MonitorsIndividualProcessed,
		MonitorsIndividualErrored,
		MonitorsNextCheckInfo,
	)
}
