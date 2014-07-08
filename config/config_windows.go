package config

var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     ".",
	Pidfile:  "mackerel-agent.pid",
	Conffile: "mackerel-agent.conf",
	Roles:    []string{},
	Verbose:  false,
	Connection: ConnectionConfig{
		Post_Metrics_Dequeue_Delay_Seconds: 30,
		Post_Metrics_Retry_Delay_Seconds:   60,
		Post_Metrics_Retry_Max:             10,
		Post_Metrics_Buffer_Size:           30,
	},
}
