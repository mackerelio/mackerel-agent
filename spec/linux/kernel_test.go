// +build linux

package linux

import (
	"testing"
)

func TestKernelKey(t *testing.T) {
	g := &KernelGenerator{}

	if g.Key() != "kernel" {
		t.Error("key should be kernel")
	}
}

func TestKernelGenerate(t *testing.T) {
	g := &KernelGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}

	kernel, typeOk := value.(map[string]string)
	if !typeOk {
		t.Errorf("value should be map. %+v", value)
	}

	if len(kernel["name"]) == 0 {
		t.Error("kernel.name should be filled")
	}
}
