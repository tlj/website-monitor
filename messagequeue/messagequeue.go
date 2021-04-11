package messagequeue

import (
	"fmt"
	"github.com/streadway/amqp"
)

type MessageQueue struct {
	Conn              *amqp.Connection
	Ch                *amqp.Channel
	NotificationQueue amqp.Queue
	ScheduleQueue     amqp.Queue
}

type MessageQueuePublisher interface {
	Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error
}

func NewMessageQueue(messageQueueURL string) (*MessageQueue, error) {
	mq := &MessageQueue{}

	var err error
	mq.Conn, err = amqp.Dial(messageQueueURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %s", err)
	}

	mq.Ch, err = mq.Conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %s", err)
	}

	err = mq.Ch.Qos(
		1,
		0,
		false)
	if err != nil {
		return nil, fmt.Errorf("failed to set channel QOS: %s", err)
	}

	mq.NotificationQueue, err = mq.Ch.QueueDeclare(
		"notifications",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %s", err)
	}

	mq.ScheduleQueue, err = mq.Ch.QueueDeclare(
		"schedule",
		false,
		false,
		false,
		false,
		nil)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %s", err)
	}

	return mq, nil
}

func (mq *MessageQueue) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	return mq.Ch.Publish(exchange, key, mandatory, immediate, msg)
}

func (mq *MessageQueue) Close() error {
	if err := mq.Ch.Close(); err != nil {
		return err
	}

	if err := mq.Conn.Close(); err != nil {
		return err
	}

	return nil
}
