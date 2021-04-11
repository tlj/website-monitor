package scheduler_test

import (
	"github.com/benbjohnson/clock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/mock"
	"runtime"
	"testing"
	"time"
	"website-monitor/interfaces"
	"website-monitor/messagequeue"
	"website-monitor/monitors"
	"website-monitor/schedule"
	"website-monitor/scheduler"
)

type MonitorRepositoryMock struct {
	mock.Mock
}

func (m *MonitorRepositoryMock) Find(id int64) (*monitors.Monitor, error) {
	panic("implement me")
}

func (m *MonitorRepositoryMock) All() ([]*monitors.Monitor, error) {
	args := m.Called()
	return args.Get(0).([]*monitors.Monitor), args.Error(1)
}

func (m *MonitorRepositoryMock) FindByUserId(userId int64) ([]*monitors.Monitor, error) {
	panic("implement me")
}

type ScheduleEmitterMock struct {
	mock.Mock
}

func (s *ScheduleEmitterMock) HandleMonitor(mon interfaces.ShouldUpdaterIdName, publisher messagequeue.MessageQueuePublisher, exchange string, routingKey string) {
	s.Called(mon, publisher, exchange, routingKey)
}

func TestScheduler_Run(t *testing.T) {
	log.SetLevel(log.FatalLevel)

	cl := clock.NewMock()
	stop := make(chan bool)

	zero := 0
	monTest1 := &monitors.Monitor{
		Id:   1,
		Name: "Test",
		Type: "monitor",
		Scheduler: &schedule.Schedule{
			Interval:                    1 * time.Second,
			IntervalVariationPercentage: &zero,
			Hours:                       schedule.Hours{},
			Days:                        schedule.Days{},
		},
	}
	monTest2 := &monitors.Monitor{
			Id:                 2,
			Name:               "Test2",
			Type:               "monitor",
			Scheduler: &schedule.Schedule{
				Interval: 1 * time.Second,
				IntervalVariationPercentage: &zero,
				Hours: schedule.Hours{},
				Days: schedule.Days{},
			},
	}

	mq := &MessageQueueMock{}
	monRepo := &MonitorRepositoryMock{}
	monRepo.On("All").Return([]*monitors.Monitor{monTest1, monTest2}, nil)

	em := &ScheduleEmitterMock{}
	em.On("HandleMonitor", monTest1, mq, "", "notifications")
	em.On("HandleMonitor", monTest2, mq, "", "notifications")

	sched := scheduler.NewScheduler(mq, monRepo, em, "", "notifications")

	go sched.Run(cl, 1*time.Second, stop)

	runtime.Gosched()

	cl.Add(3 * time.Second)

	mq.AssertExpectations(t)
	monRepo.AssertExpectations(t)
	em.AssertExpectations(t)
	em.AssertNumberOfCalls(t, "HandleMonitor", 6)
}
