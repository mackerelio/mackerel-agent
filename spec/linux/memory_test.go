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

	if _, ok := memory["total"]; !ok {
		t.Error("memory should has total")
	}

	if _, ok := memory["free"]; !ok {
		t.Error("memory should has free")
	}

	if _, ok := memory["buffers"]; !ok {
		t.Error("memory should has buffers")
	}

	if _, ok := memory["cached"]; !ok {
		t.Error("memory should has cached")
	}

	if _, ok := memory["active"]; !ok {
		t.Error("memory should has active")
	}

	if _, ok := memory["inactive"]; !ok {
		t.Error("memory should has inactive")
	}

	if _, ok := memory["high_total"]; !ok {
		t.Log("Skip: memory should has high_total")
	}

	if _, ok := memory["high_free"]; !ok {
		t.Log("Skip: memory should has high_free")
	}

	if _, ok := memory["low_total"]; !ok {
		t.Log("Skip: memory should has low_tatal")
	}

	if _, ok := memory["low_free"]; !ok {
		t.Log("Skip: memory should has low_free")
	}

	if _, ok := memory["dirty"]; !ok {
		t.Error("memory should has dirty")
	}

	if _, ok := memory["writeback"]; !ok {
		t.Error("memory should has writeback")
	}

	if _, ok := memory["anon_pages"]; !ok {
		t.Error("memory should has anon_pages")
	}

	if _, ok := memory["mapped"]; !ok {
		t.Error("memory should has mapped")
	}

	if _, ok := memory["slab"]; !ok {
		t.Error("memory should has slab")
	}

	if _, ok := memory["slab_reclaimable"]; !ok {
		t.Error("memory should has slab_reclaimable")
	}

	if _, ok := memory["slab_unreclaim"]; !ok {
		t.Error("memory should has slab_unreclaim")
	}

	if _, ok := memory["page_tables"]; !ok {
		t.Error("memory should has page_tables")
	}

	if _, ok := memory["nfs_unstable"]; !ok {
		t.Error("memory should has nfs_unstable")
	}

	if _, ok := memory["bounce"]; !ok {
		t.Error("memory should has bounce")
	}

	if _, ok := memory["commit_limit"]; !ok {
		t.Error("memory should has commit_limmit")
	}

	if _, ok := memory["committed_as"]; !ok {
		t.Error("memory should has committed_as")
	}

	if _, ok := memory["vmalloc_total"]; !ok {
		t.Error("memory should has vmalloc_total")
	}

	if _, ok := memory["vmalloc_used"]; !ok {
		t.Error("memory should has vmalloc_used")
	}

	if _, ok := memory["vmalloc_chunk"]; !ok {
		t.Error("memory should has vmalloc_chunk")
	}

	if _, ok := memory["swap_cached"]; !ok {
		t.Error("memory should has swap_cached")
	}

	if _, ok := memory["swap_total"]; !ok {
		t.Error("memory should has swap_total")
	}

	if _, ok := memory["swap_free"]; !ok {
		t.Error("memory should has swap_free")
	}
}
