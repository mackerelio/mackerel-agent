package config

var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     "/var/lib/mackerel-agent",
	Pidfile:  "/var/run/mackerel-agent.pid",
	Conffile: "/etc/mackerel-agent/mackerel-agent.conf",
	Roles:    []string{},
	Verbose:  false,
	Metrics: MetricsConfig{
		Dequeue_Delay: 30,
		Retry_Delay:   60,
		Retry_Max:     10,
		Buffer_Size:   30,
	},
}
