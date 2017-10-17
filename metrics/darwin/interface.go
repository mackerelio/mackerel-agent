// +build darwin

package darwin

import (
	"time"

	"github.com/mackerelio/go-osstat/network"
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
collect network interface I/O

`interface.{interface}.{metric}.delta`: The increased amount of network I/O per minute retrieved from the result of netstat -bni

interface = "en0", "en1" and so on...
*/

// InterfaceGenerator XXX
type InterfaceGenerator struct {
	Interval time.Duration
}

// metrics for posting to Mackerel

var interfaceLogger = logging.GetLogger("metrics.interface")

// Generate XXX
func (g *InterfaceGenerator) Generate() (metrics.Values, error) {
	prevValues, err := g.collectIntarfacesValues()
	if err != nil {
		return nil, err
	}

	time.Sleep(g.Interval)

	currValues, err := g.collectIntarfacesValues()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for name, value := range prevValues {
		currValue, ok := currValues[name]
		if ok {
			ret[name+".delta"] = float64(currValue-value) / g.Interval.Seconds()
		}
	}

	return metrics.Values(ret), nil
}

func (g *InterfaceGenerator) collectIntarfacesValues() (map[string]uint64, error) {
	networks, err := network.Get()
	if err != nil {
		interfaceLogger.Errorf("failed to get network statistics: %s", err)
		return nil, err
	}
	if len(networks) == 0 {
		return nil, nil
	}
	results := make(map[string]uint64, len(networks)*2)
	for _, network := range networks {
		results["interface."+network.Name+".rxBytes"] = network.RxBytes
		results["interface."+network.Name+".txBytes"] = network.TxBytes
	}
	return results, nil
}
