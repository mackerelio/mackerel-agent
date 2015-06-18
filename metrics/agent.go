package metrics

import (
	"runtime"
)

// AgentGenerator is generator of metrics
// about the runnning agent itself
type AgentGenerator struct {
}

var memStats = new(runtime.MemStats)

// Generate generates the memory usage of the running agent itself
func (g *AgentGenerator) Generate() (Values, error) {
	runtime.ReadMemStats(memStats)

	ret := map[string]float64{
		"custom.agent.memory.alloc":     float64(memStats.Alloc),
		"custom.agent.memory.sys":       float64(memStats.Sys),
		"custom.agent.memory.heapAlloc": float64(memStats.HeapAlloc),
		"custom.agent.memory.heapSys":   float64(memStats.HeapSys),
	}

	return ret, nil
}
