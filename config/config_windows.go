package config

import (
	"github.com/mackerelio/mackerel-agent/util/windows"
	"log"
	"path/filepath"
)

func execdirInit() string {
	path, err := windows.ExecPath()
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Dir(path)
}

var execdir = execdirInit()

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
