// +build linux

package linux

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
MemoryGenerator collect memory usage

`memory.{metric}`: using memory size[KiB] retrieved from /proc/meminfo

metric = "total", "free", "buffers", "cached", "active", "inactive", "swap_cached", "swap_total", "swap_free"

Metrics "used" is caluculated here like (total - free - buffers - cached) for ease.
This calculation may be going to be done in server side in the future.

graph: stacks `memory.{metric}`
*/
type MemoryGenerator struct {
}

var memoryLogger = logging.GetLogger("metrics.memory")

// Generate generate metrics values
func (g *MemoryGenerator) Generate() (metrics.Values, error) {
	out, err := ioutil.ReadFile("/proc/meminfo")
	if err != nil {
		memoryLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}
	return parseMeminfo(out)
}

var memItems = map[string]string{
	"MemTotal":     "total",
	"MemFree":      "free",
	"MemAvailable": "available",
	"Buffers":      "buffers",
	"Cached":       "cached",
	"Active":       "active",
	"Inactive":     "inactive",
	"SwapCached":   "swap_cached",
	"SwapTotal":    "swap_total",
	"SwapFree":     "swap_free",
}

func parseMeminfo(out []byte) (metrics.Values, error) {
	scanner := bufio.NewScanner(bytes.NewReader(out))

	ret := make(map[string]float64)
	var total, unused, available float64
	usedCnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		// ex.) MemTotal:        3916792 kB
		kvAndUnit := strings.Fields(line)
		if len(kvAndUnit) < 3 {
			continue
		}
		name := strings.TrimRight(kvAndUnit[0], ":")
		k, ok := memItems[name]
		if !ok {
			continue
		}
		if kvAndUnit[2] != "kB" {
			memoryLogger.Warningf("/proc/meminfo contains an invalid unit: %s", k)
			continue
		}
		value, err := strconv.ParseFloat(kvAndUnit[1], 64)
		if err != nil {
			memoryLogger.Warningf("Failed to parse memory metrics: %s", err)
			continue
		}
		ret["memory."+k] = value * 1024
		switch k {
		case "free", "buffers", "cached":
			unused += value
			usedCnt++
		case "total":
			total = value
			usedCnt++
		case "available":
			available = value
		}
	}
	if err := scanner.Err(); err != nil {
		memoryLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}
	if total > 0 && available > 0 {
		ret["memory.used"] = (total - available) * 1024
	} else if usedCnt == 4 { // 4 is free, buffers, cached and total
		ret["memory.used"] = (total - unused) * 1024
	}

	return metrics.Values(ret), nil
}
