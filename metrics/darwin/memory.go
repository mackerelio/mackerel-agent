// +build darwin

package darwin

import (
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
MemoryGenerator collect memory usage

`memory.{metric}`: using memory size retrieved from `vm_stat`

metric = "total", "free", "used", "cached", "active", "inactive"

graph: stacks `memory.{metric}`
*/
type MemoryGenerator struct {
}

var memoryLogger = logging.GetLogger("metrics.memory")

// Generate generate metrics values
func (g *MemoryGenerator) Generate() (metrics.Values, error) {
	memory, err := memory.Get()
	if err != nil {
		memoryLogger.Errorf("failed to get memory statistics: %s", err)
		return nil, err
	}
	ret := map[string]float64{
		"memory.total":      float64(memory.Total),
		"memory.used":       float64(memory.Used),
		"memory.cached":     float64(memory.Cached),
		"memory.free":       float64(memory.Free),
		"memory.active":     float64(memory.Active),
		"memory.inactive":   float64(memory.Inactive),
		"memory.swap_total": float64(memory.SwapTotal),
		"memory.swap_free":  float64(memory.SwapFree),
	}
	return metrics.Values(ret), nil
}
