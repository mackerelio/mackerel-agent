// +build windows

package windows

import (
	"errors"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// CPUUsageGenerator XXX
type CPUUsageGenerator struct {
	Interval time.Duration
}

var cpuUsageLogger = logging.GetLogger("metrics.cpuUsage")

func NewCPUUsageGenerator(interval time.Duration) (*CPUUsageGenerator, error) {
	return &CPUUsageGenerator{interval}, nil
}

func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {
	if g == nil {
		return nil, errors.New("CPUUsageGenerator is not initialized")
	}
	time.Sleep(g.Interval * time.Second)

	// TODO
	return metrics.Values{}, nil
}
