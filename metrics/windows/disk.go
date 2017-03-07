// +build windows

package windows

import (
	"fmt"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// DiskGenerator XXX
type DiskGenerator struct {
	Interval time.Duration
}

var diskLogger = logging.GetLogger("metrics.disk")

// NewDiskGenerator XXX
func NewDiskGenerator(interval time.Duration) (*DiskGenerator, error) {
	return &DiskGenerator{interval}, nil
}

type win32PerfFormattedData struct {
	Name                 string
	DiskReadBytesPersec  uint64
	DiskWriteBytesPersec uint64
}

// Generate XXX
func (g *DiskGenerator) Generate() (metrics.Values, error) {
	time.Sleep(g.Interval)

	var records []win32PerfFormattedData
	err := wmi.Query("SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk ", &records)
	if err != nil {
		return nil, err
	}

	results := make(map[string]float64)
	for _, record := range records {
		name := record.Name
		// Collect metrics for only drives
		if len(name) != 2 || name[1] != ':' {
			continue
		}
		name = name[:1]
		results[fmt.Sprintf(`disk.%s.reads.delta`, name)] = float64(record.DiskReadBytesPersec)
		results[fmt.Sprintf(`disk.%s.writes.delta`, name)] = float64(record.DiskWriteBytesPersec)
	}
	diskLogger.Debugf("%q", results)
	return results, nil
}
