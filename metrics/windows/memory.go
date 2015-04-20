// +build windows

package windows

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
	"unsafe"
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

	var memoryStatusEx windows.MEMORY_STATUS_EX
	memoryStatusEx.Length = uint32(unsafe.Sizeof(memoryStatusEx))
	r, _, err := windows.GlobalMemoryStatusEx.Call(uintptr(unsafe.Pointer(&memoryStatusEx)))
	if r == 0 {
		return nil, err
	}

	free := float64(memoryStatusEx.AvailPhys)
	total := float64(memoryStatusEx.TotalPhys)
	ret["memory.free"] = free
	ret["memory.total"] = total
	ret["memory.used"] = total - free
	ret["memory.pagefile_total"] = float64(memoryStatusEx.TotalVirtual) / 1024
	ret["memory.pagefile_free"] = float64(memoryStatusEx.AvailVirtual) / 1024

	memoryLogger.Debugf("memory : %s", ret)
	return metrics.Values(ret), nil
}
