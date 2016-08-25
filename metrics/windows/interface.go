// +build windows

package windows

import (
	"fmt"
	"net"
	"strings"
	"syscall"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
	Interval time.Duration
	query    syscall.Handle
	counters []*windows.CounterInfo
}

var interfaceLogger = logging.GetLogger("metrics.interface")

// NewInterfaceGenerator XXX
func NewInterfaceGenerator(interval time.Duration) (*InterfaceGenerator, error) {
	g := &InterfaceGenerator{interval, 0, nil}

	var err error
	g.query, err = windows.CreateQuery()
	if err != nil {
		interfaceLogger.Criticalf(err.Error())
		return nil, err
	}

	ifs, err := net.Interfaces()
	if err != nil {
		interfaceLogger.Criticalf(err.Error())
		return nil, err
	}

	ai, err := windows.GetAdapterList()
	if err != nil {
		interfaceLogger.Criticalf(err.Error())
		return nil, err
	}

	first := ai

	for _, ifi := range ifs {
		for ai = first; ai != nil; ai = ai.Next {
			if ifi.Index == int(ai.Index) {
				name, err := windows.AnsiBytePtrToString(&ai.Description[0])
				if err != nil {
					name = windows.BytePtrToString(&ai.Description[0])
				}
				name = strings.Replace(name, "(", "[", -1)
				name = strings.Replace(name, ")", "]", -1)
				name = strings.Replace(name, "#", "_", -1)
				var counter *windows.CounterInfo

				counter, err = windows.CreateCounter(
					g.query,
					fmt.Sprintf(`interface.nic%d.rxBytes.delta`, ifi.Index),
					fmt.Sprintf(`\Network Interface(%s)\Bytes Received/sec`, name))
				if err != nil {
					interfaceLogger.Criticalf(err.Error())
					return nil, err
				}
				g.counters = append(g.counters, counter)
				counter, err = windows.CreateCounter(
					g.query,
					fmt.Sprintf(`interface.nic%d.txBytes.delta`, ifi.Index),
					fmt.Sprintf(`\Network Interface(%s)\Bytes Sent/sec`, name))
				if err != nil {
					interfaceLogger.Criticalf(err.Error())
					return nil, err
				}
				g.counters = append(g.counters, counter)
			}
		}
	}

	r, _, err := windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			interfaceLogger.Infof("this metric has not data. ")
			return nil, err
		}
		interfaceLogger.Criticalf(err.Error())
		return nil, err
	}

	return g, nil
}

// Generate XXX
func (g *InterfaceGenerator) Generate() (metrics.Values, error) {

	time.Sleep(g.Interval)

	r, _, err := windows.PdhCollectQueryData.Call(uintptr(g.query))
	if r != 0 && err != nil {
		if r == windows.PDH_NO_DATA {
			interfaceLogger.Infof("this metric has not data. ")
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

	interfaceLogger.Debugf("%q", results)

	return results, nil
}
