package app_test

import (
	"testing"
	"website-monitor/app"
	"website-monitor/content_checkers"
)

func TestParse(t *testing.T) {
	data := []byte(`loglevel: debug
global:
  expected_status_code: 200
  interval: 30
  headers:
    User-Agent: "Mozilla/5.0"
  notifiers:
    - name: Slack
      type: slack
      webhook: "https://hooks.slack.com/services/1/2/3"
monitors:
  - name: "A regex expected monitor"
    url: "https://monitor.example/expected"
    regex_expected: "Expected test"
`)

	cfg := &app.Config{}
	err := cfg.Parse(data)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if len(cfg.Monitors) != 1 {
		t.Errorf("Number of monitors was incorrect, got: %d, want: %d.", len(cfg.Monitors), 1)
	}

	monitor := cfg.Monitors[0]
	if len(monitor.ContentChecks) != 1 {
		t.Errorf("Number of content checks was incorrect, got: %d, want: %d.", len(monitor.ContentChecks), 1)
	}

	if len(monitor.Notifiers) != 1 {
		t.Errorf("Number of notifiers was incorrect, got: %d, want: %d.", len(monitor.Notifiers), 1)
	}

	if monitor.ExpectedStatusCode != 200 {
		t.Errorf("ExpectedStatusCode was incorrect, got: %d, want: %d.", monitor.ExpectedStatusCode, 200)
	}

	if monitor.Interval != 30 {
		t.Errorf("Interval was incorrect, got: %d, want: %d.", monitor.Interval, 30)
	}

	if monitor.Headers == nil {
		t.Errorf("Headers was not set")
	}

	if _, ok := monitor.Headers["User-Agent"]; !ok {
		t.Errorf("User-Agent was not set.")
	}

	if monitor.Headers["User-Agent"] != "Mozilla/5.0" {
		t.Errorf("User-Agent was incorrect, got: %s, want: %s.", monitor.Headers["User-Agent"], "Mozilla/5.0")
	}

	if monitor.Url != "https://monitor.example/expected" {
		t.Errorf("URL was incorrect, got: %s, want: %s.", monitor.Url, "https://monitor.example/expected")
	}

	if monitor.Url != monitor.DisplayUrl {
		t.Errorf("DisplayURL was incorrect, got: %s, want: %s.", monitor.Url, monitor.DisplayUrl)
	}

	if _, ok := monitor.Headers["Referer"]; !ok {
		t.Errorf("Referer was not set.")
	}

	if monitor.Headers["Referer"] != "https://monitor.example/" {
		t.Errorf("Referer was incorrect, got: %s, want: %s.", monitor.Headers["Referer"], "https://monitor.example/")
	}

	contentCheck, ok := monitor.ContentChecks[0].(*content_checkers.RegexChecker)
	if !ok {
		t.Errorf("ContentChecker of wrong type")
	}

	if contentCheck.Name != "Expected test" {
		t.Errorf("Content check name was incorrect, got: %s, want: %s.", contentCheck.Name, "Expected test")
	}

	if contentCheck.Regex != "Expected test" {
		t.Errorf("Content check regex was incorrect, got: %s, want: %s.", contentCheck.Name, "Expected test")
	}

	if !contentCheck.ExpectedExisting {
		t.Errorf("Content check expect existing was incorrect, got: %t, want: %t.", contentCheck.ExpectedExisting, true)
	}
}


