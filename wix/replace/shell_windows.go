package main

import (
	"path/filepath"
	"syscall"
	"unsafe"
)

const (
	// (snip) the maximum length for a path is MAX_PATH, which is defined as 260 characters.
	// https://docs.microsoft.com/en-us/windows/win32/fileio/naming-a-file#maximum-path-length-limitation
	maxPathLen = 260

	// This value is defined in shlobj.h, but it's published only name.
	// But its value is used, such as JNI, so probably it isn't changed.
	// https://docs.microsoft.com/en-us/windows/win32/shell/csidl
	csidlProgramFilesX86 = 0x2a
)

var (
	shell32                     = syscall.NewLazyDLL("shell32")
	procSHGetSpecialFolderPathW = shell32.NewProc("SHGetSpecialFolderPathW")
)

// FallbackConfigDir returns an other ProgramFiles location if there.
// wdir is current installed folder of the mackerel-agent.exe.
func FallbackConfigDir(wdir string) string {
	// The host that has installed mackerel-agent x86 edition is having id, pid and mackerel-agent.conf files
	// in C:\Program Files (x86)\Mackerel.
	// Though mackerel-agent refers to a directory containing mackerel-agent.exe,
	// it should still refer these old files in x86 folder even if mackerel-agent is upgraded to x64 edition.

	dir := getProgramFilesX86()
	if dir == "" {
		return ""
	}
	dir = filepath.Join(dir, "Mackerel", "mackerel-agent")
	if dir == wdir {
		return ""
	}
	return dir
}

func getProgramFilesX86() string {
	var buf [maxPathLen]uint16
	p := unsafe.Pointer(&buf[0])
	rv, _, _ := procSHGetSpecialFolderPathW.Call(0, uintptr(p), csidlProgramFilesX86, 0)
	if rv == 0 {
		return ""
	}
	return syscall.UTF16ToString(buf[:])
}
