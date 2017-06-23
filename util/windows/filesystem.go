// +build windows

package windows

import (
	"fmt"
	"strings"
	"syscall"
	"unsafe"

	"github.com/mackerelio/golib/logging"
)

// FilesystemInfo XXX
type FilesystemInfo struct {
	PercentUsed string
	KbUsed      float64
	KbSize      float64
	KbAvailable float64
	Mount       string
	Label       string
	VolumeName  string
	FsType      string
}

var windowsLogger = logging.GetLogger("windows")

// CollectFilesystemValues XXX
func CollectFilesystemValues() (map[string]FilesystemInfo, error) {
	filesystems := make(map[string]FilesystemInfo)

	drivebuf := make([]byte, 256)

	r, _, err := GetLogicalDriveStrings.Call(
		uintptr(len(drivebuf)),
		uintptr(unsafe.Pointer(&drivebuf[0])))
	if r == 0 {
		return nil, err
	}

	drives := []string{}
	for _, v := range drivebuf {
		if v >= 65 && v <= 90 {
			drive := string(v)
			r, _, _ = GetDriveType.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive + `:\`))))
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
		r, _, _ = GetDiskFreeSpaceEx.Call(
			uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(drive))),
			uintptr(unsafe.Pointer(&freeBytesAvailable)),
			uintptr(unsafe.Pointer(&totalNumberOfBytes)),
			0)
		if r == 0 {
			continue
		}
		filesystems[drive] = FilesystemInfo{
			PercentUsed: fmt.Sprintf("%d%%", 100*(totalNumberOfBytes-freeBytesAvailable)/totalNumberOfBytes),
			KbUsed:      float64((totalNumberOfBytes - freeBytesAvailable) / 1024),
			KbSize:      float64(totalNumberOfBytes / 1024),
			KbAvailable: float64(freeBytesAvailable / 1024),
			Mount:       drive,
			Label:       syscall.UTF16ToString(drivebuf),
			VolumeName:  syscall.UTF16ToString(volumebuf),
			FsType:      strings.ToLower(syscall.UTF16ToString(fsnamebuf)),
		}
	}

	return filesystems, nil
}
