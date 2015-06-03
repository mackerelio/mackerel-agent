package metrics

import (
	"testing"
)

func TestAgentGenerate(t *testing.T) {
	g := &AgentGenerator{}
	values, _ := g.Generate()

	agentMetricNames := []string{
		"custom.agent.memory.alloc", "custom.agent.memory.sys",
		"custom.agent.memory.heapAlloc", "custom.agent.memory.heapSys",
	}

	for _, name := range agentMetricNames {
		value, ok := values[name]
		if !ok {
			t.Errorf("AgentGenerator should generate metric value for '%s'", name)
		} else {
			t.Logf("Agent Status '%s' collected: %+v", name, value)
		}
	}
}
