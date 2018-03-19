// +build linux

package linux

import (
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/golib/logging"
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

// Generate memory values
func (g *MemoryGenerator) Generate() (metrics.Values, error) {
	mem, err := memory.Get()
	if err != nil {
		memoryLogger.Errorf("failed to get memory statistics: %s", err)
		return nil, err
	}

	ret := map[string]float64{
		"memory.total":       float64(mem.Total),
		"memory.used":        float64(mem.Used),
		"memory.active":      float64(mem.Active),
		"memory.inactive":    float64(mem.Inactive),
		"memory.swap_total":  float64(mem.SwapTotal),
		"memory.swap_cached": float64(mem.SwapCached),
		"memory.swap_free":   float64(mem.SwapFree),
	}

	if mem.MemAvailableEnabled {
		ret["memory.mem_available"] = float64(mem.Available)
	} else {
		ret["memory.buffers"] = float64(mem.Buffers)
		ret["memory.cached"] = float64(mem.Cached)
		ret["memory.free"] = float64(mem.Free)
	}

	return metrics.Values(ret), nil
}
