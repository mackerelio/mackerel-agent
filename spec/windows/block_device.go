// +build windows

package windows

import (
	"syscall"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	. "github.com/mackerelio/mackerel-agent/util/windows"
)

type BlockDeviceGenerator struct {
}

func (g *BlockDeviceGenerator) Key() string {
	return "block_device"
}

var blockDeviceLogger = logging.GetLogger("spec.block_device")

func (g *BlockDeviceGenerator) Generate() (interface{}, error) {
	results := make(map[string]map[string]interface{})

	drivebuf := make([]byte, 256)
	_, r, err := GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))
	if r != 0 {
		return nil, err
	}

	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			removable := false
			r, _, err = GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
			if r == DRIVE_REMOVABLE {
				removable = true
			}
			freeBytesAvailable := int64(0)
			totalNumberOfBytes := int64(0)
			totalNumberOfFreeBytes := int64(0)
			r, _, err = GetDiskFreeSpaceEx.Call(
				uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
				uintptr(unsafe.Pointer(&freeBytesAvailable)),
				uintptr(unsafe.Pointer(&totalNumberOfBytes)),
				uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)))
			if r == 0 {
				continue
			}
			results[drive] = map[string]interface{}{
				"size":      totalNumberOfFreeBytes,
				"removable": removable,
			}
		}
	}

	return results, nil
}
