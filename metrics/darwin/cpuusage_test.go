// +build darwin

package darwin

import "testing"

func TestCpuusageGenerator_Generate(t *testing.T) {
	g := &CpuusageGenerator{}
	values, err := g.Generate()

	if err != nil {
		t.Error("error should not have occured: %s", err)
	}

	metricName := []string{"cpu.user.percentage", "cpu.system.percentage", "cpu.idle.percentage"}
	for _, n := range metricName {
		if _, ok := values[n]; !ok {
			t.Errorf("should have '%s': %v", n, values)
		}
	}
}
