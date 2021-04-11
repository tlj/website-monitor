package scheduler

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"website-monitor/interfaces"
	"website-monitor/messagequeue"
	"website-monitor/prometheus"
)

type ScheduleEmitterInterface interface {
	HandleMonitor(
		mon interfaces.ShouldUpdaterIdName,
		publisher messagequeue.MessageQueuePublisher,
		exchange string,
		routingKey string)
}

type ScheduleEmitter struct{}

func NewScheduleEmitter() *ScheduleEmitter {
	return &ScheduleEmitter{}
}

func (s *ScheduleEmitter) HandleMonitor(
	mon interfaces.ShouldUpdaterIdName,
	publisher messagequeue.MessageQueuePublisher,
	exchange string,
	routingKey string) {
	if mon.ShouldUpdate() {
		log.Infof("Queuing %s...", mon.GetName())

		msg, err := json.Marshal(messagequeue.ScheduleJob{
			MonitorID: mon.GetId(),
			Name:      mon.GetName(),
			Type:      "monitor",
		})
		if err != nil {
			log.Warnf("Error marshalling message for %s", mon.GetName())
			return
		}

		err = publisher.Publish(
			exchange,   // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        msg,
			})
		if err != nil {
			log.Warnf("Error publishing message '%s': %s", msg, err)
		}

		prometheus.JobQueueGauge.Inc()
	}
}
