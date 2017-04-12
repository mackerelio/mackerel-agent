// +build linux freebsd netbsd

package config

import "fmt"

func init() {
	agentName := getAgentName()
	DefaultConfig = &Config{
		Apibase:    getApibase(),
		Root:       fmt.Sprintf("/var/lib/%s", agentName),
		Pidfile:    fmt.Sprintf("/var/run/%s.pid", agentName),
		Conffile:   fmt.Sprintf("/etc/%[1]s/%[1]s.conf", agentName),
		Connection: defaultConnectionConfig,
	}
}
