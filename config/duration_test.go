package config

import (
	"fmt"
	"testing"

	"github.com/BurntSushi/toml"
)

func TestParseDuration(t *testing.T) {
	testCases := []struct {
		src      string
		expected int32
		err      string
	}{
		{
			src:      "0",
			expected: 0,
		},
		{
			src:      "10",
			expected: 10,
		},
		{
			src: "-10",
			err: "duration out of range: -10",
		},
		{
			src:      "10m",
			expected: 10,
		},
		{
			src:      "1h10m",
			expected: 70,
		},
		{
			src:      "2.5h",
			expected: 150,
		},
		{
			src: "-10m",
			err: "duration out of range: -10m0s",
		},
		{
			src: "1s",
			err: "duration not multiple of 1m: 1s",
		},
		{
			src: "1ms",
			err: "duration not multiple of 1m: 1ms",
		},
		{
			src: "1.1m",
			err: "duration not multiple of 1m: 1m6s",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.src, func(t *testing.T) {
			var m struct{ Duration *duration }
			if _, err := toml.Decode(fmt.Sprintf(`duration = %q`, tc.src), &m); err != nil {
				if tc.err == "" {
					t.Fatalf("duration %q, got error: %v", tc.src, err)
				} else if err.Error() != tc.err {
					t.Fatalf("duration %q, expected error: %v, got error: %v", tc.src, tc.err, err)
				}
			}
			got := m.Duration.Minutes()
			if *got != tc.expected {
				t.Errorf("duration %q, expected: %v, got: %v", tc.src, tc.expected, *got)
			}
		})
	}
}
