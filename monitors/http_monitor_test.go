package monitors_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"website-monitor/content_checkers"
	"website-monitor/monitors"
)

func TestHttpMonitor_Check(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		checkers []content_checkers.ContentChecker
		result   bool
		err      error
	}{
		{
			name: "one checker, success",
			data: "Some sort of text test.",
			checkers: []content_checkers.ContentChecker{
				content_checkers.NewRegexChecker("regex", "sort of text", true),
			},
			result: true,
			err:    nil,
		},
		{
			name: "one checker, fail",
			data: "Some sort of text test.",
			checkers: []content_checkers.ContentChecker{
				content_checkers.NewRegexChecker("regex", "no such text", true),
			},
			result: false,
			err:    nil,
		},
		{
			name: "two checkers - success and fail equals fail",
			data: "Some sort of text test.",
			checkers: []content_checkers.ContentChecker{
				content_checkers.NewRegexChecker("regex", "sort of text", true),
				content_checkers.NewRegexChecker("regex", "sort of text", false),
			},
			result: false,
			err:    nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			webCalls := 0
			ts := httptest.NewServer(
				http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					webCalls++
					_, _ = fmt.Fprintln(w, test.data)
				}))
			defer ts.Close()

			ch := monitors.Check{
				Name:               test.name,
				Url:                ts.URL,
				ExpectedStatusCode: 200,
				ContentChecks:      test.checkers,
			}
			hm := monitors.HttpMonitor{}
			res, err := hm.Check(ch)
			if err != test.err {
				t.Errorf("got err: %v, expected nil", test.err)
			}

			if res == nil {
				t.Fatal("result is nil")
			}

			if res.AllTrue() != test.result {
				t.Errorf("got allTrue: %t, expected %t", res.AllTrue(), test.result)
			}

			if webCalls != 1 {
				t.Errorf("expected webCalls: %d, expected: %d", webCalls, 1)
			}
		})
	}
}
