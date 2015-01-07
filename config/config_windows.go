package config

var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     ".",
	Pidfile:  "mackerel-agent.pid",
	Conffile: "mackerel-agent.conf",
	Roles:    []string{},
	Verbose:  false,
	Connection: ConnectionConfig{
		PostMetricsDequeueDelaySeconds: 30,
		PostMetricsRetryDelaySeconds:   60,
		PostMetricsRetryMax:            10,
		PostMetricsBufferSize:          30,
	},
}
