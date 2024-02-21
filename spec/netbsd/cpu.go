//go:build netbsd
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

// MEMO: sysctl -a machdep.cpu.brand_string

// Generate collects CPU specs.
func (g *CPUGenerator) Generate() (any, error) {
	brandBytes, err := exec.Command("sysctl", "-n", "hw.model").Output()
	if err != nil {
		cpuLogger.Errorf("Failed: %s", err)
		return nil, err
	}

	return mackerel.CPU{
		{"model_name": string(brandBytes)},
	}, nil
}
