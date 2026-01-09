//go:build windows

package windows

import (
	"fmt"
	"net"
	"regexp"
	"sort"
	"strings"
	"syscall"
	"time"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
	IgnoreRegexp *regexp.Regexp
	Interval     time.Duration
	query        syscall.Handle
	counters     []*windows.CounterInfo
}

var interfaceLogger = logging.GetLogger("metrics.interface")

func normalizeName(s string) string {
	return strings.Map(func(r rune) rune {
		if ('0' <= r && r <= '9') || ('a' <= r && r <= 'z') || ('A' <= r && r <= 'Z') || r == '-' {
			return r
		}
		return '_'
	}, s)
}

// NewInterfaceGenerator XXX
func NewInterfaceGenerator(ignoreReg *regexp.Regexp, interval time.Duration, useAdapterMetric bool) (*InterfaceGenerator, error) {
	g := &InterfaceGenerator{ignoreReg, interval, 0, nil}

	var err error
	g.query, err = windows.CreateQuery()
	if err != nil {
		interfaceLogger.Criticalf("%s", err.Error())
		return nil, err
	}

	ifs, err := net.Interfaces()
	if err != nil {
		interfaceLogger.Criticalf("%s", err.Error())
		return nil, err
	}

	ai, err := windows.GetAdapterList()
	if err != nil {
		interfaceLogger.Criticalf("%s", err.Error())
		return nil, err
	}

	// make sorted list of names to escape device names.
	names := []string{}
	for _, ad := range ai {
		names = append(names, ad.Name)
	}
	sort.Strings(names)

	// make map for original/escaped. if the names can be duplicated, following
	// name should be renamed with underbar-ed suffixes.
	nameMap := make(map[string]string)
	for _, name := range names {
		escaped := normalizeName(name)
		for {
			if _, ok := nameMap[escaped]; !ok {
				break
			}
			escaped += "_"
		}
		nameMap[name] = escaped
	}

	for _, ifi := range ifs {
		if ifi.Flags&net.FlagLoopback != 0 {
			continue
		}
		for _, ad := range ai {
			if ifi.Index == ad.Index {
				name := ad.Name

				// convert to escaped name
				escaped, ok := nameMap[name]
				if !ok {
					continue
				}
				if g.IgnoreRegexp != nil && g.IgnoreRegexp.MatchString(name) {
					continue
				}
				name = strings.ReplaceAll(name, "(", "[")
				name = strings.ReplaceAll(name, ")", "]")
				name = strings.ReplaceAll(name, "#", "_")
				name = strings.ReplaceAll(name, "/", "_")
				name = strings.ReplaceAll(name, `\`, "_")
				var counter *windows.CounterInfo

				queryType := "Interface"
				if useAdapterMetric {
					queryType = "Adapter"
				}
				counter, err = windows.CreateCounter(
					g.query,
					fmt.Sprintf(`interface.%s.rxBytes.delta`, escaped),
					fmt.Sprintf(`\Network %s(%s)\Bytes Received/sec`, queryType, name))
				if err != nil {
					interfaceLogger.Criticalf(err.Error())
					return nil, err
				}
				g.counters = append(g.counters, counter)
				counter, err = windows.CreateCounter(
					g.query,
					fmt.Sprintf(`interface.%s.txBytes.delta`, escaped),
					fmt.Sprintf(`\Network %s(%s)\Bytes Sent/sec`, queryType, name))
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

	interfaceLogger.Debugf("interface: %#v", results)

	return results, nil
}
