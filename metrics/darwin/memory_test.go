// +build darwin

package darwin

import "testing"

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()

	if err != nil {
		t.Errorf("error should be nil but got: %s", err)
	}

	metricNames := []string{
		"total",
		"free",
		"cached",
		"used",
		"swap_total",
		"swap_free",
	}

	for _, name := range metricNames {
		if _, ok := values["memory."+name]; !ok {
			t.Errorf("memory should have %s", name)
		}
	}

	t.Logf("memory metrics: %+v", values)
}
