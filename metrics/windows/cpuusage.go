// +build windows

package windows

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util"
)

// CPUUsageGenerator XXX
type CPUUsageGenerator struct {
}

var cpuUsageLogger = logging.GetLogger("metrics.cpuusage")

// NewCPUUsageGenerator XXX
func NewCPUUsageGenerator() (*CPUUsageGenerator, error) {
	return &CPUUsageGenerator{}, nil
}

// Generate XXX
func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {
	cpuusage := make(map[string]float64, 1)

	cpuusageValue, err := util.GetWmicToFloat("cpu", "loadpercentage")
	if err != nil {
		cpuusageValue = 0
	}
	cpuusage["cpu.user.percentage"] = cpuusageValue
	cpuusage["cpu.idle.percentage"] = 100 - cpuusageValue
	cpuUsageLogger.Debugf("cpuusage : %s", cpuusage)
	return metrics.Values(cpuusage), nil
}
