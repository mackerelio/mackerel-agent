// +build !windows

package metrics

import (
	"testing"
)

func TestLoadAvgGenerate(t *testing.T) {
	g := &LoadavgGenerator{}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should be nil but got: %s", err)
	}

	metricName := []string{"loadavg1", "loadavg5", "loadavg15"}
	for _, n := range metricName {
		if _, ok := values[n]; !ok {
			t.Errorf("loadavg metrics should have '%s': %v", n, values)
		}
	}

	t.Logf("loadavg metrics: %+v", values)
}
