// +build windows

package windows

import (
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// MemoryGenerator XXX
type MemoryGenerator struct {
}

var memoryLogger = logging.GetLogger("metrics.memory")

// NewMemoryGenerator XXX
func NewMemoryGenerator() (*MemoryGenerator, error) {
	return &MemoryGenerator{}, nil
}

// Generate XXX
func (g *MemoryGenerator) Generate() (metrics.Values, error) {
	ret := make(map[string]float64)

	var memoryStatusEx windows.MemoryStatusEx
	memoryStatusEx.Length = uint32(unsafe.Sizeof(memoryStatusEx))
	r, _, err := windows.GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memoryStatusEx)))
	if r == 0 {
		return nil, err
	}

	ret["memory.total"] = float64(memoryStatusEx.TotalPhys) / 1024
	ret["memory.free"] = float64(memoryStatusEx.AvailPhys) / 1024
	ret["memory.swap_free"] = 0
	ret["memory.swap_total"] = 0

	return metrics.Values(ret), nil
}
