// +build windows

package windows

import (
	"testing"
)

func TestMemoryKey(t *testing.T) {
	g := &MemoryGenerator{}

	if g.Key() != "memory" {
		t.Error("key should be memory")
	}
}

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	memory, typeOk := value.(map[string]interface{})
	if !typeOk {
		t.Errorf("value should be map. %+v", value)
	}

	if _, ok := memory["total"]; !ok {
		t.Error("memory should have total")
	}

	if _, ok := memory["free"]; !ok {
		t.Error("memory should have free")
	}
}
