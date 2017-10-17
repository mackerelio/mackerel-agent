// +build linux

package linux

import (
	"bufio"
	"io"
	"os"
	"strings"

	"github.com/mackerelio/golib/logging"
)

var memItems = map[string]string{
	"MemTotal:":     "total",
	"MemFree:":      "free",
	"Buffers:":      "buffers",
	"Cached:":       "cached",
	"Active:":       "active",
	"Inactive:":     "inactive",
	"HighTotal:":    "high_total",
	"HighFree:":     "high_free",
	"LowTotal:":     "low_total",
	"LowFree:":      "low_free",
	"Dirty:":        "dirty",
	"Writeback:":    "writeback",
	"AnonPages:":    "anon_pages",
	"Mapped:":       "mapped",
	"Slab:":         "slab",
	"SReclaimable:": "slab_reclaimable",
	"SUnreclaim:":   "slab_unreclaim",
	"PageTables:":   "page_tables",
	"NFS_Unstable:": "nfs_unstable",
	"Bounce:":       "bounce",
	"CommitLimit:":  "commit_limit",
	"Committed_AS:": "committed_as",
	"VmallocTotal:": "vmalloc_total",
	"VmallocUsed:":  "vmalloc_used",
	"VmallocChunk:": "vmalloc_chunk",
	"SwapCached:":   "swap_cached",
	"SwapTotal:":    "swap_total",
	"SwapFree:":     "swap_free",
}

// MemoryGenerator collects the host's memory specs.
type MemoryGenerator struct {
}

// Key XXX
func (g *MemoryGenerator) Key() string {
	return "memory"
}

var memoryLogger = logging.GetLogger("spec.memory")

// Generate XXX
func (g *MemoryGenerator) Generate() (interface{}, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		memoryLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}
	defer file.Close()
	return generateMemorySpec(file)
}

func generateMemorySpec(out io.Reader) (map[string]string, error) {
	scanner := bufio.NewScanner(out)

	result := make(map[string]string)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) < 2 {
			continue
		}
		if k, ok := memItems[fields[0]]; ok {
			// ex) MemTotal:  3916792 kB
			//   -> "total": "3916782kB"
			result[k] = strings.Join(fields[1:], "")
		}
	}
	if err := scanner.Err(); err != nil {
		memoryLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}

	return result, nil
}
