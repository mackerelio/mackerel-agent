// +build windows

package windows

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"fmt"
	"strings"
	"os/exec"
	"strconv"
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

	cpuGet, err := exec.Command("wmic", "cpu", "get", "loadpercentage").Output()
	if err != nil {
		cpuUsageLogger.Errorf("Failed to invoke 'wmic': %s", err)
		return nil, err
	}

	percentages := string(cpuGet)

	lines := strings.Split(percentages, "\r\r\n")

	if len(lines) <= 2 {
		return nil, fmt.Errorf("wmic result malformed: [%q]", lines)
	}

	cpuusage := make(map[string]float64, 1)

	value, err := strconv.ParseFloat(strings.Trim(lines[1], " "), 64)
	if err != nil {
		value = 0
	}

	cpuusage["cpu.user.percentage"] = value
	cpuusage["cpu.idle.percentage"] = 100 - value
	cpuUsageLogger.Debugf("cpuusage : %s", cpuusage)
	return metrics.Values(cpuusage), nil
}
