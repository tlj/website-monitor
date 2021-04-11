package scheduler_test

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/mock"
	"testing"
	"website-monitor/messagequeue"
	"website-monitor/scheduler"
)

type MonitorMock struct {
	mock.Mock
}

func (m *MonitorMock) ShouldUpdate() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MonitorMock) GetId() int64 {
	args := m.Called()
	return int64(args.Int(0))
}

func (m *MonitorMock) GetName() string {
	args := m.Called()
	return args.String(0)
}

type MessageQueueMock struct {
	mock.Mock
}

func (m *MessageQueueMock) Publish(exchange, key string, mandatory, immediate bool, msg amqp.Publishing) error {
	args := m.Called(exchange, key, mandatory, immediate, msg)
	return args.Error(0)
}

func TestScheduler_HandleMonitor(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	mon := &MonitorMock{}

	msg, _ := json.Marshal(messagequeue.ScheduleJob{
		MonitorID: 1,
		Name:      "Test",
		Type:      "monitor",
	})

	schedEmitter := scheduler.NewScheduleEmitter()
	mq := &MessageQueueMock{}

	mon.On("ShouldUpdate").Return(true)
	mon.On("GetId").Return(1)
	mon.On("GetName").Return("Test")
	mq.On("Publish", "", "notifications", false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	}).Return(nil)
	mon.On("ShouldUpdate").Return(false)

	schedEmitter.HandleMonitor(mon, mq, "", "notifications")
	schedEmitter.HandleMonitor(mon, mq, "", "notifications")

	mq.AssertExpectations(t)
	mon.AssertExpectations(t)
}

