// +build darwin

package darwin

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

// KernelGenerator Generates specs about the kernel.
// Keys below are expected.
// - name:    the operating system name ("Linux")
// - release: the operating system release ("2.6.32-5-686")
// - version: the operating system version ("#1 SMP Sun Sep 23 09:49:36 UTC 2012")
// - machine: the machine hardware name ("i686")
// - os:      the operating system name ("GNU/Linux")
type KernelGenerator struct {
}

// Key XXX
func (g *KernelGenerator) Key() string {
	return "kernel"
}

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate XXX
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
