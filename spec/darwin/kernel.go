// +build darwin

package darwin

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

// KernelGenerator Generates specs about the kernel.
type KernelGenerator struct {
}

// Key XXX
func (g *KernelGenerator) Key() string {
	return "kernel"
}

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate collects specs from `uname` command and `sw_vers` command
func (g *KernelGenerator) Generate() (interface{}, error) {
	// foundamental information from `uname` command
	unameArgs := map[string][]string{
		"release": {"-r"},
		"version": {"-v"},
		"machine": {"-m"},
		"os":      {"-s"},
	}

	unames := make(map[string]string, len(unameArgs))

	for field, args := range unameArgs {
		out, err := exec.Command("/usr/bin/uname", args...).Output()
		if err != nil {
			kernelLogger.Errorf("Failed to run uname %s (skip this field): %s", args, err)
			continue
		}

		unames[field] = strings.TrimSpace(string(out))
	}

	// platform information from `sw_vers` command
	swVerArgs := map[string][]string{
		"productName":    {"-productName"},
		"productVersion": {"-productVersion"},
	}

	swVers := make(map[string]string, len(swVerArgs))

	for field, args := range swVerArgs {
		out, err := exec.Command("/usr/bin/sw_vers", args...).Output()
		if err != nil {
			kernelLogger.Errorf("Failed to run sw_vers %s (skip this field): %s", args, err)
			continue
		}

		swVers[field] = strings.TrimSpace(string(out))
	}

	results := map[string]string{
		"release":          unames["release"],
		"version":          unames["version"],
		"machine":          unames["machine"],
		"os":               unames["os"],
		"name":             unames["os"], // same as name
		"platform_name":    swVers["productName"],
		"platform_version": swVers["productVersion"],
	}

	return results, nil
}
