package config

import (
	"os"
	"path/filepath"
)

func init() {
	mackerelRoot := filepath.Join(os.Getenv("HOME"), "Library", getAgentName())
	DefaultConfig = &Config{
		Apibase:    getApibase(),
		Root:       mackerelRoot,
		Pidfile:    filepath.Join(mackerelRoot, "pid"),
		Conffile:   filepath.Join(mackerelRoot, getAgentName()+".conf"),
		Connection: defaultConnectionConfig,
	}
}
