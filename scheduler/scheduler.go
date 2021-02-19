package scheduler

import (
	log "github.com/sirupsen/logrus"
	"math/rand"
	"sort"
	"time"
)

type Scheduler struct {
	interval                    time.Duration
	intervalVariationPercentage int
	hours                       []int
	days                        []time.Weekday
}

func NewScheduler(interval time.Duration, intervalVariationPercentage int, hours []int, days []time.Weekday) *Scheduler {
	s := &Scheduler{
		interval: interval,
		intervalVariationPercentage: intervalVariationPercentage,
		hours: hours,
		days: days,
	}
	sort.Ints(s.hours)

	return s
}

func (s *Scheduler) IsWithinDays(t time.Time) bool {
	if len(s.days) == 0 {
		return true
	}

	for _, d := range s.days {
		if t.Weekday() == d {
			return true
		}
	}

	return false
}

func (s *Scheduler) IsWithinHours(t time.Time) bool {
	if len(s.hours) == 0 {
		return true
	}

	for _, h := range s.hours {
		if t.Hour() == h {
			return true
		}
	}

	return false
}

func (s *Scheduler) IsWithinSchedule(t time.Time) bool {
	return s.IsWithinDays(t) && s.IsWithinHours(t)
}

func (s *Scheduler) NextDay(from time.Time) time.Time {
	to := from

	if len(s.days) == 0 {
		return time.Date(to.Year(), to.Month(), to.Day(), 00, 00, 00, 00, time.UTC)
	}

	if !s.IsWithinDays(to) {
		for {
			to = to.Add(24 * time.Hour)
			if s.IsWithinDays(to) {
				break
			}
		}
		to = time.Date(to.Year(), to.Month(), to.Day(), 00, 00, 00, 00, time.UTC)
	}

	return to
}

func (s *Scheduler) NextHour(from time.Time) time.Time {
	to := from

	if from.Hour() > s.hours[len(s.hours) - 1] {
		to = s.NextDay(to)
	}

	if !s.IsWithinHours(to) {
		for {
			to = to.Add(1 * time.Hour)
			if s.IsWithinHours(to) {
				break
			}
		}
		to = time.Date(to.Year(), to.Month(), to.Day(), to.Hour(), 00, 00, 00, time.UTC)
	}

	return to
}

func (s *Scheduler) CalculateNextFrom(from time.Time) time.Time {
	incBy := s.interval.Seconds()

	if s.intervalVariationPercentage > 0 {
		var p float64
		p = incBy * (float64(s.intervalVariationPercentage) / 100)
		min := incBy - p
		max := incBy + p
		rand.Seed(time.Now().UnixNano())
		iv := rand.Intn(int(max)-int(min)+1) + int(min)
		log.Printf("delay %d...", iv)
	}

	to := from.Add(s.interval)

	if s.IsWithinSchedule(to) {
		return to
	}

	if !s.IsWithinDays(to) {
		to = s.NextDay(to)
	}

	if !s.IsWithinHours(to) {
		to = s.NextHour(to)
	}

	return to
}