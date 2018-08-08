// +build windows

package windows

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestKernelGenerate(t *testing.T) {
	g := &KernelGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %s", err)
	}

	kernel, typeOk := value.(mackerel.Kernel)
	if !typeOk {
		t.Errorf("value should be mackerel.Kernel. %+v", value)
	}

	if len(kernel["name"]) == 0 {
		t.Error("kernel.name should be filled")
	}
}
