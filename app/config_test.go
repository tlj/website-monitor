package app_test

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"time"
	"website-monitor/app"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
	"website-monitor/notifiers"
	"website-monitor/scheduler"

	"github.com/google/go-cmp/cmp"
)

func TestParseConfig(t *testing.T) {
	data := []byte(`loglevel: debug
global:
  expected_status_code: 200
  interval: 30
  interval_variable_percentage: 20
  schedule:
    days: "1-5"
    hours: "10-17"
  render_server_urn: "ws://localhost:9222"
  headers:
    User-Agent: "Mozilla/5.0"
  notifiers:
    - name: Slack
      type: slack
      webhook: "https://hooks.slack.com/services/1/2/3"
monitors:
  - name: "A regex expected monitor with default settings"
    url: "https://monitor.example/expected"
    regex_expected: "Expected test"
  - name: "A regex expected monitor with default settings, expanded"
    url: "https://monitor.example/expected"
    content_checks:
    - name: Expected test
      type: "Regex"
      expected: "Expected test"
  - name: "A regex expected monitor without schedule"
    url: "https://monitor.example/expected"
    regex_expected: "Expected test"
    schedule: {}
  - name: "A regex unexpected monitor with custom intervals"
    url: "https://monitor.example/expected"
    regex_expected: "Expected test"
    interval: 60
    interval_variable_percentage: 0
  - name: "A Json Path example"
    url: "https://monitor.example/data.json"
    display_url: "https://monitor.example/data"
    content_checks:
    - name: Expected test
      type: "JsonPath"
      path: "//SomeProperty"
      expected: "whatever"
  - name: "A HTML XPath example"
    url: "https://monitor.example/data.html"
    content_checks:
    - name: Expected test
      type: "HtmlXPath"
      path: "//h1"
      expected: "whatever"
  - name: "A http render example"
    url: "https://monitor.example/httprender"
    type: http_render
    content_checks:
    - name: Css Selector Check
      type: HtmlRenderSelector
      path: "html h1"
      not_expected: "Expected header"
  - name: "Two regex tests, only one is required to report success"
    url: "https://monitor.example/expected"
    require_some: true
    content_checks:
    - name: Either this test
      type: "Regex"
      expected: "Expected test"
    - name: Or this test
      type: "Regex"
      expected: "Alternative test"
`)

	expectedIntervalVariable := 20
	zeroInt := 0
	defaultSchedule := scheduler.NewScheduler(
		30*time.Second,
		&expectedIntervalVariable,
		[]int{10, 11, 12, 13, 14, 15, 16, 17},
		[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
	)
	expectedCfg := &app.Config{
		LogLevel: "debug",
		Monitors: []*monitors.Check{
			{
				Name:               "A regex expected monitor with default settings",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Type:               monitors.HttpMonitorType,
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          defaultSchedule,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Expected test", "Expected test", true),
				},
			},
			{
				Name:               "A regex expected monitor with default settings, expanded",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Type:               monitors.HttpMonitorType,
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          defaultSchedule,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Expected test", "Expected test", true),
				},
			},
			{
				Name:               "A regex expected monitor without schedule",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Type:               monitors.HttpMonitorType,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          scheduler.NewScheduler(time.Duration(30) * time.Second, &expectedIntervalVariable, nil, nil),
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Expected test", "Expected test", true),
				},
			},
			{
				Name:               "A regex unexpected monitor with custom intervals",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Type:               monitors.HttpMonitorType,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler: scheduler.NewScheduler(
					time.Duration(60)*time.Second,
					&zeroInt,
					[]int{10, 11, 12, 13, 14, 15, 16, 17},
					[]time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday},
				),
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Expected test", "Expected test", true),
				},
			},
			{
				Name:               "A Json Path example",
				Url:                "https://monitor.example/data.json",
				DisplayUrl:         "https://monitor.example/data",
				RenderServerURN:    "ws://localhost:9222",
				Type:               monitors.HttpMonitorType,
				Headers:            map[string]string{"Referer": "https://monitor.example/data", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          defaultSchedule,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewJsonPathChecker("Expected test", "//SomeProperty", "whatever", true),
				},
			},
			{
				Name:               "A HTML XPath example",
				Url:                "https://monitor.example/data.html",
				DisplayUrl:         "https://monitor.example/data.html",
				RenderServerURN:    "ws://localhost:9222",
				Type:               monitors.HttpMonitorType,
				Headers:            map[string]string{"Referer": "https://monitor.example/data.html", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          defaultSchedule,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewHtmlXPathChecker("Expected test", "//h1", "whatever", true),
				},
			},
			{
				Name:               "A http render example",
				Url:                "https://monitor.example/httprender",
				DisplayUrl:         "https://monitor.example/httprender",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/httprender", "User-Agent": "Mozilla/5.0"},
				Type:               monitors.HttpRenderMonitorType,
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          defaultSchedule,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewHtmlRenderSelectorChecker("Css Selector Check", "html h1", "Expected header", false),
				},
			},
			{
				Name:               "Two regex tests, only one is required to report success",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				Type:               monitors.HttpMonitorType,
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{notifiers.NewSlackNotifier("https://hooks.slack.com/services/1/2/3")},
				Scheduler:          defaultSchedule,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Either this test", "Expected test", true),
					content_checkers.NewRegexChecker("Or this test", "Alternative test", true),
				},
			},
		},
	}

	cfg := &app.Config{}
	err := cfg.Parse(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	diff := cmp.Diff(
		cfg,
		expectedCfg,
		cmpopts.IgnoreUnexported(scheduler.Scheduler{}, monitors.Check{}, content_checkers.HtmlRenderSelectorChecker{}),
	)
	if diff != "" {
		t.Error(diff)
	}
}
