// +build darwin

package darwin

import "testing"

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	metricNames := []string{
		"total",
		"free",
		"cached",
		"active",
		"inactive",
		"used",
		"swap_total",
		"swap_free",
	}

	for _, name := range metricNames {
		if v, ok := values["memory."+name]; !ok {
			t.Errorf("memory should has %s", name)
		} else {
			t.Logf("memory '%s' collected: %+v", name, v)
		}
	}
}
