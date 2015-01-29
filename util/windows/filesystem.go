// +build windows

package windows

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"
	"github.com/mackerelio/mackerel-agent/logging"
)

var windowsLogger = logging.GetLogger("windows")

// CollectDfValues XXX
func CollectDfValues() (map[string]map[string]interface{}, error) {
	filesystems := make(map[string]map[string]interface{})

	drivebuf := make([]byte, 256)

	r, _, err := GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))
	if r == 0 {
		e := GetLastError.Call
		windowsLogger.Errorf("error code is  [%q]", e)
		return nil, err
	}

	drives := []string{}
	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			r, _, err = GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
			if r != DRIVE_FIXED {
				continue
			}
			drives = append(drives, drive+":")
		}
	}

	for _, drive := range drives {
		drivebuf := make([]uint16, 256)
		r, _, err := QueryDosDevice.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&drivebuf[0])),
			uintptr(len(drivebuf)))
		if r == 0 {
			windowsLogger.Debugf("do not get DosDevice [%q]", drivebuf)
			return nil, err
		}
		volumebuf := make([]uint16, 256)
		fsnamebuf := make([]uint16, 256)
		r, _, err = GetVolumeInformationW.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive+`\`))),
			uintptr(unsafe.Pointer(&volumebuf[0])),
			uintptr(len(volumebuf)),
			0,
			0,
			0,
			uintptr(unsafe.Pointer(&fsnamebuf[0])),
			uintptr(len(fsnamebuf)))
		if r == 0 {
			windowsLogger.Debugf("do not get volume [%q] or fsname [%q]", volumebuf, fsnamebuf)
			return nil, err
		}
		freeBytesAvailable := int64(0)
		totalNumberOfBytes := int64(0)
		r, _, err = GetDiskFreeSpaceEx.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&freeBytesAvailable)),
			uintptr(unsafe.Pointer(&totalNumberOfBytes)),
			0)
		if r == 0 {
			continue
		}
		filesystems[drive] = map[string]interface{}{
			"percent_used": fmt.Sprintf("%d%%", 100*(totalNumberOfBytes-freeBytesAvailable)/totalNumberOfBytes),
			"kb_used":      float64((totalNumberOfBytes - freeBytesAvailable) / 1024),
			"kb_size":      float64(totalNumberOfBytes / 1024),
			"kb_available": float64(freeBytesAvailable / 1024),
			"mount":        drive,
			"label":        syscall.UTF16ToString(drivebuf),
			"volume_name":  syscall.UTF16ToString(volumebuf),
			"fs_type":      strings.ToLower(syscall.UTF16ToString(fsnamebuf)),
		}
	}

	return filesystems, nil
}
