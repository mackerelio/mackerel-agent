// +build windows

package windows

import (
	"syscall"
	"unsafe"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// BlockDeviceGenerator XXX
type BlockDeviceGenerator struct {
}

// Key XXX
func (g *BlockDeviceGenerator) Key() string {
	return "block_device"
}

var blockDeviceLogger = logging.GetLogger("spec.block_device")

// Generate XXX
func (g *BlockDeviceGenerator) Generate() (interface{}, error) {
	results := make(map[string]map[string]interface{})

	drivebuf := make([]byte, 256)
	r, _, err := windows.GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))
	if r == 0 {
		return nil, err
	}

	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			removable := false
			r, _, err = windows.GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
			if r == windows.DRIVE_REMOVABLE {
				removable = true
			}
			freeBytesAvailable := int64(0)
			totalNumberOfBytes := int64(0)
			totalNumberOfFreeBytes := int64(0)
			r, _, _ = windows.GetDiskFreeSpaceEx.Call(
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
