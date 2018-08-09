// +build windows

package windows

import (
	"fmt"
	"unsafe"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/util/windows"
)

// CPUGenerator collects CPU specs
type CPUGenerator struct {
}

var cpuLogger = logging.GetLogger("spec.cpu")

// Generate collects CPU specs.
func (g *CPUGenerator) Generate() (interface{}, error) {
	var results mackerel.CPU

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
