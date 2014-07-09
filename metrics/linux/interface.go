// +build linux

package linux

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
collect network interface I/O

`interface.{interface}.{metric}.delta`: The increased amount of network I/O per minute retrieved from /proc/net/dev

interface = "eth0", "eth1" and so on...

metric = "rxBytes", "rxPackets", "rxErrors", "rxDrops", "rxFifo", "rxFrame", "rxCompressed", "rxMulticast", "txBytes", "txPackets", "txErrors", "txDrops", "txFifo", "txColls", "txCarrier", "txCompressed"

cat /proc/net/dev sample:
	Inter-|   Receive                                                |  Transmit
	 face |bytes    packets errs drop fifo frame compressed multicast|bytes    packets errs drop fifo colls carrier compressed
	eth0: 5461472598 24386569    0    2    0     0          0         0 7215710422 6079810    0    0    0     0       0          0
	lo: 7779878638 1952628    0    0    0     0          0         0 7779878638 1952628    0    0    0     0       0          0
	docker0: 250219988  333736    0    0    0     0          0         0 2024726607 1409929    0    0    0     0       0          0
*/
type InterfaceGenerator struct {
	Interval time.Duration
}

var interfaceMetrics = []string{
	"rxBytes", "rxPackets", "rxErrors", "rxDrops",
	"rxFifo", "rxFrame", "rxCompressed", "rxMulticast",
	"txBytes", "txPackets", "txErrors", "txDrops",
	"txFifo", "txColls", "txCarrier", "txCompressed",
}

// metrics for posting to Mackerel
var postInterfaceMetricsRegexp = regexp.MustCompile(`^interface\..+\.(rxBytes|txBytes)$`)

var interfaceLogger = logging.GetLogger("metrics.interface")

func (g *InterfaceGenerator) Generate() (metrics.Values, error) {
	prevValues, err := g.collectIntarfacesValues()
	if err != nil {
		return nil, err
	}

	interval := g.Interval * time.Second
	time.Sleep(interval)

	currValues, err := g.collectIntarfacesValues()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for name, value := range prevValues {
		if !postInterfaceMetricsRegexp.MatchString(name) {
			continue
		}
		currValue, ok := currValues[name]
		if ok {
			ret[name+".delta"] = (currValue - value) / interval.Seconds()
		}
	}

	return metrics.Values(ret), nil
}

func (g *InterfaceGenerator) collectIntarfacesValues() (metrics.Values, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		interfaceLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}

	lineScanner := bufio.NewScanner(bufio.NewReader(file))
	results := make(map[string]float64)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if matches := regexp.MustCompile(`^\s*([^:]+):\s*(.*)$`).FindStringSubmatch(line); matches != nil {
			name := regexp.MustCompile(`[^A-Za-z0-9_-]`).ReplaceAllString(matches[1], "_")
			if name == "lo" {
				continue
			}

			cols := regexp.MustCompile(`\s+`).Split(matches[2], len(interfaceMetrics))
			if len(cols) < len(interfaceMetrics) {
				continue
			}

			interfaceResult := make(map[string]float64)
			hasNonZeroValue := false
			for i, _ := range interfaceMetrics {
				key := fmt.Sprintf("interface.%s.%s", name, interfaceMetrics[i])
				value, err := strconv.ParseFloat(cols[i], 64)
				if err != nil {
					interfaceLogger.Warningf("Failed to parse host interfaces: %s", err)
					break
				}
				if value != 0 {
					hasNonZeroValue = true
				}
				interfaceResult[key] = value
			}
			if hasNonZeroValue {
				for k, v := range interfaceResult {
					results[k] = v
				}
			}
		}
	}

	// results (eth0) sample
	/** [%!s(*metrics.Value=&{interface.eth0.rxBytes 6.6074069281e+10})
	     %!s(*metrics.Value=&{interface.eth0.rxPackets 1.0483646e+08})
	     %!s(*metrics.Value=&{interface.eth0.rxErrors 0})
	     %!s(*metrics.Value=&{interface.eth0.rxDrops 1})
	     %!s(*metrics.Value=&{interface.eth0.rxFifo 0})
	     %!s(*metrics.Value=&{interface.eth0.rxFrame 0})
	     %!s(*metrics.Value=&{interface.eth0.rxCompressed 0})
	     %!s(*metrics.Value=&{interface.eth0.rxMulticast 0})
	     %!s(*metrics.Value=&{interface.eth0.txBytes 9.180531994e+09})
	     %!s(*metrics.Value=&{interface.eth0.txPackets 5.3107958e+07})
	     %!s(*metrics.Value=&{interface.eth0.txErrors 0})
	     %!s(*metrics.Value=&{interface.eth0.txDrops 0})
	     %!s(*metrics.Value=&{interface.eth0.txFifo 0})
	     %!s(*metrics.Value=&{interface.eth0.txColls 0})
	     %!s(*metrics.Value=&{interface.eth0.txCarrier 0})
	     %!s(*metrics.Value=&{interface.eth0.txCompressed 0})
	    ]
	**/

	return metrics.Values(results), nil
}
