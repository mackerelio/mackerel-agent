//go:build darwin

package darwin

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
)

// KernelGenerator Generates specs about the kernel.
type KernelGenerator struct {
}

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate collects specs from `uname` command and `sw_vers` command
func (g *KernelGenerator) Generate() (any, error) {
	unameCommand := "/usr/bin/uname"
	swVersCommand := "/usr/bin/sw_vers"

	commands := map[string][]string{
		"release":          {unameCommand, "-r"},
		"version":          {unameCommand, "-v"},
		"machine":          {unameCommand, "-m"},
		"os":               {unameCommand, "-s"},
		"platform_name":    {swVersCommand, "-productName"},
		"platform_version": {swVersCommand, "-productVersion"},
	}

	// +1 is for `name`
	results := make(mackerel.Kernel, len(commands)+1)

	for field, commandAndArgs := range commands {
		out, err := exec.Command(commandAndArgs[0], commandAndArgs[1:]...).Output()
		if err != nil {
			kernelLogger.Errorf("Failed to run %s (skip this field): %s", commandAndArgs, err)
			continue
		}

		results[field] = strings.TrimSpace(string(out))
	}

	results["name"] = results["os"]

	return results, nil
}
