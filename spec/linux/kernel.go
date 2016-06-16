// +build linux

package linux

import (
	"os/exec"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

// KernelGenerator XXX
type KernelGenerator struct {
}

// Key XXX
func (g *KernelGenerator) Key() string {
	return "kernel"
}

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate XXX
func (g *KernelGenerator) Generate() (interface{}, error) {
	commands := map[string][]string{
		"name":    {"uname", "-s"},
		"release": {"uname", "-r"},
		"version": {"uname", "-v"},
		"machine": {"uname", "-m"},
		"os":      {"uname", "-o"},
	}

	results := make(map[string]string)
	for key, command := range commands {
		out, err := exec.Command(command[0], command[1]).Output()
		if err != nil {
			kernelLogger.Errorf("Failed to run %s %s (skip this spec): %s", command[0], command[1], err)
			return nil, err
		}
		str := strings.TrimSpace(string(out))

		results[key] = str
	}

	return results, nil
}
