package result_test

import (
	"testing"
	"website-monitor/result"
)

func TestResults_AllTrue(t *testing.T) {
	tests := []struct {
		name     string
		results  result.Results
		expected bool
	}{
		{
			name: "one, true",
			results: result.Results{
				Results: []result.Result{
					{
						Result: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "two, true",
			results: result.Results{
				Results: []result.Result{
					{
						Result: true,
					},
					{
						Result: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "one true, one false",
			results: result.Results{
				Results: []result.Result{
					{
						Result: true,
					},
					{
						Result: false,
					},
				},
			},
			expected: false,
		},
		{
			name: "two, false",
			results: result.Results{
				Results: []result.Result{
					{
						Result: false,
					},
					{
						Result: false,
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.results.AllTrue() != test.expected {
				t.Errorf("allTrue got: %t, expected: %t", test.results.AllTrue(), test.expected)
			}
		})
	}
}

func TestResults_SomeTrue(t *testing.T) {
	tests := []struct {
		name     string
		results  result.Results
		expected bool
	}{
		{
			name: "one, true",
			results: result.Results{
				Results: []result.Result{
					{
						Result: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "two, true",
			results: result.Results{
				Results: []result.Result{
					{
						Result: true,
					},
					{
						Result: true,
					},
				},
			},
			expected: true,
		},
		{
			name: "one true, one false",
			results: result.Results{
				Results: []result.Result{
					{
						Result: true,
					},
					{
						Result: false,
					},
				},
			},
			expected: true,
		},
		{
			name: "two, false",
			results: result.Results{
				Results: []result.Result{
					{
						Result: false,
					},
					{
						Result: false,
					},
				},
			},
			expected: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if test.results.SomeTrue() != test.expected {
				t.Errorf("allTrue got: %t, expected: %t", test.results.AllTrue(), test.expected)
			}
		})
	}
}