package metrics

import (
	"runtime"

	mkr "github.com/mackerelio/mackerel-client-go"
)

// AgentGenerator is generator of metrics
// about the running agent itself
type AgentGenerator struct {
}

var memStats = new(runtime.MemStats)

// Generate generates the memory usage of the running agent itself
func (g *AgentGenerator) Generate() (Values, error) {
	runtime.ReadMemStats(memStats)

	return Values{
		"custom.agent.memory.alloc":          NewValueAttribute(float64(memStats.Alloc)),
		"custom.agent.memory.sys":            NewValueAttribute(float64(memStats.Sys)),
		"custom.agent.memory.heapAlloc":      NewValueAttribute(float64(memStats.HeapAlloc)),
		"custom.agent.memory.heapSys":        NewValueAttribute(float64(memStats.HeapSys)),
		"custom.agent.runtime.goroutine_num": NewValueAttribute(float64(runtime.NumGoroutine())),
	}, nil
}

// CustomIdentifier for PluginGenerator interface
func (g *AgentGenerator) CustomIdentifier() *string {
	return nil
}

// PrepareGraphDefs for PluginGenerator interface
func (g *AgentGenerator) PrepareGraphDefs() ([]*mkr.GraphDefsParam, error) {
	meta := &pluginMeta{
		Graphs: map[string]customGraphDef{
			"agent.memory": {
				Label: "Agent Memory",
				Unit:  "bytes",
				Metrics: []customGraphMetricDef{
					{Name: "alloc", Label: "Alloc"},
					{Name: "sys", Label: "Sys"},
					{Name: "heapAlloc", Label: "Heap Alloc"},
					{Name: "heapSys", Label: "Heap Sys"},
				},
			},
			"agent.runtime": {
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
