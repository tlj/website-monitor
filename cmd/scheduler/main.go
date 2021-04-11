package main

import (
	"github.com/benbjohnson/clock"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"time"
	"website-monitor/content_checkers"
	"website-monitor/messagequeue"
	"website-monitor/monitors"
	"website-monitor/scheduler"
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

	schedEmitter := scheduler.NewScheduleEmitter()
	sched := scheduler.NewScheduler(mq, monitorRepo, schedEmitter, "", mq.ScheduleQueue.Name)

	forever := make(chan bool)

	go func() {
		cl := clock.New()
		stop := make(chan bool)
		sched.Run(cl, 1 * time.Second, stop)
	}()

	log.Printf(" [*] Scheduling. To exit press CTRL-C.")

	<-forever
}
