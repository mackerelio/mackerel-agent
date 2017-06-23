// +build linux

package linux

import (
	"bufio"
	"os"
	"regexp"

	"github.com/mackerelio/golib/logging"
)

var memItems = map[string]*regexp.Regexp{
	"total":            regexp.MustCompile(`^MemTotal:\s+(\d+) (.+)$`),
	"free":             regexp.MustCompile(`^MemFree:\s+(\d+) (.+)$`),
	"buffers":          regexp.MustCompile(`^Buffers:\s+(\d+) (.+)$`),
	"cached":           regexp.MustCompile(`^Cached:\s+(\d+) (.+)$`),
	"active":           regexp.MustCompile(`^Active:\s+(\d+) (.+)$`),
	"inactive":         regexp.MustCompile(`^Inactive:\s+(\d+) (.+)$`),
	"high_total":       regexp.MustCompile(`^HighTotal:\s+(\d+) (.+)$`),
	"high_free":        regexp.MustCompile(`^HighFree:\s+(\d+) (.+)$`),
	"low_total":        regexp.MustCompile(`^LowTotal:\s+(\d+) (.+)$`),
	"low_free":         regexp.MustCompile(`^LowFree:\s+(\d+) (.+)$`),
	"dirty":            regexp.MustCompile(`^Dirty:\s+(\d+) (.+)$`),
	"writeback":        regexp.MustCompile(`^Writeback:\s+(\d+) (.+)$`),
	"anon_pages":       regexp.MustCompile(`^AnonPages:\s+(\d+) (.+)$`),
	"mapped":           regexp.MustCompile(`^Mapped:\s+(\d+) (.+)$`),
	"slab":             regexp.MustCompile(`^Slab:\s+(\d+) (.+)$`),
	"slab_reclaimable": regexp.MustCompile(`^SReclaimable:\s+(\d+) (.+)$`),
	"slab_unreclaim":   regexp.MustCompile(`^SUnreclaim:\s+(\d+) (.+)$`),
	"page_tables":      regexp.MustCompile(`^PageTables:\s+(\d+) (.+)$`),
	"nfs_unstable":     regexp.MustCompile(`^NFS_Unstable:\s+(\d+) (.+)$`),
	"bounce":           regexp.MustCompile(`^Bounce:\s+(\d+) (.+)$`),
	"commit_limit":     regexp.MustCompile(`^CommitLimit:\s+(\d+) (.+)$`),
	"committed_as":     regexp.MustCompile(`^Committed_AS:\s+(\d+) (.+)$`),
	"vmalloc_total":    regexp.MustCompile(`^VmallocTotal:\s+(\d+) (.+)$`),
	"vmalloc_used":     regexp.MustCompile(`^VmallocUsed:\s+(\d+) (.+)$`),
	"vmalloc_chunk":    regexp.MustCompile(`^VmallocChunk:\s+(\d+) (.+)$`),
	"swap_cached":      regexp.MustCompile(`^SwapCached:\s+(\d+) (.+)$`),
	"swap_total":       regexp.MustCompile(`^SwapTotal:\s+(\d+) (.+)$`),
	"swap_free":        regexp.MustCompile(`^SwapFree:\s+(\d+) (.+)$`),
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

	scanner := bufio.NewScanner(file)

	result := make(map[string]interface{})
	for scanner.Scan() {
		line := scanner.Text()

		for k, v := range memItems {
			if matches := v.FindStringSubmatch(line); matches != nil {
				// ex.) MemTotal:        3916792 kB
				// matches[1] = 3916792, matches[2] = kB
				result[k] = matches[1] + matches[2]
				break
			}
		}
	}
	if err := scanner.Err(); err != nil {
		memoryLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}

	return result, nil
}
