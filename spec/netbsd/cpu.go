// +build netbsd

package netbsd

import (
	"os/exec"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
)

// CPUGenerator Collects CPU specs
type CPUGenerator struct {
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

	return mackerel.CPU{
		{"model_name": string(brandBytes)},
	}, nil
}
