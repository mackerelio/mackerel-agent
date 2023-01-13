//go:build windows
// +build windows

package windows

import (
	"strings"
	"unsafe"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/util/windows"
	"github.com/yusufpapurcu/wmi"
)

const registryKey = `Software\Microsoft\Windows NT\CurrentVersion`

type Win32_OperatingSystem struct {
	Caption    string
	Version    string
	CSDVersion string
}

// KernelGenerator XXX
type KernelGenerator struct {
}

// Generate XXX
func (g *KernelGenerator) Generate() (interface{}, error) {
	results := make(mackerel.Kernel)

	var dst []Win32_OperatingSystem
	var osname, version, release string
	q := wmi.CreateQuery(&dst, "")
	if err := wmi.Query(q, &dst); err != nil {
		return nil, err
	}
	for _, v := range dst {
		osname = v.Caption
		version = v.Version
		release = v.CSDVersion
		break //nolint
	}

	edition, _, err := windows.RegGetString(
		windows.HKEY_LOCAL_MACHINE, registryKey, `EditionID`)
	if err != nil {
		return nil, err
	}

	if edition != "" && !strings.Contains(osname, edition) {
		osname += " (" + edition + ")"
	}

	results["name"] = "Microsoft Windows"
	results["os"] = osname
	results["version"] = version
	results["release"] = release

	var systemInfo windows.SYSTEM_INFO
	_, _, _ = windows.GetSystemInfo.Call(uintptr(unsafe.Pointer(&systemInfo)))
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
