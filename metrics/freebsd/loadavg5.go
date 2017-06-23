// +build freebsd

package freebsd

import (
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// Loadavg5Generator XXX
type Loadavg5Generator struct {
}

var loadavg5Logger = logging.GetLogger("metrics.loadavg5")

// Generate generate metric values
// % sysctl -n vm.loadavg
// { 2.26 2.08 2.00 }
func (g *Loadavg5Generator) Generate() (metrics.Values, error) {
	outputBytes, err := exec.Command("sysctl", "-n", "vm.loadavg").Output()
	if err != nil {
		loadavg5Logger.Errorf("Failed to run sysctl -n vm.loadavg: %s", err)
		return nil, err
	}

	output := string(outputBytes)

	// fields will be "{", <loadavg1>, <loadavg5>, <loadavg15>, "}"
	fields := strings.Fields(output)

	if len(fields) != 5 || fields[0] != "{" || fields[len(fields)-1] != "}" {
		loadavg5Logger.Errorf("sysctl -n vm.loadavg result malformed: %s", output)
		return nil, err
	}

	loadavg5, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		loadavg5Logger.Errorf("Failed to parse loadavg5 string: %s", err)
		return nil, err
	}

	return metrics.Values{"loadavg5": loadavg5}, nil
}
