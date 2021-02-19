package app_test

import (
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"website-monitor/app"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
	"website-monitor/notifiers"

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
  - name: "A http render example"
    url: "https://monitor.example/httprender"
    type: http_render
    content_checks:
    - name: Css Selector Check
      type: HtmlRenderSelector
      path: "html h1"
      not_expected: "Expected header"
`)

	expectedIntervalVariable := 20
	zeroInt := 0
	expectedCfg := &app.Config{
		LogLevel: "debug",
		Global: app.ConfigGlobal{
			Headers:                    map[string]string{"User-Agent": "Mozilla/5.0"},
			ExpectedStatusCode:         200,
			Interval:                   30,
			IntervalVariablePercentage: &expectedIntervalVariable,
			RenderServerURN:            "ws://localhost:9222",
			Schedule: &monitors.Schedule{
				Days:  "1-5",
				Hours: "10-17",
			},
		},
		Monitors: []*monitors.Check{
			{
				Name:               "A regex expected monitor with default settings",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				RegexExpected:      "Expected test",
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{&notifiers.SlackNotifier{WebhookUrl: "https://hooks.slack.com/services/1/2/3"}},
				Interval:                   30,
				IntervalVariablePercentage: &expectedIntervalVariable,
				Schedule:                   &monitors.Schedule{Days: "1-5", Hours: "10-17"},
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Expected test", "Expected test", true),
				},
			},
			{
				Name:               "A regex expected monitor with default settings, expanded",
				Url:                "https://monitor.example/expected",
				DisplayUrl:         "https://monitor.example/expected",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/expected", "User-Agent": "Mozilla/5.0"},
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{&notifiers.SlackNotifier{WebhookUrl: "https://hooks.slack.com/services/1/2/3"}},
				Interval:                   30,
				IntervalVariablePercentage: &expectedIntervalVariable,
				Schedule:                   &monitors.Schedule{Days: "1-5", Hours: "10-17"},
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
				RegexExpected:      "Expected test",
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{&notifiers.SlackNotifier{WebhookUrl: "https://hooks.slack.com/services/1/2/3"}},
				Interval:                   30,
				IntervalVariablePercentage: &expectedIntervalVariable,
				Schedule:                   &monitors.Schedule{},
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
				RegexExpected:      "Expected test",
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{&notifiers.SlackNotifier{WebhookUrl: "https://hooks.slack.com/services/1/2/3"}},
				Interval:                   60,
				IntervalVariablePercentage: &zeroInt,
				Schedule:                   &monitors.Schedule{Days: "1-5", Hours: "10-17"},
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("Expected test", "Expected test", true),
				},
			},
			{
				Name:               "A http render example",
				Url:                "https://monitor.example/httprender",
				DisplayUrl:         "https://monitor.example/httprender",
				RenderServerURN:    "ws://localhost:9222",
				Headers:            map[string]string{"Referer": "https://monitor.example/httprender", "User-Agent": "Mozilla/5.0"},
				Type:               "http_render",
				ExpectedStatusCode: 200,
				Notifiers:          []notifiers.Notifier{&notifiers.SlackNotifier{WebhookUrl: "https://hooks.slack.com/services/1/2/3"}},
				Interval:                   30,
				IntervalVariablePercentage: &expectedIntervalVariable,
				Schedule:                   &monitors.Schedule{Days: "1-5", Hours: "10-17"},
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewHtmlRenderSelectorChecker("Css Selector Check", "html h1", "Expected header", true),
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
		cmpopts.IgnoreUnexported(content_checkers.HtmlRenderSelectorChecker{}, monitors.Check{}),
		cmpopts.IgnoreFields(monitors.Check{}, "ContentChecksConfig", "NotifiersConfig"),
		cmpopts.IgnoreFields(app.ConfigGlobal{}, "NotifiersConfig"),
		)
	if diff != "" {
		t.Error(diff)
	}
}
