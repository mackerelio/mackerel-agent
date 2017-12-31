// +build !windows

package metrics

import (
	"testing"
)

func TestLoadAvg5Generate(t *testing.T) {
	g := &Loadavg5Generator{}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should be nil but got: %s", err)
	}

	metricName := []string{"loadavg5"}
	for _, n := range metricName {
		if _, ok := values[n]; !ok {
			t.Errorf("loadavg5 metrics should have '%s': %v", n, values)
		}
	}

	t.Logf("loadavg5 metrics: %+v", values)
}
