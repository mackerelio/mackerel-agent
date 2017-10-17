// +build linux

package linux

import (
	"reflect"
	"strings"
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

	memory, typeOk := value.(map[string]string)
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

func TestGenerateMemorySpec(t *testing.T) {
	got, err := generateMemorySpec(strings.NewReader(
		`MemTotal:       15434208 kB
MemFree:         3009856 kB
MemAvailable:   10008916 kB
Buffers:          443104 kB
Cached:          6305168 kB
SwapCached:            0 kB
Active:          7662608 kB
Inactive:        3375044 kB
Active(anon):    4474764 kB
Inactive(anon):   276180 kB
Active(file):    3187844 kB
Inactive(file):  3098864 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:             0 kB
SwapFree:              0 kB
Dirty:               516 kB
Writeback:             0 kB
AnonPages:       4289480 kB
Mapped:           258168 kB
Shmem:            461560 kB
Slab:            1169796 kB
SReclaimable:     965756 kB
SUnreclaim:       204040 kB
KernelStack:       42688 kB
PageTables:        61024 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     7717104 kB
Committed_AS:   23796544 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       42608 kB
VmallocChunk:   34359578580 kB
HardwareCorrupted:     0 kB
AnonHugePages:         0 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:       65536 kB
DirectMap2M:     3211264 kB
DirectMap1G:    12582912 kB
`))
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	expected := map[string]string{
		"total":            "15434208kB",
		"free":             "3009856kB",
		"buffers":          "443104kB",
		"cached":           "6305168kB",
		"active":           "7662608kB",
		"inactive":         "3375044kB",
		"dirty":            "516kB",
		"writeback":        "0kB",
		"anon_pages":       "4289480kB",
		"mapped":           "258168kB",
		"slab":             "1169796kB",
		"slab_reclaimable": "965756kB",
		"slab_unreclaim":   "204040kB",
		"page_tables":      "61024kB",
		"nfs_unstable":     "0kB",
		"bounce":           "0kB",
		"commit_limit":     "7717104kB",
		"committed_as":     "23796544kB",
		"vmalloc_total":    "34359738367kB",
		"vmalloc_used":     "42608kB",
		"vmalloc_chunk":    "34359578580kB",
		"swap_cached":      "0kB",
		"swap_total":       "0kB",
		"swap_free":        "0kB",
	}
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("invalid memory stat: %+v (expected: %+v)", got, expected)
	}
}
