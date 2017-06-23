// +build windows

package windows

import (
	"fmt"
	"unsafe"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// MemoryGenerator collects the host's memory specs.
type MemoryGenerator struct {
}

// Key XXX
func (g *MemoryGenerator) Key() string {
	return "memory"
}

var memoryLogger = logging.GetLogger("spec.memory")

// Generate XXX
func (g *MemoryGenerator) Generate() (interface{}, error) {
	result := make(map[string]interface{})

	var memoryStatusEx windows.MEMORY_STATUS_EX
	memoryStatusEx.Length = uint32(unsafe.Sizeof(memoryStatusEx))
	r, _, err := windows.GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memoryStatusEx)))
	if r == 0 {
		return nil, err
	}

	result["total"] = fmt.Sprintf("%dkb", memoryStatusEx.TotalPhys/1024)
	result["free"] = fmt.Sprintf("%dkb", memoryStatusEx.AvailPhys/1024)

	return result, nil
}
