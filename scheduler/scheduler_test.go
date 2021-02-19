package scheduler_test

import (
	"testing"
	"time"
	"website-monitor/scheduler"
)

func Test(t *testing.T) {
	tests := []struct {
		name     string
		from     time.Time
		interval time.Duration
		days     []time.Weekday
		hours    []int
		expected time.Time
	}{
		{
			name:     "regular interval",
			from:     time.Date(2021, 02, 17, 10, 00, 00, 00, time.UTC),
			interval: 14 * time.Second,
			expected: time.Date(2021, 02, 17, 10, 00, 14, 00, time.UTC),
		},
		{
			name:     "interval within days, hours",
			from:     time.Date(2021, 02, 17, 10, 00, 00, 00, time.UTC),
			interval: 1 * time.Hour,
			hours:    []int{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			days:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			expected: time.Date(2021, 02, 17, 11, 00, 00, 00, time.UTC),
		},
		{
			name:     "outside of days, no hours",
			from:     time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			interval: 1 * time.Hour,
			days:     []time.Weekday{time.Monday},
			expected: time.Date(2021, 02, 22, 00, 00, 00, 00, time.UTC),
		},
		{
			name:     "interval before hours",
			from:     time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			interval: 1 * time.Hour,
			hours:    []int{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			days:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			expected: time.Date(2021, 02, 17, 07, 00, 00, 00, time.UTC),
		},
		{
			name:     "interval after hours",
			from:     time.Date(2021, 02, 17, 18, 00, 00, 00, time.UTC),
			interval: 1 * time.Hour,
			hours:    []int{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			days:     []time.Weekday{time.Monday, time.Tuesday},
			expected: time.Date(2021, 02, 22, 07, 00, 00, 00, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scheduler.NewScheduler(
				test.interval,
				0,
				test.hours,
				test.days,
			)

			to := s.CalculateNextFrom(test.from)
			if to != test.expected {
				t.Errorf("got: %s, expected: %s", to.String(), test.expected.String())
			}
		})
	}
}

func TestScheduler_IsWithinSchedule(t *testing.T) {
	tests := []struct {
		name     string
		days     []time.Weekday
		hours    []int
		value    time.Time
		expected bool
	}{
		{
			name:     "no days or hours",
			value:    time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			expected: true,
		},
		{
			name:     "no days, not within hours",
			value:    time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			hours:    []int{7, 8},
			expected: false,
		},
		{
			name:     "no days, within hours",
			value:    time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			hours:    []int{2, 8},
			expected: true,
		},
		{
			name:     "no hours, not within days",
			value:    time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			days:     []time.Weekday{time.Monday, time.Tuesday},
			expected: false,
		},
		{
			name:     "no hours, within days",
			value:    time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			days:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday},
			expected: true,
		},
		{
			name:     "within hours and days",
			value:    time.Date(2021, 02, 17, 02, 00, 00, 00, time.UTC),
			days:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday},
			hours:    []int{1, 2, 3},
			expected: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scheduler.NewScheduler(
				0,
				0,
				test.hours,
				test.days,
			)
			result := s.IsWithinSchedule(test.value)
			if result != test.expected {
				t.Errorf("got: %t, expected: %t", result, test.expected)
			}
		})
	}
}
