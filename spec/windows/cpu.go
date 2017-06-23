// +build windows

package windows

import (
	"fmt"
	"unsafe"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// CPUGenerator XXX
type CPUGenerator struct {
}

// Key XXX
func (g *CPUGenerator) Key() string {
	return "cpu"
}

var cpuLogger = logging.GetLogger("spec.cpu")

// Generate XXX
func (g *CPUGenerator) Generate() (interface{}, error) {
	var results []map[string]interface{}

	var systemInfo windows.SYSTEM_INFO
	windows.GetSystemInfo.Call(uintptr(unsafe.Pointer(&systemInfo)))

	for i := uint32(0); i < systemInfo.NumberOfProcessors; i++ {
		processorName, _, err := windows.RegGetString(
			windows.HKEY_LOCAL_MACHINE,
			fmt.Sprintf(`HARDWARE\DESCRIPTION\System\CentralProcessor\%d`, i),
			`ProcessorNameString`)
		if err != nil {
			return nil, err
		}
		processorMHz, _, err := windows.RegGetInt(
			windows.HKEY_LOCAL_MACHINE,
			fmt.Sprintf(`HARDWARE\DESCRIPTION\System\CentralProcessor\%d`, i),
			`~MHz`)
		if err != nil {
			return nil, err
		}
		vendorIdentifier, _, err := windows.RegGetString(
			windows.HKEY_LOCAL_MACHINE,
			fmt.Sprintf(`HARDWARE\DESCRIPTION\System\CentralProcessor\%d`, i),
			`VendorIdentifier`)
		if err != nil {
			return nil, err
		}
		results = append(results, map[string]interface{}{
			"model_name": processorName,
			"mhz":        processorMHz,
			"model":      systemInfo.ProcessorArchitecture,
			"vendor_id":  vendorIdentifier,
		})
	}
	return results, nil
}
