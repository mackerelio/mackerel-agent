// +build darwin

package darwin

import "testing"

func TestCPUUsageGenerator_Generate(t *testing.T) {
	g := &CPUUsageGenerator{}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should not have occurred: %s", err)
	}

	metricName := []string{"cpu.user.percentage", "cpu.system.percentage", "cpu.idle.percentage"}
	for _, n := range metricName {
		if _, ok := values[n]; !ok {
			t.Errorf("should have '%s': %v", n, values)
		}
	}
}
