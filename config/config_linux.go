package config

var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     "/var/lib/mackerel-agent",
	Pidfile:  "/var/run/mackerel-agent.pid",
	Conffile: "/etc/mackerel-agent/mackerel-agent.conf",
	Roles:    []string{},
	Verbose:  false,
}
