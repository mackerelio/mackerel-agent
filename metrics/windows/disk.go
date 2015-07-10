// +build windows

package windows

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// DiskGenerator XXX
type DiskGenerator struct {
	Interval time.Duration
	query    syscall.Handle
	counters []*windows.CounterInfo
}

var diskLogger = logging.GetLogger("metrics.disk")

// NewDiskGenerator XXX
func NewDiskGenerator(interval time.Duration) (*DiskGenerator, error) {
	g := &DiskGenerator{interval, 0, nil}

	var err error
	g.query, err = windows.CreateQuery()
	if err != nil {
		diskLogger.Criticalf(err.Error())
		return nil, err
	}

	drivebuf := make([]byte, 256)
	r, _, err := windows.GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))

	if r == 0 {
		diskLogger.Criticalf(err.Error())
		return nil, err
	}

	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			r, _, err = windows.GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
			if r != windows.DRIVE_FIXED {
				continue
			}
			var counter *windows.CounterInfo

			counter, err = windows.CreateCounter(
				g.query,
				fmt.Sprintf(`disk.%s.reads.delta`, drive),
				fmt.Sprintf(`\PhysicalDisk(0 %s:)\Disk Reads/sec`, drive))
			if err != nil {
				diskLogger.Criticalf(err.Error())
				return nil, err
			}
			g.counters = append(g.counters, counter)

			counter, err = windows.CreateCounter(
				g.query,
				fmt.Sprintf(`disk.%s.writes.delta`, drive),
				fmt.Sprintf(`\PhysicalDisk(0 %s:)\Disk Writes/sec`, drive))
			if err != nil {
				diskLogger.Criticalf(err.Error())
				return nil, err
			}
			g.counters = append(g.counters, counter)
		}
	}

	r, _, err = windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			diskLogger.Infof("this metric has not data. ")
			return nil, err
		}
		diskLogger.Criticalf(err.Error())
		return nil, err
	}

	return g, nil
}

// Generate XXX
func (g *DiskGenerator) Generate() (metrics.Values, error) {
	time.Sleep(g.Interval)

	r, _, err := windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			diskLogger.Infof("this metric has not data. ")
			return nil, err
		}
		return nil, err
	}

	results := make(map[string]float64)
	for _, v := range g.counters {
		var fmtValue windows.PDH_FMT_COUNTERVALUE_DOUBLE
		r, _, err := windows.PdhGetFormattedCounterValue.Call(uintptr(v.Counter), windows.PDH_FMT_DOUBLE, uintptr(0), uintptr(unsafe.Pointer(&fmtValue)))
		if r != 0 && r != windows.PDH_INVALID_DATA {
			return nil, err
		}
		results[v.PostName] = fmtValue.DoubleValue
	}

	diskLogger.Debugf("%q", results)
	return results, nil
}
