// +build linux

package linux

import (
	"testing"
)

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()
	if err != nil {
		t.Error("should not raise error")
	}

	for _, name := range []string{
		"total",
		"free",
		"buffers",
		"cached",
		"active",
		"inactive",
		// "high_total",
		// "high_free",
		// "low_total",
		// "low_free",
		// "dirty",
		// "writeback",
		// "anon_pages",
		// "mapped",
		// "slab",
		// "slab_reclaimable",
		// "slab_unreclaim",
		// "page_tables",
		// "nfs_unstable",
		// "bounce",
		// "commit_limit",
		// "committed_as",
		// "vmalloc_total",
		// "vmalloc_used",
		// "vmalloc_chunk",
		"swap_cached",
		"swap_total",
		"swap_free",
		"used",
	} {
		if _, ok := values["memory."+name]; !ok {
			t.Errorf("memory should has %s", name)
		}
	}
}
