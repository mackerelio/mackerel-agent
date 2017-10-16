// +build linux

package linux

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

	memItemKeys := []string{
		"total",
		"free",
		"buffers",
		"cached",
		"active",
		"inactive",
		"dirty",
		"writeback",
		"anon_pages",
		"mapped",
		"slab",
		"slab_reclaimable",
		"slab_unreclaim",
		"page_tables",
		"nfs_unstable",
		"bounce",
		"commit_limit",
		"committed_as",
		"vmalloc_total",
		"vmalloc_used",
		"vmalloc_chunk",
		"swap_cached",
		"swap_total",
		"swap_free",
	}

	for _, key := range memItemKeys {
		if _, ok := memory[key]; !ok {
			t.Errorf("memory spec should have %s", key)
		}
	}
}
