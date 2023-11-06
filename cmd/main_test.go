package main

import (
	"testing"
)

func TestSliceIntersect(t *testing.T) {
	testCases := []struct {
		no       int
		version  string
		expected string
	}{
		{
			no:       1,
			version:  "const AutoVersion = \"0.1.0\"",
			expected: "const AutoVersion = \"0.1.1\"",
		},
		{
			no:       2,
			version:  "const AutoVersion = \"0.1.1\"",
			expected: "const AutoVersion = \"0.1.2\"",
		},
		{
			no:       3,
			version:  "const AutoVersion = \"0.1.120\"",
			expected: "const AutoVersion = \"0.1.121\"",
		},
		{
			no:       4,
			version:  "const AutoVersion = \"0.1.075\"",
			expected: "const AutoVersion = \"0.1.76\"",
		},
		{
			no:       5,
			version:  "const AutoVersion = \"99.0.3\"",
			expected: "const AutoVersion = \"99.0.4\"",
		},
		{
			no:       6,
			version:  "const AutoVersion = \"25\"",
			expected: "const AutoVersion = \"26\"",
		},
		{
			no:       7,
			version:  "const AutoVersion = \"0\"",
			expected: "const AutoVersion = \"1\"",
		},
		{
			no:       8,
			version:  "const AutoVersion = \"2.7\"",
			expected: "const AutoVersion = \"2.8\"",
		},
		{
			no:       9,
			version:  "const AutoVersion = \"1.11.1 test version\"",
			expected: "const AutoVersion = \"1.11.2 test version\"",
		},
		{
			no:       10,
			version:  "const AutoVersion = \"435345.23452.98288\"",
			expected: "const AutoVersion = \"435345.23452.98289\"",
		},
		{
			no:       11,
			version:  "const AutoVersion = \"0.999.3233\"",
			expected: "const AutoVersion = \"0.999.3234\"",
		},
		{
			no:       12,
			version:  "const AutoVersion = \"0.999.3233 Dirty Version, do NOT ship to customers (I MEAN IT) 0.5!!\"",
			expected: "const AutoVersion = \"0.999.3233 Dirty Version, do NOT ship to customers (I MEAN IT) 0.6!!\"",
		},
		{
			no:       13,
			version:  "const AutoVersion = \"We hate numbers!\"",
			expected: "const AutoVersion = \"We hate numbers!\"",
		},
		{
			no:       14,
			version:  "const AutoVersion = \"\"",
			expected: "const AutoVersion = \"\"",
		},
	}
	//logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	for _, tc := range testCases {
		//result := increaseVersionNumber(tc.version, logger)
		result := increaseVersionNumber(tc.version)
		if result != tc.expected {
			t.Errorf("test %d failed, got:%s expected:%s\n", tc.no, result, tc.expected)
		}
	}
}
