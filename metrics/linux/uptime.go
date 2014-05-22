// +build linux

package linux

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
collect uptime

`uptime`: uptime[day] retrieved from /proc/uptime

graph: `uptime`
*/
type UptimeGenerator struct {
}

var uptimeLogger = logging.GetLogger("metrics.uptime")

func (g *UptimeGenerator) Generate() (metrics.Values, error) {
	contentbytes, err := ioutil.ReadFile("/proc/uptime")
	if err != nil {
		uptimeLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}
	content := string(contentbytes)
	cols := strings.Split(content, " ")

	f, err := strconv.ParseFloat(cols[0], 64)
	if err != nil {
		uptimeLogger.Errorf("Failed to parse values (skip these metrics): %s", err)
		return nil, err
	}

	return metrics.Values(map[string]float64{"uptime": f / 86400}), nil
}
