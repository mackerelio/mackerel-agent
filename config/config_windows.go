package config

import (
	"log"
	"path/filepath"
	"syscall"
	"unsafe"
)

func execdirInit() string {
	var (
		kernel32              = syscall.NewLazyDLL("kernel32")
		procGetModuleFileName = kernel32.NewProc("GetModuleFileNameW")
	)
	var wpath [syscall.MAX_PATH]uint16
	r1, _, err := procGetModuleFileName.Call(0, uintptr(unsafe.Pointer(&wpath[0])), uintptr(len(wpath)))
	if r1 == 0 {
		log.Fatal(err)
	}
	return syscall.UTF16ToString(wpath[:])
}

var execdir = filepath.Dir(execdirInit())

// The default configuration for windows
var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     execdir,
	Pidfile:  filepath.Join(execdir, "mackerel-agent.pid"),
	Conffile: filepath.Join(execdir, "mackerel-agent.conf"),
	Roles:    []string{},
	Verbose:  false,
	Connection: ConnectionConfig{
		PostMetricsDequeueDelaySeconds: 30,
		PostMetricsRetryDelaySeconds:   60,
		PostMetricsRetryMax:            10,
		PostMetricsBufferSize:          30,
	},
}
