// +build linux

package linux

import (
	"math"
)
import (
	"testing"
	"time"
)

func TestCPUUsageGenerate(t *testing.T) {
	g := &CPUUsageGenerator{1 * time.Second}
	values, _ := g.Generate()

	sumPercentage := float64(0)
	for _, metricName := range cpuUsageMetricNames {
		metricName += ".percentage"
		value, ok := values[metricName]
		if !ok {
			t.Errorf("CPUUsageGenerator should generate metric value for '%s'", metricName)
		} else {
			t.Logf("CPUUsage '%s' collected: %+v", metricName, value)
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

	// Checks any errors will not occure
	// when the number of retrieved values from /proc/spec is less than length of cpuUsageMetricNames
	cpuUsageMetricNames = append(cpuUsageMetricNames, "unimplemented-new-metric")
	defer func() { cpuUsageMetricNames = cpuUsageMetricNames[0 : len(cpuUsageMetricNames)-1] }()

	g.Generate()
}
