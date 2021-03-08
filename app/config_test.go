package app_test

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"testing"
	"website-monitor/app"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
	"website-monitor/notifiers"
	"website-monitor/scheduler"
)

func SchedulerWithoutError(str string) *scheduler.Scheduler {
	ret, _ := scheduler.NewSchedulerFromString(str)

	return ret
}

func SlackNotifierWithoutError(name string, options map[string]string) *notifiers.SlackNotifier {
	s, _ := notifiers.NewSlackNotifier(name, options)

	return s
}

func PushSaferNotifierWithoutError(name string, options map[string]string) *notifiers.PushSaferNotifier {
	s, _ := notifiers.NewPushSaferNotifier(name, options)

	return s
}

func TestParseConfig(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		expected *app.Config
		err      error
	}{
		{
			name: "empty schedule",
			data: []byte(`
loglevel: debug
monitors:
  - name: "what"
`),
			expected: &app.Config{
				LogLevel: "debug",
				Monitors: []*monitors.Monitor{
					{
						Name: "what",
						Headers: map[string]string{"Referer":""},
					},
				},
			},
		},
		{
			name: "notifiers",
			data: []byte(`
monitors:
  - name: "notifier"
    notifiers:
    - name: Slack
      type: slack
      options:
        webhook: http://example.com/slack
`),
			expected: &app.Config{
				Monitors: []*monitors.Monitor{
					{
						Name: "notifier",
						Headers: map[string]string{"Referer":""},
						Notifiers: []notifiers.NotifierHolder{
							{
								SlackNotifierWithoutError("Slack", map[string]string{"webhook": "http://example.com/slack"}),
							},
						},
					},
				},
			},
		},
		{
			name: "Full config for simple regex check",
			data: []byte(`
loglevel: debug
monitors:
  - name: "yaml test"
    url: "http://example.com/test"
    type: "http"
    expected_status_code: 200
    checks:
      - name: "A monitor for regex"
        type: regex
        value: "Some monitored text"
        is_expected: true
    schedule: 
      interval: 300s
      interval_variation_percentage: 20
      days: 1,2,3,4,5
      hours: 8,9,10,11,12,13,14,15,16
    headers:
      User-Agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0"
`),
			expected: &app.Config{
				LogLevel: "debug",
				Monitors: []*monitors.Monitor{
					{
						Name:               "yaml test",
						Url:                "http://example.com/test",
						DisplayUrl:         "http://example.com/test",
						Type:               "http",
						ExpectedStatusCode: 200,
						ContentChecks: []content_checkers.ContentCheckerHolder{
							{
								content_checkers.NewRegexChecker("A monitor for regex", "Some monitored text", true),
							},
						},
						Scheduler: SchedulerWithoutError("300;20;1,2,3,4,5;8,9,10,11,12,13,14,15,16"),
						Headers: map[string]string{
							"Referer": "http://example.com/test",
							"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0",
						},
					},
				},
			},
		},
		{
			name: "Defaults",
			data: []byte(`
loglevel: debug
defaults:
  type: "http"
  expected_status_code: 200
  schedule: 
    interval: 300s
    interval_variation_percentage: 20
    days: 1-5
    hours: 7-12
  notifiers:
    - name: Slack
      type: slack
      options:
        webhook: "https://hooks.slack.com/services/X/Y/Z"
    - name: PushSafer
      type: pushsafer
      options:
        private_key: PRIVATEKEY
  headers:
    User-Agent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0"
monitors:
  - name: "yaml test"
    url: "http://example.com/test"
    checks:
      - name: "A monitor for regex"
        type: regex
        value: "Some monitored text"
        is_expected: true
`),
			expected: &app.Config{
				LogLevel: "debug",
				Default: &monitors.Monitor{
					Type: "http",
					ExpectedStatusCode: 200,
					Scheduler: SchedulerWithoutError("300;20;1,2,3,4,5;7,8,9,10,11,12"),
					Headers: map[string]string{
						"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0",
					},
					Notifiers: []notifiers.NotifierHolder{
						{
							SlackNotifierWithoutError("Slack", map[string]string{"webhook": "https://hooks.slack.com/services/X/Y/Z"}),
						},
						{
							PushSaferNotifierWithoutError("PushSafer", map[string]string{"private_key": "PRIVATEKEY"}),
						},
					},
				},
				Monitors: []*monitors.Monitor{
					{
						Name:               "yaml test",
						Url:                "http://example.com/test",
						DisplayUrl:         "http://example.com/test",
						Type:               "http",
						ExpectedStatusCode: 200,
						ContentChecks: []content_checkers.ContentCheckerHolder{
							{
								content_checkers.NewRegexChecker("A monitor for regex", "Some monitored text", true),
							},
						},
						Scheduler: SchedulerWithoutError("300;20;1,2,3,4,5;7,8,9,10,11,12"),
						Headers: map[string]string{
							"Referer": "http://example.com/test",
							"User-Agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.16; rv:85.0) Gecko/20100101 Firefox/85.0",
						},
						Notifiers: []notifiers.NotifierHolder{
							{
								SlackNotifierWithoutError("Slack", map[string]string{"webhook": "https://hooks.slack.com/services/X/Y/Z"}),
							},
							{
								PushSaferNotifierWithoutError("PushSafer", map[string]string{"private_key": "PRIVATEKEY"}),
							},
						},
					},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := &app.Config{}
			err := cfg.LoadConfig(test.data)
			if err != test.err {
				t.Fatalf("got err %v, expected %v", err, test.err)
			}

			diff := cmp.Diff(
				cfg,
				test.expected,
				cmpopts.IgnoreUnexported(scheduler.Scheduler{}, monitors.Monitor{}, content_checkers.HtmlRenderSelectorChecker{}),
			)
			if diff != "" {
				t.Error(diff)
			}
		})
	}
}
