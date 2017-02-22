// +build darwin

package darwin

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

// CPUGenerator Collects CPU specs
type CPUGenerator struct {
}

// Key XXX
func (g *CPUGenerator) Key() string {
	return "cpu"
}

var cpuLogger = logging.GetLogger("spec.cpu")

type cpuSpec map[string]interface{}

var sysCtlKeyMap = map[string]string{
	"core_count":   "cores",
	"brand_string": "model_name",
	"model":        "model",
	"vendor":       "vendor_id",
	"family":       "family",
	"stepping":     "stepping",
}

func (g *CPUGenerator) parseSysCtlBytes(res []byte) (cpuSpec, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(res))

	results := cpuSpec{}

	for scanner.Scan() {
		line := scanner.Text()
		kv := strings.SplitN(line, ":", 2)
		if len(kv) < 2 {
			continue
		}
		key := strings.TrimPrefix(strings.TrimSpace(kv[0]), "machdep.cpu.")
		val := strings.TrimSpace(kv[1])
		if label, ok := sysCtlKeyMap[key]; ok {
			results[label] = val
		}
	}

	if err := scanner.Err(); err != nil {
		cpuLogger.Errorf("Failed (skip this spec): %s", err)
		return nil, err
	}
	return results, nil
}

func (g *CPUGenerator) getCPUCount() (*int, error) {
	countBytes, err := exec.Command("sysctl", "-n", "hw.logicalcpu").Output()
	if err != nil {
		cpuLogger.Errorf("Failed: %s", err)
		return nil, err
	}
	count, err := strconv.Atoi(strings.TrimSpace(string(countBytes)))
	if err != nil {
		cpuLogger.Errorf("Failed to parse: %s", err)
		return nil, err
	}
	return &count, nil
}

// MEMO: sysctl -a machdep.cpu

// Generate collects CPU specs.
// Returns an array of cpuSpec.
// Each spec is expected to have keys below:
// - model_name (used in Web)
// - vendor_id
// - family
// - model
// - stepping
// - physical_id
// - core_id
// - cores
// - mhz
// - cache_size
// - flags
func (g *CPUGenerator) Generate() (interface{}, error) {
	cpuInfoBytes, err := exec.Command("sysctl", "-a", "machdep.cpu").Output()
	cpuInfo, err := g.parseSysCtlBytes(cpuInfoBytes)
	if err != nil {
		cpuLogger.Errorf("Failed: %s", err)
		return nil, err
	}
	cpuCount, err := g.getCPUCount()
	if err != nil {
		cpuLogger.Errorf("Failed: %s", err)
		return nil, err
	}
	results := make([]cpuSpec, *cpuCount)
	for i := 0; i < *cpuCount; i++ {
		results[i] = cpuInfo
	}
	return results, nil
}
