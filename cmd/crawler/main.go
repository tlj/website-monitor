package main

import (
	"encoding/json"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"website-monitor/content_checkers"
	"website-monitor/messagequeue"
	"website-monitor/monitors"
)

type Specification struct {
	Debug           bool   `default:"false"`
	PostgresDsn     string `required:"true" split_words:"true"`
	MessageQueueUrl string `required:"true" split_words:"true"`
}

func main() {
	var s Specification
	err := envconfig.Process("wm", &s)
	if err != nil {
		log.Fatal(err.Error())
	}

	if s.Debug {
		log.SetLevel(log.DebugLevel)
	}

	db, err := sqlx.Connect("postgres", s.PostgresDsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	checkerRepo := content_checkers.NewRepository(db)
	monitorRepo := monitors.NewRepository(db, checkerRepo)

	mq, err := messagequeue.NewMessageQueue(s.MessageQueueUrl)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := mq.Ch.Consume(
		mq.ScheduleQueue.Name,
		"crawler-1",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %s", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Message received: %s", d.Body)
			sj := messagequeue.ScheduleJob{}
			err := json.Unmarshal(d.Body, &sj)
			if err != nil {
				log.Errorf("unable to unmarshal schedule job '%s': %s", d.Body, err)
				continue
			}

			if sj.Type == "monitor" {
				m, err := monitorRepo.Find(sj.MonitorID)
				if err != nil {
					log.Fatal(err)
				}

				results, err := m.Run()
				if err != nil {
					log.Errorf("unable to run job '%s': %s", sj.Name, err)
					continue
				}

				cr := &messagequeue.CrawlResult{
					MonitorID: sj.MonitorID,
					Name:      m.Name,
					Type:      sj.Type,
				}

				for _, result := range results.Results {
					cr.Results = append(cr.Results, messagequeue.CrawlCheckerResult{
						Type:    result.ContentChecker.Type(),
						Message: result.ContentChecker.String(),
						Result:  result.Result,
						Err:     result.Err,
					})
				}

				if m.RequireSome {
					cr.Result = results.SomeTrue()
				} else {
					cr.Result = results.AllTrue()
				}

				crJson, err := json.Marshal(cr)
				if err != nil {
					log.Errorf("unable to marshal crawlResult message: %s", err)
					continue
				}

				err = mq.Ch.Publish(
					"", // exchange
					"notifications",              // routing key
					false,           // mandatory
					false,           // immediate
					amqp.Publishing{
						ContentType: "application/json",
						Body:        crJson,
					})
			}
		}
	}()

	log.Printf(" [*] Wating for crawler messages. To exit press CTRL-C.")

	<-forever
}
