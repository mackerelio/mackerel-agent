//go:build windows
// +build windows

package windows

import (
	"syscall"
	"unsafe"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/util/windows"
)

// BlockDeviceGenerator XXX
type BlockDeviceGenerator struct {
}

// Generate XXX
func (g *BlockDeviceGenerator) Generate() (interface{}, error) {
	results := make(mackerel.BlockDevice)

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
			d, err := syscall.UTF16PtrFromString(drive + `:\`)
			if err != nil {
				return nil, err
			}
			r, _, _ = windows.GetDriveType.Call(uintptr(unsafe.Pointer(d)))
			if r == windows.DRIVE_REMOVABLE {
				removable = true
			}
			freeBytesAvailable := int64(0)
			totalNumberOfBytes := int64(0)
			totalNumberOfFreeBytes := int64(0)
			d, err = syscall.UTF16PtrFromString(drive)
			if err != nil {
				return nil, err
			}
			r, _, _ = windows.GetDiskFreeSpaceEx.Call(
				uintptr(unsafe.Pointer(d)),
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
