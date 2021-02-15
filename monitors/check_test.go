package monitors_test

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
	"website-monitor/notifiers"
)

func TestHttpMonitor_CheckNotification(t *testing.T) {
	notificationCount := 0
	slackServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notificationCount++
			_, _ = fmt.Fprint(w, "ok")
		}))
	defer slackServer.Close()

	tests := []struct {
		name              string
		check             monitors.Check
		notificationCount int
	}{
		{
			name: "send notification when state changes",
			check: monitors.Check{
				Name:               "Test notifications",
				ExpectedStatusCode: 200,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("regex", "sort of text", true),
				},
				Notifiers: []notifiers.Notifier{
					notifiers.NewSlackNotifier(slackServer.URL),
				},
				LastSeenState: false,
			},
			notificationCount: 1,
		},
		{
			name: "dont send notification when state is same",
			check: monitors.Check{
				Name:               "Test notifications",
				ExpectedStatusCode: 200,
				ContentChecks: []content_checkers.ContentChecker{
					content_checkers.NewRegexChecker("regex", "sort of text", true),
				},
				Notifiers: []notifiers.Notifier{
					notifiers.NewSlackNotifier(slackServer.URL),
				},
				LastSeenState: true,
			},
			notificationCount: 0,
		},
	}
	log.SetLevel(log.ErrorLevel)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			contentServer := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					_, _ = fmt.Fprintln(w, "This is some sort of text.")
				}))
			defer contentServer.Close()

			notificationCount = 0

			test.check.Url = contentServer.URL
			err := test.check.Run()
			if err != nil {
				t.Errorf("got err: %v, expected %v", err, nil)
			}

			if notificationCount != test.notificationCount {
				t.Errorf("notificationCount; got %d, expected %d", notificationCount, 1)
			}
		})
	}

}