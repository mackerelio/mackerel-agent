// +build windows

package windows

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestCPUGenerate(t *testing.T) {
	g := &CPUGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	cpu, typeOk := value.(mackerel.CPU)
	if !typeOk {
		t.Errorf("value should be mackerel.CPU. %+v", value)
	}

	if len(cpu) == 0 {
		t.Fatal("should have at least 1 cpu")
	}

	cpu1 := cpu[0]
	if _, ok := cpu1["vendor_id"]; !ok {
		t.Error("cpu should have vendor_id")
	}
	if _, ok := cpu1["family"]; !ok {
		//t.Error("cpu should have family")
	}
	if _, ok := cpu1["model"]; !ok {
		t.Error("cpu should have model")
	}
	if _, ok := cpu1["stepping"]; !ok {
		//t.Error("cpu should have stepping")
	}
	if _, ok := cpu1["physical_id"]; !ok {
		// fails on some environments
		// t.Error("cpu should have physical_id")
	}
	if _, ok := cpu1["core_id"]; !ok {
		// fails on some environments
		// t.Error("cpu should have core_id")
	}
	if _, ok := cpu1["cores"]; !ok {
		// fails on some environments
		// t.Error("cpu should have cores")
	}
	if _, ok := cpu1["model_name"]; !ok {
		t.Error("cpu should have model_name")
	}
	if _, ok := cpu1["mhz"]; !ok {
		t.Error("cpu should have mhz")
	}
	if _, ok := cpu1["cache_size"]; !ok {
		//t.Error("cpu should have cache_size")
	}
	if _, ok := cpu1["flags"]; !ok {
		//t.Error("cpu should have flags")
	}
	if _, ok := cpu1["flags"].([]string); !ok {
		//t.Error("cpu.flags should be slice of string")
	}
}
