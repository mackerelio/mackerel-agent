// +build !windows

package metrics

import (
	"github.com/mackerelio/go-osstat/loadavg"
	"github.com/mackerelio/golib/logging"
)

// loadavg
//   - loadavg1: load average per 1 minutes
//   - loadavg5: load average per 5 minutes
//   - loadavg15: load average per 15 minutes

// LoadavgGenerator generates load average values
type LoadavgGenerator struct {
}

var loadavgLogger = logging.GetLogger("metrics.loadavg")

// Generate load averages
func (g *LoadavgGenerator) Generate() (Values, error) {
	loadavgs, err := loadavg.Get()
	if err != nil {
		loadavgLogger.Errorf("%s", err)
		return nil, err
	}
	return Values{
		"loadavg1":  loadavgs.Loadavg1,
		"loadavg5":  loadavgs.Loadavg5,
		"loadavg15": loadavgs.Loadavg15,
	}, nil
}
