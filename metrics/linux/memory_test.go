// +build linux

package linux

import (
	"testing"
)

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should be nil but got: %s", err)
	}

	metricNames := []string{
		"total",
		"free",
		"buffers",
		"cached",
		"active",
		"inactive",
		"swap_cached",
		"swap_total",
		"swap_free",
		"used",
	}

	for _, name := range metricNames {
		if _, ok := values["memory."+name]; !ok {
			t.Errorf("memory should have %s", name)
		}
	}

	t.Logf("memory metrics: %+v", values)
}
