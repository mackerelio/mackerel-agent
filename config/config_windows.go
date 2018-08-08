package config

import (
	"log"
	"os"
	"path/filepath"
)

func init() {
	path, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	execDir := filepath.Dir(path)
	agentName := getAgentName()
	DefaultConfig = &Config{
		Apibase:  getApibase(),
		Root:     execDir,
		Pidfile:  filepath.Join(execDir, agentName+".pid"),
		Conffile: filepath.Join(execDir, agentName+".conf"),
	}
}
