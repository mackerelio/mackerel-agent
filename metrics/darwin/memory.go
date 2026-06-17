//go:build darwin

package darwin

import (
	"github.com/mackerelio/go-osstat/memory"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
MemoryGenerator collect memory usage

`memory.{metric}`: using memory size retrieved from `vm_stat`

metric = "total", "free", "used", "cached"

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
	return metrics.Values{
		"memory.total":      metrics.NewValueAttribute(float64(memory.Total)),
		"memory.used":       metrics.NewValueAttribute(float64(memory.Used)),
		"memory.cached":     metrics.NewValueAttribute(float64(memory.Cached)),
		"memory.free":       metrics.NewValueAttribute(float64(memory.Free)),
		"memory.swap_total": metrics.NewValueAttribute(float64(memory.SwapTotal)),
		"memory.swap_free":  metrics.NewValueAttribute(float64(memory.SwapFree)),
	}, nil
}
