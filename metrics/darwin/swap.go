// +build darwin

package darwin

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
SwapGenerator collect swap usage

`memory.{metric}`: using swap size retrieved from `sysctl vm.swapusage`

metric = "swap_total", "swap_free"

graph: `memory.{metric}`
*/
type SwapGenerator struct {
}

/* sysctl vm.swapusage sample
% sysctl vm.swapusage
vm.swapusage: total = 1024.00M  used = 2.50M  free = 1021.50M  (encrypted)
*/

var swapLogger = logging.GetLogger("metrics.memory.swap")
var swapReg = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)M[^0-9]*([0-9]+(?:\.[0-9]+)?)M[^0-9]*([0-9]+(?:\.[0-9]+)?)M`)

// Generate generate swap values
func (g *SwapGenerator) Generate() (metrics.Values, error) {
	outBytes, err := exec.Command("sysctl", "vm.swapusage").Output()
	if err != nil {
		swapLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}

	out := string(outBytes)
	matches := swapReg.FindStringSubmatch(out)
	if matches == nil || len(matches) != 4 {
		return nil, fmt.Errorf("faild to parse vm.swapusage result: [%q]", out)
	}
	t, _ := strconv.ParseFloat(matches[1], 64)
	// swap_used are calculated at server, so don't send it
	// u, _ := strconv.ParseFloat(matches[2], 64)
	f, _ := strconv.ParseFloat(matches[3], 64)

	const mb = 1024.0 * 1024.0
	return metrics.Values{
		"memory.swap_total": t * mb,
		"memory.swap_free":  f * mb,
	}, nil
}
