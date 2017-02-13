package metadata

import (
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestMetadataGeneratorInterval(t *testing.T) {
	tests := []struct {
		interval *int32
		expected time.Duration
	}{
		{
			interval: pint(0),
			expected: 1 * time.Minute,
		},
		{
			interval: pint(1),
			expected: 1 * time.Minute,
		},
		{
			interval: pint(30),
			expected: 30 * time.Minute,
		},
		{
			interval: nil,
			expected: 10 * time.Minute,
		},
	}
	for _, test := range tests {
		g := Generator{
			Config: &config.MetadataPlugin{
				ExecutionInterval: test.interval,
			},
		}
		if g.Interval() != test.expected {
			t.Errorf("interval should be %v but got: %v", test.expected, g.Interval())
		}
	}
}

func pint(i int32) *int32 {
	return &i
}
