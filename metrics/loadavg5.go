package metrics

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

/*
collect load average

`loadavg5`: load average per 5 minutes retrieved from /proc/loadavg

graph: `loadavg5`
*/
type Loadavg5Generator struct {
}

var loadavg5Logger = logging.GetLogger("metrics.loadavg5")

func (g *Loadavg5Generator) Generate() (Values, error) {
	contentbytes, err := ioutil.ReadFile("/proc/loadavg")
	if err != nil {
		loadavg5Logger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}
	content := string(contentbytes)
	cols := strings.Split(content, " ")

	f, err := strconv.ParseFloat(cols[1], 64)
	if err != nil {
		loadavg5Logger.Errorf("Failed to parse loadavg5 metrics (skip these metrics): %s", err)
		return nil, err
	}

	return Values(map[string]float64{"loadavg5": f}), nil
}
