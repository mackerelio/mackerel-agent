// +build linux

package linux

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util"
)

/*
collect network interface I/O

`interface.{interface}.{metric}.delta`: The increased amount of network I/O per minute retrieved from /proc/net/dev

interface = "eth0", "eth1" and so on...

see interface_test.go for sample input/output
*/

// InterfaceGenerator XXX
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
var postInterfaceMetricsRegexp = regexp.MustCompile(`^interface\..+\.(?:rxBytes|txBytes)$`)

var interfaceLogger = logging.GetLogger("metrics.interface")

// Generate XXX
func (g *InterfaceGenerator) Generate() (metrics.Values, error) {
	prevValues, err := g.collectInterfacesValues()
	if err != nil {
		return nil, err
	}

	time.Sleep(g.Interval)

	currValues, err := g.collectInterfacesValues()
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
			ret[name+".delta"] = (currValue - value) / g.Interval.Seconds()
		}
	}

	return metrics.Values(ret), nil
}

func (g *InterfaceGenerator) collectInterfacesValues() (metrics.Values, error) {
	out, err := ioutil.ReadFile("/proc/net/dev")
	if err != nil {
		interfaceLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, err
	}
	return parseNetdev(out)
}

func parseNetdev(out []byte) (metrics.Values, error) {
	lineScanner := bufio.NewScanner(bytes.NewReader(out))
	results := make(map[string]float64)
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if kv := strings.SplitN(line, ":", 2); len(kv) == 2 {
			name := util.SanitizeMetricKey(strings.TrimSpace(kv[0]))
			if name == "lo" {
				continue
			}

			cols := strings.Fields(kv[1])
			if len(cols) < len(interfaceMetrics) {
				continue
			}

			interfaceResult := make(map[string]float64)
			hasNonZeroValue := false
			for i, metricName := range interfaceMetrics {
				key := fmt.Sprintf("interface.%s.%s", name, metricName)
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

	return metrics.Values(results), nil
}
