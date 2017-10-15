// +build darwin

package darwin

import (
	"testing"
	"time"
)

func TestCPUUsageGenerator_Generate(t *testing.T) {
	g := &CPUUsageGenerator{1 * time.Second}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should not have occurred: %s", err)
	}

	metricNames := []string{"cpu.user.percentage", "cpu.system.percentage", "cpu.idle.percentage"}

	for _, name := range metricNames {
		if v, ok := values[name]; !ok {
			t.Errorf("cpu should has %s", name)
		} else {
			t.Logf("cpu '%s' collected: %+v", name, v)
		}
	}
}
