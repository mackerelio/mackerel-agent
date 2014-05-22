// +build windows

package windows

import "math"
import "testing"

var cpuusageMetricNames = []string{
	//"cpu.user",
	//"cpu.nice",
	//"cpu.system",
	//"cpu.idle",
	//"cpu.iowait",
	//"cpu.irq",
	//"cpu.softirq",
	//"cpu.steal",
	//"cpu.guest",
}

func TestCpuusageGenerate(t *testing.T) {
	g := &CpuusageGenerator{1}
	values, _ := g.Generate()

	for _, metricName := range cpuusageMetricNames {
		value, ok := values[metricName]
		if !ok {
			t.Errorf("CpuusageGenerator should generate metric value for '%s'", metricName)
		} else {
			t.Logf("Cpuusage '%s' collected: %+v", metricName, value)
		}
	}

	var sumPercentage float64 = 0
	for _, metricName := range cpuusageMetricNames {
		metricName += ".percentage"
		value, ok := values[metricName]
		if !ok {
			t.Errorf("CpuusageGenerator should generate metric value for '%s'", metricName)
		} else {
			t.Logf("Cpuusage '%s' collected: %+v", metricName, value)
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
	// when the number of retrieved values from /proc/spec is less than length of cpuusageMetricNames
	cpuusageMetricNames = append(cpuusageMetricNames, "unimplemented-new-metric")
	defer func() { cpuusageMetricNames = cpuusageMetricNames[0 : len(cpuusageMetricNames)-1] }()

	g.Generate()
}
