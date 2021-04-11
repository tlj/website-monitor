package content_checkers_test

import (
	"strings"
	"testing"
	"website-monitor/content_checkers"
)

func TestHtmlXPathChecker_Check(t *testing.T) {
	tests := []struct {
		name             string
		data             string
		path             string
		expected         string
		expectedExisting bool
		result           bool
		err              error
	}{
		{
			name:             "expected, found",
			data:             `<html><body><div id="whatever"><h1>Expected</h1></div></body></html>`,
			path:             `//div[@id="whatever"]/h1/text()`,
			expected:         "Expected",
			expectedExisting: true,
			result:           true,
			err:              nil,
		},
		{
			name:             "expected, not found",
			data:             `<html><body><div id="whatever"><h1>Expected</h1></div></body></html>`,
			path:             `//div[@id="whatever"]/h1/text()`,
			expected:         "Not Expected",
			expectedExisting: true,
			result:           false,
			err:              nil,
		},
		{
			name:             "expected, no key",
			data:             `<html><body><div id="whatever"><h1>Expected</h1></div></body></html>`,
			path:             `//div[@id="whatever2"]/h1/text()`,
			expected:         "Expected",
			expectedExisting: true,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, found",
			data:             `<html><body><div id="whatever"><h1>Expected</h1></div></body></html>`,
			path:             `//div[@id="whatever"]/h1/text()`,
			expected:         "Expected",
			expectedExisting: false,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, found",
			data:             `<html><body><div id="whatever"><h1>Expected</h1></div></body></html>`,
			path:             `//div[@id="whatever"]/h1/text()`,
			expected:         "Expected",
			expectedExisting: false,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, no key",
			data:             `<html><body><div id="whatever"><h1>Expected</h1></div></body></html>`,
			path:             `//div[@id="whatever2"]/h1/text()`,
			expected:         "Expected",
			expectedExisting: false,
			result:           true,
			err:              nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := content_checkers.NewHtmlXPathChecker(test.name, test.path, test.expected, test.expectedExisting)
			res, err := c.Check(strings.NewReader(test.data))
			if err != test.err {
				t.Errorf("%s - got err %v, expected %v", test.name, err, test.err)
			}

			if res != test.result {
				t.Errorf("%s - got result %t, expected %t", test.name, res, test.result)
			}
		})
	}
}
