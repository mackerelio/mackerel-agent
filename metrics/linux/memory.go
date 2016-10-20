// +build linux

package linux

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"regexp"
	"strconv"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
MemoryGenerator collect memory usage

`memory.{metric}`: using memory size[KiB] retrieved from /proc/meminfo

metric = "total", "free", "buffers", "cached", "active", "inactive", "swap_cached", "swap_total", "swap_free"

Metrics "used" is calculated here like (total - free - buffers - cached) for ease.
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

var memReg = regexp.MustCompile(`^([A-Za-z]+):\s+([0-9]+)\s+kB`)

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

	ret := metrics.Values{}
	used := float64(0)
	usedCnt := 0
	for scanner.Scan() {
		line := scanner.Text()
		// ex.) MemTotal:        3916792 kB
		if matches := memReg.FindStringSubmatch(line); len(matches) == 3 {
			k, ok := memItems[matches[1]]
			if !ok {
				continue
			}
			value, _ := strconv.ParseFloat(matches[2], 64)
			ret["memory."+k] = value * 1024
			switch k {
			case "free", "buffers", "cached":
				used -= value
				usedCnt++
			case "total":
				used += value
				usedCnt++
			}
		}
	}
	if err := scanner.Err(); err != nil {
		memoryLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}
	if usedCnt == 4 { // 4 is free, buffers, cached and total
		ret["memory.used"] = used * 1024
	}

	return ret, nil
}
