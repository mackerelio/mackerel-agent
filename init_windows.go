// +build windows

package main

import (
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

var (
	kernel32              = syscall.NewLazyDLL("kernel32")
	procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")
)

func getModuleFileName() (string, error) {
	var p [syscall.MAX_PATH]uint16
	result, _, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&p[0])), uintptr(len(p)))
	if result == 0 {
		return os.Args[0], err
	}
	return syscall.UTF16ToString(p[:]), nil
}

func init() {
	if p, err := getModuleFileName(); err == nil {
		os.Setenv("PATH", filepath.Dir(p)+
			string(filepath.ListSeparator)+os.Getenv("PATH"))
	}
}
