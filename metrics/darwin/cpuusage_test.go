// +build darwin

package darwin

import (
	"reflect"
	"testing"

	"github.com/hatena/mackerel-container-agent/metric"
)

func TestCPUUsageGenerator_Generate(t *testing.T) {
	g := &CPUUsageGenerator{}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should be nil but got: %s", err)
	}

	metricName := []string{"cpu.user.percentage", "cpu.system.percentage", "cpu.idle.percentage"}
	for _, n := range metricName {
		if _, ok := values[n]; !ok {
			t.Errorf("should have '%s': %v", n, values)
		}
	}

	t.Logf("cpu metrics: %+v", values)
}

func TestCPUUsageGenerator_parseIostatOutput(t *testing.T) {
	testCases := []struct {
		output string
		values metric.Values
	}{
		{
			output: `      cpu    load average
 us sy id   1m   5m   15m
 19  9 72  2.50 3.04 3.20
 16 12 72  2.50 3.04 3.20
`,
			values: metric.Values{
				"cpu.user.percentage":   16.0,
				"cpu.system.percentage": 12.0,
				"cpu.idle.percentage":   72.0,
			},
		},
		{
			output: `      cpu    load average
 us sy id   1m   5m   15m
 19  9 72  2.50 3.04 3.20
      cpu    load average
 us sy id   1m   5m   15m
 16 12 72  2.50 3.04 3.20
`,
			values: metric.Values{
				"cpu.user.percentage":   16.0,
				"cpu.system.percentage": 12.0,
				"cpu.idle.percentage":   72.0,
			},
		},
	}
	for _, testCase := range testCases {
		got, err := parseIostatOutput(testCase.output)
		if err != nil {
			t.Errorf("error should be nil but got: %s", err)
		}
		if !reflect.DeepEqual(map[string]float64(got), map[string]float64(testCase.values)) {
			t.Errorf("metric values should be %#v but got: %#v", testCase.values, got)
		}
	}
}
