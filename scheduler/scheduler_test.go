package scheduler_test

import (
	"gopkg.in/yaml.v3"
	"reflect"
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
		{
			name:     "interval after hours, just",
			from:     time.Date(2021, 02, 19, 20, 00, 00, 00, time.UTC),
			interval: 1 * time.Hour,
			hours:    []int{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			days:     []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
			expected: time.Date(2021, 02, 22, 07, 00, 00, 00, time.UTC),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := scheduler.NewScheduler(
				test.interval,
				nil,
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
				nil,
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

func TestScheduler_String(t *testing.T) {
	zero := 0
	twenty := 20
	tests := []struct {
		name      string
		scheduler scheduler.Scheduler
		expected  string
	}{
		{
			name: "interval: 3600s, variation: 0, days: 1-5, hours: 7-17",
			scheduler: *scheduler.NewScheduler(
				3600*time.Second,
				&zero,
				scheduler.Hours{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
				scheduler.Days{1, 2, 3, 4, 5}),
			expected: "3600;0;1,2,3,4,5;7,8,9,10,11,12,13,14,15,16,17",
		},
		{
			name: "interval: 600s, variation: 20, days: 1,3, hours: 8,9,11",
			scheduler: *scheduler.NewScheduler(
				600*time.Second,
				&twenty,
				scheduler.Hours{8, 9, 11},
				scheduler.Days{1, 3}),
			expected: "600;20;1,3;8,9,11",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.scheduler.String()
			if got != test.expected {
				t.Errorf("got: %s, expected: %s", got, test.expected)
			}
		})
	}
}

func TestScheduler_FromString(t *testing.T) {
	zero := 0
	twenty := 20
	tests := []struct {
		name     string
		input    string
		expected *scheduler.Scheduler
	}{
		{
			name:  "interval: 3600, variation: 0, days: 1-5, hours: 7-17",
			input: "3600;0;1,2,3,4,5;7,8,9,10,11,12,13,14,15,16,17",
			expected: scheduler.NewScheduler(
				3600*time.Second,
				&zero,
				scheduler.Hours{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
				scheduler.Days{1, 2, 3, 4, 5}),
		},
		{
			name:  "interval: 600, variation: 20, days: 1,3, hours: 8,9,11",
			input: "600;20;1,3;8,9,11",
			expected: scheduler.NewScheduler(
				600*time.Second,
				&twenty,
				scheduler.Hours{8, 9, 11},
				scheduler.Days{1, 3}),
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := scheduler.NewSchedulerFromString(test.input)
			if err != nil {
				t.Errorf("unexpected err: %v", err)
			}
			if !reflect.DeepEqual(got, test.expected) {
				t.Errorf("got: %v, expected: %v", got, test.expected)
			}
		})
	}
}

func TestDays_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected *scheduler.Days
	}{
		{
			name:     "1-5",
			data:     "1-5",
			expected: &scheduler.Days{1, 2, 3, 4, 5},
		},
		{
			name:     "1,2,3",
			data:     "1,2,3",
			expected: &scheduler.Days{1, 2, 3},
		},
		{
			name:     "1,2,3,5-6",
			data:     "1,2,3,5-6",
			expected: &scheduler.Days{1, 2, 3, 5, 6},
		},
		{
			name:     "1,2-4,6",
			data:     "1,2-4,6",
			expected: &scheduler.Days{1, 2, 3, 4, 6},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := &scheduler.Days{}
			err := yaml.Unmarshal([]byte(test.data), d)
			if err != nil {
				t.Errorf("got err %v, expected nil", err)
			}
			if !reflect.DeepEqual(d, test.expected) {
				t.Errorf("got '%v' expected '%v'", d, test.expected)
			}
		})
	}
}

func TestHours_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		expected *scheduler.Hours
	}{
		{
			name:     "1-5",
			data:     "1-5",
			expected: &scheduler.Hours{1, 2, 3, 4, 5},
		},
		{
			name:     "1,2,3",
			data:     "1,2,3",
			expected: &scheduler.Hours{1, 2, 3},
		},
		{
			name:     "1,2,3,5-6",
			data:     "1,2,3,5-6",
			expected: &scheduler.Hours{1, 2, 3, 5, 6},
		},
		{
			name:     "1,2-4,6",
			data:     "1,2-4,6",
			expected: &scheduler.Hours{1, 2, 3, 4, 6},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			d := &scheduler.Hours{}
			err := yaml.Unmarshal([]byte(test.data), d)
			if err != nil {
				t.Errorf("got err %v, expected nil", err)
			}
			if !reflect.DeepEqual(d, test.expected) {
				t.Errorf("got '%v' expected '%v'", d, test.expected)
			}
		})
	}
}

func TestHours_String(t *testing.T) {
	tests := []struct {
		name string
		h    scheduler.Hours
		want string
	}{
		{
			name: "simple list of hours",
			h: scheduler.Hours{1,2,3,4,5},
			want: "1,2,3,4,5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.h.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}