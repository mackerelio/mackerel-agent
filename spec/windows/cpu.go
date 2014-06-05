// +build windows

package windows

import (
	"fmt"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	. "github.com/mackerelio/mackerel-agent/util/windows"
)

type CPUGenerator struct {
}

func (g *CPUGenerator) Key() string {
	return "cpu"
}

var cpuLogger = logging.GetLogger("spec.cpu")

func (g *CPUGenerator) Generate() (interface{}, error) {
	results := make([]map[string]interface{}, 0)

	var systemInfo SYSTEM_INFO
	GetSystemInfo.Call(uintptr(unsafe.Pointer(&systemInfo)))

	for i := uint32(0); i < systemInfo.NumberOfProcessors; i++ {
		processorName, err := RegGetString(
			HKEY_LOCAL_MACHINE,
			fmt.Sprintf(`HARDWARE\DESCRIPTION\System\CentralProcessor\%d`, i),
			`ProcessorNameString`)
		if err != nil {
			return nil, err
		}
		processorMHz, err := RegGetInt(
			HKEY_LOCAL_MACHINE,
			fmt.Sprintf(`HARDWARE\DESCRIPTION\System\CentralProcessor\%d`, i),
			`~MHz`)
		if err != nil {
			return nil, err
		}
		vendorIdentifier, err := RegGetString(
			HKEY_LOCAL_MACHINE,
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
