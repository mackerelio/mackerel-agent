// +build netbsd

package netbsd

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
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

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate XXX
func (g *KernelGenerator) Generate() (interface{}, error) {
	unameArgs := map[string][]string{
		"release": {"-r"},
		"version": {"-v"},
		"machine": {"-m"},
		"os":      {"-s"},
	}

	results := make(mackerel.Kernel, len(unameArgs)+1)

	for field, args := range unameArgs {
		out, err := exec.Command("/usr/bin/uname", args...).Output()
		if err != nil {
			kernelLogger.Errorf("Failed to run uname %s (skip this field): %s", args, err)
			continue
		}

		results[field] = strings.TrimSpace(string(out))
	}

	results["name"] = results["os"]

	return results, nil
}
