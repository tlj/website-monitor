package monitors_test

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
	"website-monitor/notifiers"
	"website-monitor/scheduler"
)

func SlackNotifierWithoutError(name string, options map[string]string) *notifiers.SlackNotifier {
	s, _ := notifiers.NewSlackNotifier(name, options)

	return s
}

func PushSaferNotifierWithoutError(name string, options map[string]string) *notifiers.PushSaferNotifier {
	s, _ := notifiers.NewPushSaferNotifier(name, options)

	return s
}

func TestHttpMonitor_CheckNotification(t *testing.T) {
	notificationCount := 0
	slackServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			notificationCount++
			_, _ = fmt.Fprint(w, "ok")
		}))
	defer slackServer.Close()

	intZero := 0
	tests := []struct {
		name              string
		check             monitors.Monitor
		notificationCount int
	}{
		{
			name: "send notification when state changes",
			check: monitors.Monitor{
				Name:               "Test notifications",
				ExpectedStatusCode: 200,
				ContentChecks: []content_checkers.ContentCheckerHolder{
					{
						content_checkers.NewRegexChecker("regex", "sort of text", true),
					},
				},
				Notifiers: []notifiers.NotifierHolder{
					{
						SlackNotifierWithoutError("Slack", map[string]string{"webhook": slackServer.URL}),
					},
				},
				Scheduler: scheduler.NewScheduler(time.Duration(30) * time.Second, &intZero, nil, nil),
				LastSeenState: false,
			},
			notificationCount: 1,
		},
		{
			name: "dont send notification when state is same",
			check: monitors.Monitor{
				Name:               "Test notifications",
				ExpectedStatusCode: 200,
				ContentChecks: []content_checkers.ContentCheckerHolder{
					{
						content_checkers.NewRegexChecker("regex", "sort of text", true),
					},
				},
				Notifiers: []notifiers.NotifierHolder{
					{
						SlackNotifierWithoutError("Slack", map[string]string{"webhook": slackServer.URL}),
					},
				},
				Scheduler: scheduler.NewScheduler(time.Duration(30) * time.Second, &intZero, nil, nil),
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