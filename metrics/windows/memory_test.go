// +build windows

package windows

import (
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
		"swap_total",
		"swap_free",
		"used",
	} {
		if _, ok := values["memory."+name]; !ok {
			t.Errorf("memory should has %s", name)
		}
	}
}
