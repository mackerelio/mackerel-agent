//go:build linux

package linux

import (
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
MemoryGenerator collect memory usage

`memory.{metric}`: using memory size[KiB] retrieved from /proc/meminfo

metric = "total", "free", "buffers", "cached", "swap_cached", "swap_total", "swap_free"

Metrics "used" is calculated by

	(total - available) when MemAvailable is available in /proc/meminfo and
	(total - free - buffers - cached) otherwise

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

	ret := metrics.Values{
		"memory.total":       metrics.NewValueAttribute(float64(mem.Total)),
		"memory.used":        metrics.NewValueAttribute(float64(mem.Used)),
		"memory.swap_total":  metrics.NewValueAttribute(float64(mem.SwapTotal)),
		"memory.swap_cached": metrics.NewValueAttribute(float64(mem.SwapCached)),
		"memory.swap_free":   metrics.NewValueAttribute(float64(mem.SwapFree)),
	}

	if mem.MemAvailableEnabled {
		ret["memory.mem_available"] = metrics.NewValueAttribute(float64(mem.Available))
	} else {
		ret["memory.buffers"] = metrics.NewValueAttribute(float64(mem.Buffers))
		ret["memory.cached"] = metrics.NewValueAttribute(float64(mem.Cached))
		ret["memory.free"] = metrics.NewValueAttribute(float64(mem.Free))
	}

	return ret, nil
}
