package scheduler

import (
	"fmt"
	"math/rand"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Days []time.Weekday

func (d Days) String() string {
	dayStr := ""
	for _, day := range d {
		if dayStr != "" {
			dayStr += ","
		}
		dayStr += fmt.Sprintf("%d", int(day))
	}

	return dayStr
}

func (d *Days) FromString(str string) error {
	for _, s := range strings.Split(str, ",") {
		s = strings.Trim(s, " ")
		if s == "" {
			continue
		}

		if strings.Contains(s, "-") {
			fromTo := strings.Split(s, "-")
			from, err := strconv.Atoi(fromTo[0])
			if err != nil {
				return fmt.Errorf("invalid from in range '%s': %v", fromTo, err)
			}
			to, err := strconv.Atoi(fromTo[1])
			if err != nil {
				return fmt.Errorf("invalid to in range '%s': %v", fromTo, err)
			}
			for i := from; i <= to; i++ {
				*d = append(*d, time.Weekday(i))
			}

			continue
		}

		di, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*d = append(*d, time.Weekday(di))
	}

	return nil
}

func (d *Days) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	return d.FromString(str)
}

type Hours []int

func (h Hours) String() string {
	hourStr := ""
	for _, hour := range h {
		if hourStr != "" {
			hourStr += ","
		}
		hourStr += fmt.Sprintf("%d", hour)
	}

	return hourStr
}

func (h *Hours) FromString(str string) error {
	for _, s := range strings.Split(str, ",") {
		s = strings.Trim(s, " ")
		if s == "" {
			continue
		}

		if strings.Contains(s, "-") {
			fromTo := strings.Split(s, "-")
			from, err := strconv.Atoi(fromTo[0])
			if err != nil {
				return fmt.Errorf("invalid from in range '%s': %v", fromTo, err)
			}
			to, err := strconv.Atoi(fromTo[1])
			if err != nil {
				return fmt.Errorf("invalid to in range '%s': %v", fromTo, err)
			}
			for i := from; i <= to; i++ {
				*h = append(*h, i)
			}

			continue
		}

		di, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*h = append(*h, di)
	}

	return nil
}

func (h *Hours) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var str string
	if err := unmarshal(&str); err != nil {
		return err
	}

	return h.FromString(str)
}

type Scheduler struct {
	Interval                    time.Duration `yaml:"interval"`
	IntervalVariationPercentage *int          `yaml:"interval_variation_percentage"`
	Hours                       Hours         `yaml:"hours"`
	Days                        Days          `yaml:"days"`
}

func NewScheduler(interval time.Duration, intervalVariationPercentage *int, hours Hours, days Days) *Scheduler {
	s := &Scheduler{
		Interval:                    interval,
		IntervalVariationPercentage: intervalVariationPercentage,
		Hours:                       hours,
		Days:                        days,
	}
	sort.Ints(s.Hours)

	return s
}

func NewSchedulerFromString(input string) (*Scheduler, error) {
	parts := strings.Split(input, ";")
	if len(parts) != 4 {
		return nil, fmt.Errorf("need 4 parts, only %d given", len(parts))
	}

	interval, err := strconv.Atoi(parts[0])
	if err != nil {
		return nil, err
	}

	variation, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	days := &Days{}
	hours := &Hours{}

	if err := days.FromString(parts[2]); err != nil {
		return nil, err
	}

	if err := hours.FromString(parts[3]); err != nil {
		return nil, err
	}

	return &Scheduler{
		Interval:                    time.Duration(interval) * time.Second,
		IntervalVariationPercentage: &variation,
		Days:                        *days,
		Hours:                       *hours,
	}, nil
}

func (s *Scheduler) IsWithinDays(t time.Time) bool {
	if len(s.Days) == 0 {
		return true
	}

	for _, d := range s.Days {
		if t.Weekday() == d {
			return true
		}
	}

	return false
}

func (s *Scheduler) IsWithinHours(t time.Time) bool {
	if len(s.Hours) == 0 {
		return true
	}

	for _, h := range s.Hours {
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

	if len(s.Days) == 0 {
		return time.Date(to.Year(), to.Month(), to.Day(), 00, 00, 00, 00, time.UTC)
	}

	to = to.Add(24 * time.Hour)
	to = time.Date(to.Year(), to.Month(), to.Day(), 00, 00, 00, 00, time.UTC)

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

	if from.Hour() > s.Hours[len(s.Hours)-1] {
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
	incBy := int(s.Interval.Seconds())

	if s.IntervalVariationPercentage != nil && *s.IntervalVariationPercentage > 0 {
		var p float64
		p = float64(incBy) * (float64(*s.IntervalVariationPercentage) / 100)
		min := incBy - int(p)
		max := incBy + int(p)
		rand.Seed(time.Now().UnixNano())
		incBy = rand.Intn(max-min+1) + min
	}

	to := from.Add(time.Duration(incBy) * time.Second)

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

func (s *Scheduler) Equal(y *Scheduler) bool {
	if s == nil && y == nil {
		return true
	}
	if (s == nil && y != nil) || (s != nil && y == nil) {
		return false
	}
	if s.Interval != y.Interval {
		return false
	}

	if s.IntervalVariationPercentage == nil && y.IntervalVariationPercentage != nil {
		return false
	}

	if s.IntervalVariationPercentage != nil && y.IntervalVariationPercentage == nil {
		return false
	}

	if *s.IntervalVariationPercentage != *y.IntervalVariationPercentage {
		return false
	}

	if !reflect.DeepEqual(s.Days, y.Days) {
		return false
	}
	if !reflect.DeepEqual(s.Hours, y.Hours) {
		return false
	}

	return true
}

func (s *Scheduler) String() string {
	variation := 0
	if s.IntervalVariationPercentage != nil {
		variation = *s.IntervalVariationPercentage
	}
	str := fmt.Sprintf(
		"%d;%d;%s;%s",
		s.Interval/1000000000,
		variation,
		s.Days,
		s.Hours)

	return str
}

func (s *Scheduler) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type alias struct {
		Interval                    string `yaml:"interval"`
		IntervalVariationPercentage *int   `yaml:"interval_variation_percentage"`
		Days                        Days   `yaml:"days"`
		Hours                       Hours  `yaml:"hours"`
	}

	var tmp alias
	if err := unmarshal(&tmp); err != nil {
		return err
	}

	t, err := time.ParseDuration(tmp.Interval)
	if err != nil {
		return fmt.Errorf("failed to parse interval '%s' to time.Duration: %v", tmp.Interval, err)
	}

	s.Interval = t
	s.Hours = tmp.Hours
	s.Days = tmp.Days
	s.IntervalVariationPercentage = tmp.IntervalVariationPercentage

	return nil
}
