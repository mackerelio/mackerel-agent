// +build windows

package windows

import (
	"strings"
	"unsafe"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/util/windows"
)

const registryKey = `Software\Microsoft\Windows NT\CurrentVersion`

// KernelGenerator XXX
type KernelGenerator struct {
}

var kernelLogger = logging.GetLogger("spec.kernel")

// Generate XXX
func (g *KernelGenerator) Generate() (interface{}, error) {
	results := make(mackerel.Kernel)

	osname, _, err := windows.RegGetString(
		windows.HKEY_LOCAL_MACHINE, registryKey, `ProductName`)
	if err != nil {
		return nil, err
	}
	edition, _, err := windows.RegGetString(
		windows.HKEY_LOCAL_MACHINE, registryKey, `EditionID`)
	if err != nil {
		return nil, err
	}
	version, _, err := windows.RegGetString(
		windows.HKEY_LOCAL_MACHINE, registryKey, `CurrentVersion`)
	if err != nil {
		return nil, err
	}
	release, errno, err := windows.RegGetString(
		windows.HKEY_LOCAL_MACHINE, registryKey, `CSDVersion`)
	if err != nil && errno != windows.ERROR_FILE_NOT_FOUND { // CSDVersion is nullable
		return nil, err
	}

	if edition != "" && strings.Index(osname, edition) == -1 {
		osname += " (" + edition + ")"
	}

	results["name"] = "Microsoft Windows"
	results["os"] = osname
	results["version"] = version
	results["release"] = release

	var systemInfo windows.SYSTEM_INFO
	windows.GetSystemInfo.Call(uintptr(unsafe.Pointer(&systemInfo)))
	switch systemInfo.ProcessorArchitecture {
	case 0:
		results["machine"] = "x86"
	case 1:
		results["machine"] = "mips"
	case 2:
		results["machine"] = "alpha"
	case 3:
		results["machine"] = "ppc"
	case 4:
		results["machine"] = "shx"
	case 5:
		results["machine"] = "arm"
	case 6:
		results["machine"] = "ia64"
	case 7:
		results["machine"] = "alpha64"
	case 8:
		results["machine"] = "msil"
	case 9:
		results["machine"] = "amd64"
	case 10:
		results["machine"] = "ia32_on_win64"
	}

	return results, nil
}
