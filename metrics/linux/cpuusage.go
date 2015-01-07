// +build linux

package linux

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
)

/*
collect CPU usage

`cpu.{metric}.percentage`: The increased amount of CPU time per minute as percentage of total CPU cores x 100

metric = "user", "nice", "system", "idle", "iowait", "irq", "softirq", "steal", "guest"

graph: stacks `cpu.{metric}.percentage`

cat /proc/stat sample: {{{
	cpu  7792253 5479 4851396 18056319678 127239 0 146818 2383839
	cpu0 5385397 1412 1970781 4509432750 103260 0 136689 876389
	cpu1 641247 1361 782257 4516019361 7247 0 2403 452803
	cpu2 652342 1366 617100 4516172153 7762 0 2447 453509
	cpu3 1113265 1339 1481257 4514695413 8968 0 5278 601135
	intr 6664031039 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 3682869251 40969382 60 304 40427429 141 567698585 39988217 145 500771676 67725387 95 1170166889 187 33636967 83463 519692861 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0
	ctxt 14007527061
	btime 1349954031
	processes 60807520
	procs_running 1
	procs_blocked 0
}}}
*/
type CpuusageGenerator struct {
	Interval time.Duration
}

// In additions these metrics, collect *.percentage metrics
var cpuusageMetricNames = []string{
	"cpu.user", "cpu.nice", "cpu.system", "cpu.idle", "cpu.iowait",
	"cpu.irq", "cpu.softirq", "cpu.steal", "cpu.guest",
}

var cpuNumberPattern = regexp.MustCompile(`^cpu\d+\s`)

var cpuusageLogger = logging.GetLogger("metrics.cpuusage")

// Generate XXX
func (g *CpuusageGenerator) Generate() (metrics.Values, error) {
	prevValues, prevTotal, _, err := g.collectProcStatValues()
	if err != nil {
		return nil, err
	}

	time.Sleep(g.Interval * time.Second)

	currValues, currTotal, cpuCount, err := g.collectProcStatValues()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for i, name := range cpuusageMetricNames {
		// Values in /proc/stat differ in Linux kernel versions.
		// Not all metrics in cpuusageMetricNames can be retrieved.
		// ref: `man 5 proc`
		if i >= len(currValues) || i >= len(prevValues) {
			break
		}

		// percentage of increased amount of CPU time
		ret[name+".percentage"] = (currValues[i] - prevValues[i]) * 100.0 * float64(cpuCount) / (currTotal - prevTotal)
	}

	return metrics.Values(ret), nil
}

// returns values corresponding to cpuusageMetricNames, those total and the number of CPUs
func (g *CpuusageGenerator) collectProcStatValues() ([]float64, float64, uint, error) {
	file, err := os.Open("/proc/stat")
	if err != nil {
		cpuusageLogger.Errorf("Failed (skip these metrics): %s", err)
		return nil, 0, 0, err
	}

	lineScanner := bufio.NewScanner(bufio.NewReader(file))

	var cols []string
	var cpuCount uint
	var firstLine bool = true

	for lineScanner.Scan() {
		line := lineScanner.Text()

		if firstLine {
			// first line contains total values of all CPUs
			cols = strings.Fields(lineScanner.Text())[1:]
			firstLine = false
		} else if cpuNumberPattern.MatchString(line) {
			// number of cores
			cpuCount += 1
		} else {
			break
		}
	}

	values := make([]float64, len(cols))

	var totalValues float64
	for i, strValue := range cols {
		values[i], err = strconv.ParseFloat(strValue, 64)
		if err != nil {
			cpuusageLogger.Errorf("Failed to parse cpuusage metrics (skip these metrics): %s", err)
			return nil, 0, 0, err
		}
		totalValues += values[i]
	}

	return values, totalValues, cpuCount, nil
}
