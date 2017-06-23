// +build freebsd

package freebsd

import (
	"os/exec"

	"github.com/mackerelio/golib/logging"
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
	brandBytes, err := exec.Command("sysctl", "-n", "hw.model").Output()
	if err != nil {
		cpuLogger.Errorf("Failed: %s", err)
		return nil, err
	}

	return []cpuSpec{
		{"model_name": string(brandBytes)},
	}, nil
}
