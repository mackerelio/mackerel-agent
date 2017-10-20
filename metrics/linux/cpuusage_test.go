// +build linux

package linux

import (
	"math"
	"testing"
	"time"
)

func TestCPUUsageGenerate(t *testing.T) {
	g := &CPUUsageGenerator{1 * time.Second}
	values, _ := g.Generate()

	var metricNames = []string{
		"user", "nice", "system", "idle", "iowait",
		"irq", "softirq", "steal", "guest",
	}

	sumPercentage := float64(0)
	for _, name := range metricNames {
		metricName := "cpu." + name + ".percentage"
		value, ok := values[metricName]
		if !ok {
			t.Errorf("cpu values shuold have '%s': %v", metricName, values)
		}
		sumPercentage += value
	}

	percentDistFrom100 := math.Mod(sumPercentage, 100)
	// Checks sum of each persentages for cores are not so far from 100%
	if math.Min(percentDistFrom100, 100-percentDistFrom100) > 10 {
		t.Errorf("Sum of CPU usage percentage values art not N * 100%%: %f", sumPercentage)
	} else {
		t.Logf("Sum of CPU usage percentage values: %f", sumPercentage)
	}

	t.Logf("cpu metric metrics: %+v", values)
}
