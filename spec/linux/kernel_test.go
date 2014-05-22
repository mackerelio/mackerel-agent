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
		t.Error("should not raise error:", err)
	}

	kernel, typeOk := value.(map[string]string)
	if !typeOk {
		t.Error("value should be map", value)
	}

	if len(kernel["name"]) == 0 {
		t.Error("kernel.name should be filled", kernel["name"])
	}
}
