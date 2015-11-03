// +build netbsd

package netbsd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

// CPUUsageGenerator XXX
type CPUUsageGenerator struct {
}

var cpuUsageLogger = logging.GetLogger("metrics.cpuUsage")

var iostatFieldToMetricName = []string{"user", "nice", "system", "interrupt", "idle"}

// Generate returns current CPU usage of the host.
// Keys below are expected:
// - cpu.user.percentage
// - cpu.system.percentage
// - cpu.idle.percentage
func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {

	// % iostat -c2 -C
	//             CPU
	// us ni sy in id
	//  0  0  0  0 100
	//  0  0  0  0 100

	iostatBytes, err := exec.Command("iostat", "-c2", "-d", "-C").Output()
	if err != nil {
		cpuUsageLogger.Errorf("Failed to invoke iostat: %s", err)
		return nil, err
	}

	iostat := string(iostatBytes)
	lines := strings.Split(iostat, "\n")
	if len(lines) != 5 {
		return nil, fmt.Errorf("iostat result malformed: [%q]", iostat)
	}

	fields := strings.Fields(lines[3])
	if len(fields) < len(iostatFieldToMetricName) {
		return nil, fmt.Errorf("iostat result malformed: [%q]", iostat)
	}

	cpuUsage := make(map[string]float64, len(iostatFieldToMetricName))

	for i, n := range iostatFieldToMetricName {
		if i == 3 {
			continue
		}
		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, err
		}

		cpuUsage["cpu."+n+".percentage"] = value
	}

	return metrics.Values(cpuUsage), nil
}
