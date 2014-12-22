// +build freebsd

package freebsd

import (
	"regexp"
	"testing"
)

func TestMemoryGenerator_Generate(t *testing.T) {
	g := &MemoryGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() must not fail: %s", err)
	}

	memorySpecs := result.(map[string]string)
	totalMemory, ok := memorySpecs["total"]
	if !ok {
		t.Error("'total' key must exist")
	}

	matched, err := regexp.MatchString(`^\d+kB$`, totalMemory)
	if err != nil {
		t.Fatal(err)
	}
	if !matched {
		t.Errorf("Total must be of form ###kB: %q", totalMemory)
	}
}
