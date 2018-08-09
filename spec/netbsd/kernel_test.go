// +build netbsd

package netbsd

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestKernelGenerator_Generate(t *testing.T) {
	g := &KernelGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() must not fail: %s", err)
	}

	kernel, ok := result.(mackerel.Kernel)
	if !ok {
		t.Fatalf("the result must be of type mackerel.Kernel: %t", result)
	}

	_, osExists := kernel["os"]
	if !osExists {
		t.Errorf("'os' must exit: %v", kernel)
	}

	_, releaseExists := kernel["release"]
	if !releaseExists {
		t.Errorf("'release' must exit: %v", kernel)
	}
}
