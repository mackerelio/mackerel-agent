// +build windows

package windows

import (
	"errors"
	"fmt"
	"time"

	"github.com/StackExchange/wmi"
	"github.com/mackerelio/golib/logging"
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

type win32PerfFormattedDataPerfDiskPhysicalDisk struct {
	Name             string
	DiskReadsPerSec  uint64
	DiskWritesPerSec uint64
}

// Generate XXX
func (g *DiskGenerator) Generate() (metrics.Values, error) {
	time.Sleep(g.Interval)

	records, err := g.queryWmiWithTimeout()
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
		results[fmt.Sprintf(`disk.%s.reads.delta`, name)] = float64(record.DiskReadsPerSec)
		results[fmt.Sprintf(`disk.%s.writes.delta`, name)] = float64(record.DiskWritesPerSec)
	}
	diskLogger.Debugf("disk %#v", results)
	return results, nil
}

const queryWmiTimeout = 30 * time.Second

func (g *DiskGenerator) queryWmiWithTimeout() ([]win32PerfFormattedDataPerfDiskPhysicalDisk, error) {
	errCh := make(chan error)
	recordsCh := make(chan []win32PerfFormattedDataPerfDiskPhysicalDisk)
	go func() {
		var records []win32PerfFormattedDataPerfDiskPhysicalDisk
		if err := wmi.Query("SELECT * FROM Win32_PerfFormattedData_PerfDisk_LogicalDisk ", &records); err != nil {
			errCh <- err
			return
		}
		recordsCh <- records
	}()
	select {
	case <-time.After(queryWmiTimeout):
		return nil, errors.New("Timeouted while retrieving disk metrics")
	case err := <-errCh:
		return nil, err
	case records := <-recordsCh:
		return records, nil
	}
}
