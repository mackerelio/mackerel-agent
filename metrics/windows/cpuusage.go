// +build windows

package windows

import (
	"syscall"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// CPUUsageGenerator is struct of windows api
type CPUUsageGenerator struct {
	query    syscall.Handle
	counters []*windows.CounterInfo
}

var cpuUsageLogger = logging.GetLogger("cpu.user.percentage")

// NewCPUUsageGenerator is set up windows api
func NewCPUUsageGenerator() (*CPUUsageGenerator, error) {
	g := &CPUUsageGenerator{0, nil}

	var err error
	g.query, err = windows.CreateQuery()
	if err != nil {
		cpuUsageLogger.Criticalf("%s", err.Error())
		return nil, err
	}
	var counter *windows.CounterInfo

	counter, err = windows.CreateCounter(g.query, "cpu.user.percentage", `\Processor(_Total)\% User Time`)
	if err != nil {
		cpuUsageLogger.Criticalf("%s", err.Error())
		return nil, err
	}
	g.counters = append(g.counters, counter)

	counter, err = windows.CreateCounter(g.query, "cpu.system.percentage", `\Processor(_Total)\% Privileged Time`)
	if err != nil {
		cpuUsageLogger.Criticalf("%s", err.Error())
		return nil, err
	}
	g.counters = append(g.counters, counter)

	counter, err = windows.CreateCounter(g.query, "cpu.idle.percentage", `\Processor(_Total)\% Idle Time`)
	if err != nil {
		cpuUsageLogger.Criticalf("%s", err.Error())
		return nil, err
	}
	g.counters = append(g.counters, counter)
	return g, nil
}

// Generate XXX
func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {

	r, _, err := windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			cpuUsageLogger.Infof("this metric has not data. ")
			return nil, err
		}
		return nil, err
	}

	results := make(map[string]float64)
	for _, v := range g.counters {
		results[v.PostName], err = windows.GetCounterValue(v.Counter)
		if err != nil {
			return nil, err
		}
	}

	cpuUsageLogger.Debugf("cpuusage: %#v", results)

	return results, nil
}
