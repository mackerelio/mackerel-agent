// +build windows

package windows

import (
	"errors"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

type CpuusageGenerator struct {
	Interval time.Duration
}

var cpuusageLogger = logging.GetLogger("metrics.cpuusage")

func NewCpuusageGenerator(interval time.Duration) *CpuusageGenerator {
	return &CpuusageGenerator{interval}
}

func (g *CpuusageGenerator) Generate() (metrics.Values, error) {
	if g == nil {
		return nil, errors.New("CpuusageGenerator is not initialized")
	}
	time.Sleep(g.Interval * time.Second)

	// TODO
	return metrics.Values{}, nil
}
