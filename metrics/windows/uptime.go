// +build windows

package windows

import (
	"errors"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	. "github.com/mackerelio/mackerel-agent/util/windows"
)

type UptimeGenerator struct {
}

var uptimeLogger = logging.GetLogger("metrics.uptime")

func NewUptimeGenerator() *UptimeGenerator {
	return &UptimeGenerator{}
}

func (g *UptimeGenerator) Generate() (metrics.Values, error) {
	if g == nil {
		return nil, errors.New("UptimeGenerator is not initialized")
	}
	r, _, err := GetTickCount.Call()
	if r == 0 {
		return nil, err
	}

	return metrics.Values(map[string]float64{"uptime": float64(r) / 1000}), nil
}
