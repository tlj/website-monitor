package scheduler

import (
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"time"
	"website-monitor/messagequeue"
	"website-monitor/monitors"
)

type Scheduler struct {
	messageQueue    messagequeue.MessageQueuePublisher
	monitorRepo     monitors.RepositoryInterface
	scheduleEmitter ScheduleEmitterInterface
	exchange        string
	routingKey      string
}

func NewScheduler(
	messageQueue messagequeue.MessageQueuePublisher,
	monitorRepo monitors.RepositoryInterface,
	scheduleEmitter ScheduleEmitterInterface,
	exchange string,
	routingKey string,
) *Scheduler {
	return &Scheduler{
		messageQueue:    messageQueue,
		monitorRepo:     monitorRepo,
		scheduleEmitter: scheduleEmitter,
		exchange:        exchange,
		routingKey:      routingKey,
	}
}

func (s *Scheduler) Run(cl clock.Clock, interval time.Duration, stop chan bool) {
	t := cl.Timer(interval)

	mons, err := s.monitorRepo.All()
	if err != nil {
		log.Fatal(err)
	}
	for {
		select {
		case <-t.C:
			log.Debug("Looking for job...")
			for _, mon := range mons {
				log.Debugf("Checking %s", mon.Name)
				s.scheduleEmitter.HandleMonitor(mon, s.messageQueue, s.exchange, s.routingKey)
			}
			t.Reset(interval)
		case <-stop:
			t.Stop()
			return
		}
	}
}
