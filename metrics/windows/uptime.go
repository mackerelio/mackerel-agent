// +build windows

package windows

import (
	"errors"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// UptimeGenerator XXX
type UptimeGenerator struct {
}

var uptimeLogger = logging.GetLogger("metrics.uptime")

// NewUptimeGenerator XXX
func NewUptimeGenerator() (*UptimeGenerator, error) {
	return &UptimeGenerator{}, nil
}

// Generate XXX
func (g *UptimeGenerator) Generate() (metrics.Values, error) {
	if g == nil {
		return nil, errors.New("UptimeGenerator is not initialized")
	}
	r, _, err := windows.GetTickCount.Call()
	if r == 0 {
		return nil, err
	}

	return metrics.Values(map[string]float64{"uptime": float64(r) / 1000}), nil
}
