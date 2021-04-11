package content_checkers_test

import (
	"strings"
	"testing"
	"website-monitor/content_checkers"
)

func TestRegexChecker_Check(t *testing.T) {
	tests := []struct {
		name             string
		data             string
		regex            string
		expectedExisting bool
		result           bool
		err              error
	}{
		{
			name:             "expected, found",
			data:             "this is some sort of text to check",
			regex:            "sort of",
			expectedExisting: true,
			result:           true,
			err:              nil,
		},
		{
			name:             "expected, not found",
			data:             "this is some sort of text to check",
			regex:            "hello world",
			expectedExisting: true,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, found",
			data:             "this is some sort of text to check",
			regex:            "sort of",
			expectedExisting: false,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, not found",
			data:             "this is some sort of text to check",
			regex:            "hello world",
			expectedExisting: false,
			result:           true,
			err:              nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := content_checkers.NewRegexChecker(test.name, test.regex, test.expectedExisting)
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
