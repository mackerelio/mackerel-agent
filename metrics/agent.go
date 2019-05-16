package metrics

import (
	"runtime"

	mkr "github.com/mackerelio/mackerel-client-go"
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
		"custom.agent.memory.alloc":          float64(memStats.Alloc),
		"custom.agent.memory.sys":            float64(memStats.Sys),
		"custom.agent.memory.heapAlloc":      float64(memStats.HeapAlloc),
		"custom.agent.memory.heapSys":        float64(memStats.HeapSys),
		"custom.agent.runtime.goroutine_num": float64(runtime.NumGoroutine()),
	}

	return ret, nil
}

// CustomIdentifier for PluginGenerator interface
func (g *AgentGenerator) CustomIdentifier() *string {
	return nil
}

// PrepareGraphDefs for PluginGenerator interface
func (g *AgentGenerator) PrepareGraphDefs() ([]*mkr.GraphDefsParam, error) {
	meta := &pluginMeta{
		Graphs: map[string]customGraphDef{
			"agent.memory": customGraphDef{
				Label: "Agent Memory",
				Unit:  "bytes",
				Metrics: []customGraphMetricDef{
					{Name: "alloc", Label: "Alloc"},
					{Name: "sys", Label: "Sys"},
					{Name: "heapAlloc", Label: "Heap Alloc"},
					{Name: "heapSys", Label: "Heap Sys"},
				},
			},
			"agent.runtime": customGraphDef{
				Label: "Agent Runtime",
				Unit:  "integer",
				Metrics: []customGraphMetricDef{
					{Name: "goroutine_num", Label: "Goroutine Num"},
				},
			},
		},
	}
	return makeGraphDefsParam(meta), nil
}
