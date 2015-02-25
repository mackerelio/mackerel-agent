// +build windows

package windows

import (
	"syscall"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// Loadavg5Generator XXX
type Loadavg5Generator struct {
	query    syscall.Handle
	counters []*windows.CounterInfo
}

var loadavg5Logger = logging.GetLogger("metrics.loadavg5")

// NewLoadavg5Generator XXX
func NewLoadavg5Generator() (*Loadavg5Generator, error) {
	g := &Loadavg5Generator{0, nil}

	var err error
	g.query, err = windows.CreateQuery()
	if err != nil {
		loadavg5Logger.Criticalf(err.Error())
		return nil, err
	}

	counter, err := windows.CreateCounter(g.query, "loadavg5", `\System\Processor Queue Length`)
	if err != nil {
		loadavg5Logger.Criticalf(err.Error())
		return nil, err
	}
	g.counters = append(g.counters, counter)
	return g, nil
}

// Generate XXX
func (g *Loadavg5Generator) Generate() (metrics.Values, error) {

	r, _, err := windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			loadavg5Logger.Infof("this metric has not data. ")
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

	loadavg5Logger.Debugf("loadavg5: %q", results)

	return results, nil
}
