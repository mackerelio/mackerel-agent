// +build linux

package linux

import (
	"os"
	"testing"
)

func TestMemoryGenerator(t *testing.T) {
	g := &MemoryGenerator{}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
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
		if v, ok := values["memory."+name]; !ok {
			if name == "swap_cached" && os.Getenv("TRAVIS") != "" {
				t.Logf("memory '%s' is not collected in Travis", name)
			} else {
				t.Errorf("memory should has %s", name)
			}
		} else {
			t.Logf("memory '%s' collected: %+v", name, v)
		}
	}
}
