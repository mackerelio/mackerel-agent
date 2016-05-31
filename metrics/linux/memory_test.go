// +build linux

package linux

import (
	"os"
	"testing"
)

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	for _, name := range []string{
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
	} {
		if v, ok := values["memory."+name]; !ok {
			if name == "swap_cached" && os.Getenv("TRAVIS") != "" {
				t.Logf("memory '%s' is not collected in Travis", name)
			} else {
				t.Errorf("memory should has %s", name)
			}
		} else {
			t.Logf("memory '%s' collected: %+v", name, v)
		}
	}
}
