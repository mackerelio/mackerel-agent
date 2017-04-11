package config

import (
	"log"
	"path/filepath"

	"github.com/mackerelio/mackerel-agent/util/windows"
)

func init() {
	path, err := windows.ExecPath()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(path)
	agentName := getAgentName()
	DefaultConfig = &Config{
		Apibase:    getApibase(),
		Root:       execdir,
		Pidfile:    filepath.Join(execdir, agentName+".pid"),
		Conffile:   filepath.Join(execdir, agentName+".conf"),
		Connection: defaultConnectionConfig,
	}
}
