// +build windows

package windows

import (
	"syscall"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// ProcessorQueueLengthGenerator is struct of windows api
type ProcessorQueueLengthGenerator struct {
	query    syscall.Handle
	counters []*windows.CounterInfo
}

var processorQueueLengthLogger = logging.GetLogger("metrics.processor_queue_length")

// NewProcessorQueueLengthGenerator is set up windows api
func NewProcessorQueueLengthGenerator() (*ProcessorQueueLengthGenerator, error) {
	g := &ProcessorQueueLengthGenerator{0, nil}

	var err error
	g.query, err = windows.CreateQuery()
	if err != nil {
		processorQueueLengthLogger.Criticalf(err.Error())
		return nil, err
	}

	counter, err := windows.CreateCounter(g.query, "processor_queue_length", `\System\Processor Queue Length`)
	if err != nil {
		processorQueueLengthLogger.Criticalf(err.Error())
		return nil, err
	}
	g.counters = append(g.counters, counter)
	return g, nil
}

// Generate XXX
func (g *ProcessorQueueLengthGenerator) Generate() (metrics.Values, error) {

	r, _, err := windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			processorQueueLengthLogger.Infof("this metric has not data. ")
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

	processorQueueLengthLogger.Debugf("processor_queue_length: %#v", results)

	return results, nil
}
