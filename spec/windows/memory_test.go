// +build windows

package windows

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	memory, typeOk := value.(mackerel.Memory)
	if !typeOk {
		t.Errorf("value should be mackerel.Memory. %+v", value)
	}

	if _, ok := memory["total"]; !ok {
		t.Error("memory should have total")
	}

	if _, ok := memory["free"]; !ok {
		t.Error("memory should have free")
	}
}
