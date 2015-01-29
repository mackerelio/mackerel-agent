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

	counter, err := windows.CreateCounter(g.query, "loadavg5", `\Processor(_Total)\% Processor Time`)
	if err != nil {
		loadavg5Logger.Criticalf(err.Error())
		return nil, err
	}
	g.counters = append(g.counters, counter)
	return g, nil
}

// Generate XXX
func (g *Loadavg5Generator) Generate() (metrics.Values, error) {

	windows.PdhCollectQueryData.Call(uintptr(g.query))

	results := make(map[string]float64)
	for _, v := range g.counters {
		var value windows.PDH_FMT_COUNTERVALUE_ITEM_DOUBLE
		r, _, err := windows.PdhGetFormattedCounterValue.Call(uintptr(v.Counter), windows.PDH_FMT_DOUBLE, uintptr(0), uintptr(unsafe.Pointer(&value)))
		if r != 0 && r != windows.PDH_INVALID_DATA {
			return nil, err
		}
		results[v.PostName] = value.FmtValue.DoubleValue
	}

	loadavg5Logger.Debugf("loadavg5: %q", results)

	return results, nil
}
