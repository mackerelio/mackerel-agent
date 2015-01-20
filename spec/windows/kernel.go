// +build windows

package windows

import (
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
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
	results := make(map[string]string)

	name, err := windows.RegGetString(
		windows.HkeyLocalMachine,
		`Software\Microsoft\Windows NT\CurrentVersion`,
		`ProductName`)
	if err != nil {
		return nil, err
	}
	version, err := windows.RegGetString(
		windows.HkeyLocalMachine,
		`Software\Microsoft\Windows NT\CurrentVersion`,
		`CurrentVersion`)
	if err != nil {
		return nil, err
	}
	release, err := windows.RegGetString(
		windows.HkeyLocalMachine,
		`Software\Microsoft\Windows NT\CurrentVersion`,
		`CSDVersion`)
	if err != nil {
		return nil, err
	}

	results["name"] = "Microsoft Windows"
	results["os"] = name
	results["version"] = version
	results["release"] = release

	var systemInfo windows.SystemInfo
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
