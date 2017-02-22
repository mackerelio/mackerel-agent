// +build netbsd

package netbsd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

// MemoryGenerator collects the host's memory specs.
type MemoryGenerator struct {
}

// Key XXX
func (g *MemoryGenerator) Key() string {
	return "memory"
}

var memoryLogger = logging.GetLogger("spec.memory")

const bytesInKibibytes = 1024

// Generate returns memory specs.
// The returned spec must have below:
// - total (in "###kB" format, Kibibytes)
func (g *MemoryGenerator) Generate() (interface{}, error) {
	spec := map[string]string{}

	cmd := exec.Command("sysctl", "-n", "hw.physmem64")
	outputBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("sysctl -n hw.physmem64: %s", err)
	}

	output := string(outputBytes)

	memsizeInBytes, err := strconv.ParseInt(strings.TrimSpace(output), 10, 64)
	fmt.Printf("[DEBUG] memsizeInBytes: %d", memsizeInBytes)
	if err != nil {
		fmt.Printf("[DEBUG] MemoryGenerator err != nil")
		return nil, fmt.Errorf("while parsing %q: %s", output, err)
	}

	spec["total"] = fmt.Sprintf("%dkB", memsizeInBytes/bytesInKibibytes)
	fmt.Printf("[DEBUG] spec[total]: %s", spec["total"])

	return spec, nil
}
