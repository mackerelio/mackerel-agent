// +build windows

package windows

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	. "github.com/mackerelio/mackerel-agent/util/windows"
)

type DiskGenerator struct {
	Interval time.Duration
	query    syscall.Handle
	counters []*CounterInfo
}

var diskLogger = logging.GetLogger("metrics.disk")

func NewDiskGenerator(interval time.Duration) (*DiskGenerator, error) {
	g := &DiskGenerator{interval, 0, nil}

	var err error
	g.query, err = CreateQuery()
	if err != nil {
		diskLogger.Criticalf(err.Error())
		return nil, err
	}

	drivebuf := make([]byte, 256)
	_, r, err := GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))
	if r != 0 {
		diskLogger.Criticalf(err.Error())
		return nil, err
	}

	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			r, _, err = GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
			if r != DRIVE_FIXED {
				continue
			}
			var counter *CounterInfo

			counter, err = CreateCounter(
				g.query,
				fmt.Sprintf(`disk.%s.reads.delta`, drive),
				fmt.Sprintf(`\PhysicalDisk(0 %s:)\Disk Reads/sec`, drive))
			if err != nil {
				diskLogger.Criticalf(err.Error())
				return nil, err
			}
			g.counters = append(g.counters, counter)

			counter, err = CreateCounter(
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

	r, _, err = PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		diskLogger.Criticalf(err.Error())
		return nil, err
	}
	return g, nil
}

func (g *DiskGenerator) Generate() (metrics.Values, error) {
	interval := g.Interval * time.Second
	time.Sleep(interval)

	r, _, err := PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 {
		return nil, err
	}

	results := make(map[string]float64)
	for _, v := range g.counters {
		var value PDH_FMT_COUNTERVALUE_ITEM_DOUBLE
		r, _, err = PdhGetFormattedCounterValue.Call(uintptr(v.Counter), PDH_FMT_DOUBLE, uintptr(0), uintptr(unsafe.Pointer(&value)))
		if r != 0 && r != PDH_INVALID_DATA {
			return nil, err
		}
		results[v.PostName] = value.FmtValue.DoubleValue
	}
	return results, nil
}
