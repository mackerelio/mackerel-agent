// +build freebsd

package freebsd

import "testing"

func TestCPUGenerator_Generate(t *testing.T) {
	g := &CPUGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() must not fail: %s", err)
	}

	cpus, ok := result.([]cpuSpec)
	if !ok {
		t.Fatalf("the result must be of type []cpuSpec: %T", result)
	}

	if len(cpus) == 0 {
		t.Fatalf("cpu specs must not be empty: %v", cpus)
	}

	_, modelNameExists := cpus[0]["model_name"]
	if !modelNameExists {
		t.Errorf("'model_name' must exit: %v", cpus[0])
	}
}
