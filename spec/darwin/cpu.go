// +build darwin

package darwin

import (
	"bufio"
	"bytes"
	"os/exec"
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
}

func (g *CPUGenerator) parseSysCtlBytes(res []byte) (interface{}, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(res))

	results := map[string]interface{}{}

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

// MEMO: sysctl -a machdep.cpu.brand_string

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
	if err != nil {
		cpuLogger.Errorf("Failed: %s", err)
		return nil, err
	}
	return g.parseSysCtlBytes(cpuInfoBytes)
}
