// +build windows

package windows

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// Key XX
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var filesystemLogger = logging.GetLogger("spec.filesystem")

// Generate XXX
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	filesystems := make(map[string]map[string]interface{})

	drivebuf := make([]byte, 256)
	_, r, err := windows.GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))
	if r != 0 {
		return nil, err
	}

	drives := []string{}
	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			r, _, err = windows.GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
			if r != windows.DriveFixed {
				continue
			}
			drives = append(drives, drive+":")
		}
	}

	for _, drive := range drives {
		drivebuf := make([]uint16, 256)
		r, _, err := windows.QueryDosDevice.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&drivebuf[0])),
			uintptr(len(drivebuf)))
		if r == 0 {
			return nil, err
		}
		volumebuf := make([]uint16, 256)
		fsnamebuf := make([]uint16, 256)
		r, _, err = windows.GetVolumeInformationW.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive+`\`))),
			uintptr(unsafe.Pointer(&volumebuf[0])),
			uintptr(len(volumebuf)),
			0,
			0,
			0,
			uintptr(unsafe.Pointer(&fsnamebuf[0])),
			uintptr(len(fsnamebuf)))
		if r == 0 {
			return nil, err
		}
		freeBytesAvailable := int64(0)
		totalNumberOfBytes := int64(0)
		r, _, err = windows.GetDiskFreeSpaceEx.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&freeBytesAvailable)),
			uintptr(unsafe.Pointer(&totalNumberOfBytes)),
			0)
		if r == 0 {
			continue
		}
		filesystems[drive] = map[string]interface{}{
			"percent_used": fmt.Sprintf("%d%%", 100*(totalNumberOfBytes-freeBytesAvailable)/totalNumberOfBytes),
			"kb_used":      (totalNumberOfBytes - freeBytesAvailable) / 1024 / 1024,
			"kb_size":      totalNumberOfBytes / 1024 / 1024,
			"kb_available": freeBytesAvailable / 1024 / 1024,
			"mount":        drive,
			"label":        syscall.UTF16ToString(drivebuf),
			"volume_name":  syscall.UTF16ToString(volumebuf),
			"fs_type":      strings.ToLower(syscall.UTF16ToString(fsnamebuf)),
		}
	}

	return filesystems, nil
}
