package config

var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     "/var/lib/mackerel-agent",
	Pidfile:  "/var/run/mackerel-agent.pid",
	Conffile: "/etc/mackerel-agent/mackerel-agent.conf",
	Roles:    []string{},
	Verbose:  false,
	Connection: ConnectionConfig{
		Metrics_Dequeue_Delay: 30,
		Metrics_Retry_Delay:   60,
		Metrics_Retry_Max:     10,
		Metrics_Buffer_Size:   30,
	},
}
