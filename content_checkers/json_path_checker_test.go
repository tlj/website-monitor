package content_checkers_test

import (
	"strings"
	"testing"
	"website-monitor/content_checkers"
)

func TestJsonPathChecker_Check(t *testing.T) {
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
			data:             `{"somedata":{"somekey":"somevalue"}}`,
			path:             "//somedata/somekey",
			expected:         "somevalue",
			expectedExisting: true,
			result:           true,
			err:              nil,
		},
		{
			name:             "expected, not matching",
			data:             `{"somedata":{"somekey":"somevalue"}}`,
			path:             "//somedata/somekey",
			expected:         "anothervalue",
			expectedExisting: true,
			result:           false,
			err:              nil,
		},
		{
			name:             "expected, key not exists",
			data:             `{"somedata":{"somekey":"somevalue"}}`,
			path:             "//somedata/anotherkey",
			expected:         "somevalue",
			expectedExisting: true,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, found",
			data:             `{"somedata":{"somekey":"somevalue"}}`,
			path:             "//somedata/somekey",
			expected:         "somevalue",
			expectedExisting: false,
			result:           false,
			err:              nil,
		},
		{
			name:             "not expected, key not exists",
			data:             `{"somedata":{"somekey":"somevalue"}}`,
			path:             "//somedata/anotherkey",
			expected:         "somevalue",
			expectedExisting: false,
			result:           true,
			err:              nil,
		},
		{
			name:             "not expected, not found",
			data:             `{"somedata":{"somekey":"somevalue"}}`,
			path:             "//somedata/somekey",
			expected:         "whatever value",
			expectedExisting: false,
			result:           true,
			err:              nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := content_checkers.NewJsonPathChecker(test.name, test.path, test.expected, test.expectedExisting)
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
