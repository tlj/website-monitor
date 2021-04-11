package notifier

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"website-monitor/messagequeue"
	"website-monitor/monitors"
	"website-monitor/notifiers"
)

type Notifier struct {
	monitorRepo  monitors.RepositoryInterface
	notifierRepo notifiers.RepositoryInterface
}

func NewNotifier(monitorRepo monitors.RepositoryInterface, notifierRepo notifiers.RepositoryInterface) *Notifier {
	return &Notifier{
		monitorRepo:  monitorRepo,
		notifierRepo: notifierRepo,
	}
}

func (n *Notifier) Run(msgs <-chan amqp.Delivery, stop chan bool) {
	for {
		select {
		case d := <-msgs:
			log.Printf("Message received: %s", d.Body)

			cr := &messagequeue.CrawlResult{}
			err := json.Unmarshal(d.Body, cr)
			if err != nil {
				log.Errorf("unable to unmarshal message '%s': %s", d.Body, err)
				continue
			}

			mon, err := n.monitorRepo.Find(cr.MonitorID)
			if err != nil {
				log.Error(err)
				continue
			}

			nfs, err := n.notifierRepo.FindByMonitorId(cr.MonitorID)
			if err != nil {
				log.Error(err)
				continue
			}

			for _, not := range nfs {
				if err := not.Notify(mon.Name, mon.DisplayUrl, cr); err != nil {
					log.Errorf("unable to send notification to '%s': %s", not.Name(), err)
				}
			}
		case <-stop:
			return
		}
	}
}
