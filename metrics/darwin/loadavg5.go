// +build darwin

package darwin

import (
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// Loadavg5Generator XXX
type Loadavg5Generator struct {
}

var loadavg5Logger = logging.GetLogger("metrics.loadavg5")

// Generate load averages
func (g *Loadavg5Generator) Generate() (metrics.Values, error) {
	loadavgs, err := loadavg.Get()
	if err != nil {
		loadavg5Logger.Errorf("%s", err)
		return nil, err
	}
	return metrics.Values{"loadavg5": loadavgs.Loadavg5}, nil
}
