//go:build darwin
// +build darwin

package darwin

import (
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

// CPUUsageGenerator XXX
type CPUUsageGenerator struct {
}

var cpuUsageLogger = logging.GetLogger("metrics.cpuUsage")

var systemCPUUtilGauge metric.Float64Gauge

var iostatFieldToMetricName = []string{"user", "system", "idle"}

func init() {
	var err error
	systemCPUUtilGauge, err = meter.Float64Gauge(
		semconv.SystemCPUUtilizationName,
		metric.WithUnit(semconv.SystemCPUUtilizationUnit),
		metric.WithDescription(semconv.SystemCPUUtilizationDescription),
	)
	if err != nil {
		panic(err)
	}
}

// Generate returns current CPU usage of the host.
// Keys below are expected:
// - cpu.user.percentage
// - cpu.system.percentage
// - cpu.idle.percentage
func (g *CPUUsageGenerator) Generate() (metrics.Values, error) {

	// $ iostat -n0 -c2
	//         cpu     load average
	//    us sy id   1m   5m   15m
	//    13  7 81  1.93 2.23 2.65
	//    13  7 81  1.93 2.23 2.65
	iostatBytes, err := exec.Command("iostat", "-n0", "-c2").Output()
	if err != nil {
		cpuUsageLogger.Errorf("Failed to invoke iostat: %s", err)
		return nil, err
	}

	return parseIostatOutput(string(iostatBytes))
}

func parseIostatOutput(output string) (metrics.Values, error) {
	ctx := context.TODO()
	lines := strings.Split(output, "\n")

	var fields []string
	for i := len(lines) - 1; i >= 0; i-- {
		if lines[i] == "" {
			continue
		}
		xs := strings.Fields(lines[i])
		if len(xs) >= len(iostatFieldToMetricName) {
			fields = xs
			break
		}
	}

	if len(fields) == 0 {
		return nil, fmt.Errorf("iostat result malformed: [%q]", output)
	}

	cpuUsage := make(map[string]float64, len(iostatFieldToMetricName))

	for i, n := range iostatFieldToMetricName {
		value, err := strconv.ParseFloat(fields[i], 64)
		if err != nil {
			return nil, err
		}

		cpuUsage["cpu."+n+".percentage"] = value

		systemCPUUtilGauge.Record(ctx, value/100, metric.WithAttributes(
			attribute.Key("state").String(n),
		))
	}

	return metrics.Values(cpuUsage), nil
}
