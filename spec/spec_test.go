package spec

import (
	"fmt"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

type testCPUGenerator struct{}

func (g *testCPUGenerator) Generate() (interface{}, error) {
	return mackerel.CPU{{"cores": "2"}}, nil
}

type testKernelGenerator struct{}

func (g *testKernelGenerator) Generate() (interface{}, error) {
	return mackerel.Kernel{"name": "Linux"}, nil
}

type testErrorGenerator struct{}

func (g *testErrorGenerator) Generate() (interface{}, error) {
	return nil, fmt.Errorf("error")
}

func TestCollect(t *testing.T) {
	generators := []Generator{
		&testCPUGenerator{},
		&testKernelGenerator{},
		&testErrorGenerator{},
	}
	specs := Collect(generators)

	if len(specs.CPU) != 1 {
		t.Error("cpu spec should be collected")
	}

	if len(specs.Kernel) != 1 {
		t.Error("kernel spec should be collected")
	}
}
