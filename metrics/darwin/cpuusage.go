// +build darwin

package darwin

import (
	"errors"
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// CPUUsageGenerator XXX
type CPUUsageGenerator struct {
	Interval time.Duration
}

var cpuUsageLogger = logging.GetLogger("metrics.cpuUsage")

// Generate returns current CPU usage of the host.
// Keys below are expected:
// - cpu.user.percentage
// - cpu.system.percentage
// - cpu.idle.percentage
func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {
	before, err := cpu.Get()
	if err != nil {
		cpuUsageLogger.Errorf("failed to get cpu statistics: %s", err)
		return nil, err
	}

	time.Sleep(g.Interval)

	after, err := cpu.Get()
	if err != nil {
		cpuUsageLogger.Errorf("failed to get cpu statistics: %s", err)
		return nil, err
	}

	if before.Total == after.Total {
		err := errors.New("cpu total counter did not change")
		cpuUsageLogger.Errorf("%s", err)
		return nil, err
	}

	cpuUsage := make(map[string]float64, 3)
	total := float64(after.Total - before.Total)
	cpuUsage["cpu.user.percentage"] = float64(after.User-before.User) / total * 100.0
	cpuUsage["cpu.system.percentage"] = float64(after.System-before.System) / total * 100.0
	cpuUsage["cpu.idle.percentage"] = float64(after.Idle-before.Idle) / total * 100.0
	return metrics.Values(cpuUsage), nil
}
