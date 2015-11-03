// +build netbsd

package netbsd

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
	unameArgs := map[string][]string{
		"release": []string{"-r"},
		"version": []string{"-v"},
		"machine": []string{"-m"},
		"os":      []string{"-s"},
	}

	results := make(map[string]string, len(unameArgs)+1)

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
