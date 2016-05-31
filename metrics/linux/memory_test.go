// +build linux

package linux

import (
	"os"
	"reflect"
	"testing"

	"github.com/mackerelio/mackerel-agent/metrics"
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

func TestParseMeminfo(t *testing.T) {
	out := []byte(`MemTotal:        1922196 kB
MemFree:          166416 kB
Buffers:          171724 kB
Cached:           647172 kB
SwapCached:        13564 kB
Active:           829688 kB
Inactive:         762348 kB
Active(anon):     338616 kB
Inactive(anon):   434700 kB
Active(file):     491072 kB
Inactive(file):   327648 kB
Unevictable:           0 kB
Mlocked:               0 kB
SwapTotal:       2097148 kB
SwapFree:        2050772 kB
Dirty:               216 kB
Writeback:             8 kB
AnonPages:        760120 kB
Mapped:            17284 kB
Shmem:               176 kB
Slab:             130012 kB
SReclaimable:     107300 kB
SUnreclaim:        22712 kB
KernelStack:        1440 kB
PageTables:         6024 kB
NFS_Unstable:          0 kB
Bounce:                0 kB
WritebackTmp:          0 kB
CommitLimit:     3058244 kB
Committed_AS:    1306640 kB
VmallocTotal:   34359738367 kB
VmallocUsed:       11492 kB
VmallocChunk:   34359722904 kB
HardwareCorrupted:     0 kB
AnonHugePages:    417792 kB
HugePages_Total:       0
HugePages_Free:        0
HugePages_Rsvd:        0
HugePages_Surp:        0
Hugepagesize:       2048 kB
DirectMap4k:        8180 kB
DirectMap2M:     2088960 kB
`)

	result, err := parseMeminfo(out)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := metrics.Values{
		"memory.total":       1968328704,
		"memory.free":        170409984,
		"memory.inactive":    780644352,
		"memory.swap_total":  2147479552,
		"memory.used":        959369216,
		"memory.buffers":     175845376,
		"memory.cached":      662704128,
		"memory.swap_cached": 13889536,
		"memory.active":      849600512,
		"memory.swap_free":   2099990528,
	}
	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %#v", result)
	}
}

func TestParseMeminfoWithMemAvailable(t *testing.T) {
	out := []byte(`
MemTotal: 32767512 kB
MemFree: 263928 kB
MemAvailable: 29702072 kB
Buffers: 342100 kB
Cached: 5376976 kB
SwapCached: 104 kB
Active: 4945908 kB
Inactive: 2857752 kB
Active(anon): 2047984 kB
Inactive(anon): 70596 kB
Active(file): 2897924 kB
Inactive(file): 2787156 kB
Unevictable: 0 kB
Mlocked: 0 kB
SwapTotal: 2097148 kB
SwapFree: 2096476 kB
Dirty: 188 kB
Writeback: 0 kB
AnonPages: 2084328 kB
Mapped: 121500 kB
Shmem: 33996 kB
Slab: 24165452 kB
SReclaimable: 24006480 kB
SUnreclaim: 158972 kB
KernelStack: 11616 kB
PageTables: 258528 kB
NFS_Unstable: 0 kB
Bounce: 0 kB
WritebackTmp: 0 kB
CommitLimit: 18480904 kB
Committed_AS: 7728616 kB
VmallocTotal: 34359738367 kB
VmallocUsed: 328548 kB
VmallocChunk: 34359384060 kB
HardwareCorrupted: 0 kB
AnonHugePages: 22528 kB
HugePages_Total: 0
HugePages_Free: 0
HugePages_Rsvd: 0
HugePages_Surp: 0
Hugepagesize: 2048 kB
DirectMap4k: 181700 kB
DirectMap2M: 11339776 kB
DirectMap1G: 24117248 kB
`)

	result, err := parseMeminfo(out)
	if err != nil {
		t.Errorf("error should be nil but: %s", err)
	}

	expect := metrics.Values{
		"memory.active":      5064609792,
		"memory.swap_total":  2147479552,
		"memory.used":        3139010560,
		"memory.total":       33553932288,
		"memory.buffers":     350310400,
		"memory.cached":      5506023424,
		"memory.swap_cached": 106496,
		"memory.inactive":    2926338048,
		"memory.swap_free":   2146791424,
		"memory.free":        270262272,
		"memory.available":   30414921728,
	}

	if !reflect.DeepEqual(result, expect) {
		t.Errorf("result is not expected one: %#v", result)
	}
}
