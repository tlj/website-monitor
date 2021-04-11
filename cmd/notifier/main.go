package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"website-monitor/content_checkers"
	"website-monitor/messagequeue"
	"website-monitor/monitors"
	"website-monitor/notifier"
	"website-monitor/notifiers"
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

	notifierRepo := notifiers.NewRepository(db)
	checkerRepo := content_checkers.NewRepository(db)
	monitorRepo := monitors.NewRepository(db, checkerRepo)

	mq, err := messagequeue.NewMessageQueue(s.MessageQueueUrl)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := mq.Ch.Consume(
		mq.NotificationQueue.Name,
		"",
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

	noti := notifier.NewNotifier(monitorRepo, notifierRepo)
	stop := make(chan bool)

	go func() {
		noti.Run(msgs, stop)
	}()

	log.Printf(" [*] Waiting for message to send as notifications. To exit press CTRL-C.")

	<-forever
}
