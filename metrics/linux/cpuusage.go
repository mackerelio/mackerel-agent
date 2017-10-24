// +build linux

package linux

import (
	"time"

	"github.com/mackerelio/go-osstat/cpu"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
collect CPU usage

`cpu.{metric}.percentage`: The increased amount of CPU time per minute as percentage of total CPU cores x 100

metric = "user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest"

graph: stacks `cpu.{metric}.percentage`
*/

// CPUUsageGenerator generates CPU metric values
type CPUUsageGenerator struct {
	Interval time.Duration
}

var cpuUsageLogger = logging.GetLogger("metrics.cpuUsage")

// Generate CPU metric values
func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {
	previous, err := g.collectProcStatValues()
	if err != nil {
		return nil, err
	}

	time.Sleep(g.Interval)

	current, err := g.collectProcStatValues()
	if err != nil {
		return nil, err
	}

	totalDiff := float64(current.Total - previous.Total)
	cpuCount := float64(current.CPUCount)

	ret := map[string]float64{
		"cpu.user.percentage":    float64((current.User-current.Guest)-(previous.User-previous.Guest)) * cpuCount * 100.0 / totalDiff,
		"cpu.nice.percentage":    float64(current.Nice-previous.Nice) * cpuCount * 100.0 / totalDiff,
		"cpu.system.percentage":  float64(current.System-previous.System) * cpuCount * 100.0 / totalDiff,
		"cpu.idle.percentage":    float64(current.Idle-previous.Idle) * cpuCount * 100.0 / totalDiff,
		"cpu.iowait.percentage":  float64(current.Iowait-previous.Iowait) * cpuCount * 100.0 / totalDiff,
		"cpu.irq.percentage":     float64(current.Irq-previous.Irq) * cpuCount * 100.0 / totalDiff,
		"cpu.softirq.percentage": float64(current.Softirq-previous.Softirq) * cpuCount * 100.0 / totalDiff,
		"cpu.steal.percentage":   float64(current.Steal-previous.Steal) * cpuCount * 100.0 / totalDiff,
		"cpu.guest.percentage":   float64(current.Guest-previous.Guest) * cpuCount * 100.0 / totalDiff,
		// guest_nice is not yet supported in Mackerel
		// "cpu.guest_nice.percentage": float64(current.GuestNice - previous.GuestNice) * cpuCount * 100.0 / totalDiff,
	}
	return metrics.Values(ret), nil
}

// returns values corresponding to cpuUsageMetricNames, those total and the number of CPUs
func (g *CPUUsageGenerator) collectProcStatValues() (*cpu.Stats, error) {
	cpu, err := cpu.Get()
	if err != nil {
		cpuUsageLogger.Errorf("failed to get cpu statistics: %s", err)
		return nil, err
	}
	return cpu, nil
}
