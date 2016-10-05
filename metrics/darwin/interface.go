// +build darwin

package darwin

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
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

`interface.{interface}.{metric}.delta`: The increased amount of network I/O per minute retrieved from the result of netstat -bni

interface = "en0", "en1" and so on...

metric = "rxPackets", "rxErrors", "rxBytes", "txPackets", "txErrors", "txBytes",  "colls"

netstat -bni sample:
Name  Mtu   Network       Address            Ipkts Ierrs     Ibytes    Opkts Oerrs     Obytes  Coll
lo0   16384 <Link#1>                       4504403     0 3063793000  4504403     0 3063793000     0
lo0   16384 localhost   ::1                4504403     - 3063793000  4504403     - 3063793000     -
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
			ret[name+".delta"] = (currValue - value) / g.Interval.Seconds()
		}
	}

	return metrics.Values(ret), nil
}

func (g *InterfaceGenerator) collectIntarfacesValues() (metrics.Values, error) {
	out, err := exec.Command("netstat", "-bni").Output()
	if err != nil {
		interfaceLogger.Errorf("Failed to invoke netstat: %s", err)
		return nil, err
	}

	lineScanner := bufio.NewScanner(bytes.NewReader(out))

	results := make(map[string]float64)
	hasNonZeroValue := false
	for lineScanner.Scan() {
		line := lineScanner.Text()
		fields := strings.Fields(line)
		name := util.SanitizeMetricKey(regexp.MustCompile(`\*`).ReplaceAllString(fields[0], ""))
		if match, _ := regexp.MatchString(`^lo\d+$`, name); match {
			continue
		}
		if match, _ := regexp.MatchString(`^<Link#\d+>$`, fields[2]); match {
			rxIndex, txIndex := getFieldIndex(fields)
			rxKey := fmt.Sprintf("interface.%s.rxBytes", name)
			rxValue, rxErr := strconv.ParseFloat(fields[rxIndex], 64)
			if rxErr != nil {
				interfaceLogger.Warningf("Failed to parse host interfaces: %s", err)
				break
			}
			results[rxKey] = rxValue
			if rxValue != 0 {
				hasNonZeroValue = true
			}
			txKey := fmt.Sprintf("interface.%s.txBytes", name)
			txValue, txErr := strconv.ParseFloat(fields[txIndex], 64)
			if txErr != nil {
				interfaceLogger.Warningf("Failed to parse host interfaces: %s", err)
				break
			}
			results[txKey] = txValue
			if txValue != 0 {
				hasNonZeroValue = true
			}
		}
	}
	// results (eth0) sample
	/** [%!s(*metrics.Value=&{interface.eth0.rxBytes 6.6074069281e+10})
	     %!s(*metrics.Value=&{interface.eth0.txBytes 9.180531994e+09})
	    ]
	**/
	if hasNonZeroValue {
		return metrics.Values(results), nil
	}
	return nil, nil
}

func getFieldIndex(fields []string) (rxIndex, txIndex int) {
	if len(fields) == 11 {
		rxIndex = 6
		txIndex = 9
	} else {
		rxIndex = 5
		txIndex = 8
	}
	return
}
