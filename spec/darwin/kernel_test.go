// +build darwin

package darwin

import "testing"

func TestKernelGenerator_Generate(t *testing.T) {
	g := &KernelGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() must not fail: %s", err)
	}

	kernel, ok := result.(map[string]string)
	if !ok {
		t.Fatalf("the result must be of type map[string]string: %t", result)
	}

	_, osExists := kernel["os"]
	if !osExists {
		t.Errorf("'os' must exit: %v", kernel)
	}

	_, releaseExists := kernel["release"]
	if !releaseExists {
		t.Errorf("'release' must exit: %v", kernel)
	}

	_, platformNameExists := kernel["platform_name"]
	if !platformNameExists {
		t.Errorf("'platform_name' must exit: %v", kernel)
	}
}
