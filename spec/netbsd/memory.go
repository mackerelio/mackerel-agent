//go:build netbsd
// +build netbsd

package netbsd

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
)

// MemoryGenerator collects the host's memory specs.
type MemoryGenerator struct {
}

var memoryLogger = logging.GetLogger("spec.memory")

const bytesInKibibytes = 1024

// Generate returns memory specs.
// The returned spec must have below:
// - total (in "###kB" format, Kibibytes)
func (g *MemoryGenerator) Generate() (any, error) {
	spec := make(mackerel.Memory)

	cmd := exec.Command("sysctl", "-n", "hw.physmem64")
	outputBytes, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("sysctl -n hw.physmem64: %s", err)
	}

	output := string(outputBytes)

	memsizeInBytes, err := strconv.ParseInt(strings.TrimSpace(output), 10, 64)
	memoryLogger.Debugf("memsizeInBytes: %d", memsizeInBytes)
	if err != nil {
		memoryLogger.Debugf("MemoryGenerator err != nil")
		return nil, fmt.Errorf("while parsing %q: %s", output, err)
	}

	spec["total"] = fmt.Sprintf("%dkB", memsizeInBytes/bytesInKibibytes)
	memoryLogger.Debugf("spec[total]: %s", spec["total"])

	return spec, nil
}
