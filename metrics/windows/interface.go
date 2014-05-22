// +build windows

package windows

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	. "github.com/mackerelio/mackerel-agent/util/windows"
)

type InterfaceGenerator struct {
	Interval time.Duration
	query    syscall.Handle
	counters []*CounterInfo
}

var interfaceLogger = logging.GetLogger("metrics.interface")

func NewInterfaceGenerator(interval time.Duration) *InterfaceGenerator {
	g := &InterfaceGenerator{interval, 0, nil}

	var err error
	g.query, err = CreateQuery()
	if err != nil {
		interfaceLogger.Criticalf(err.Error())
		return nil
	}

	ifs, err := net.Interfaces()
	if err != nil {
		interfaceLogger.Criticalf(err.Error())
		return nil
	}

	ai, err := GetAdapterList()
	if err != nil {
		interfaceLogger.Criticalf(err.Error())
		return nil
	}

	for _, ifi := range ifs {
		for ; ai != nil; ai = ai.Next {
			if ifi.Index == int(ai.Index) {
				name := BytePtrToString(&ai.Description[0])
				name = strings.Replace(name, "(", "[", -1)
				name = strings.Replace(name, ")", "]", -1)
				var counter *CounterInfo

				counter, err = CreateCounter(
					g.query,
					fmt.Sprintf(`interface.nic%d.rxBytes.delta`, ifi.Index),
					fmt.Sprintf(`\Network Interface(%s)\Bytes Received/sec`, name))
				if err != nil {
					interfaceLogger.Criticalf(err.Error())
					return nil
				}
				g.counters = append(g.counters, counter)
				counter, err = CreateCounter(
					g.query,
					fmt.Sprintf(`interface.nic%d.txBytes.delta`, ifi.Index),
					fmt.Sprintf(`\Network Interface(%s)\Bytes Sent/sec`, name))
				if err != nil {
					interfaceLogger.Criticalf(err.Error())
					return nil
				}
				g.counters = append(g.counters, counter)
			}
		}
	}

	r, _, err := PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 {
		interfaceLogger.Criticalf(err.Error())
		return nil
	}

	return g
}

func (g *InterfaceGenerator) Generate() (metrics.Values, error) {
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
