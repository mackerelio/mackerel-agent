package config

var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     "/var/lib/mackerel-agent",
	Pidfile:  "/var/run/mackerel-agent.pid",
	Conffile: "/etc/mackerel-agent/mackerel-agent.conf",
	Roles:    []string{},
	Verbose:  false,
	Connection: ConnectionConfig{
		Post_Metrics_Dequeue_Delay_Seconds: 30,     // Check the metric values queue for each half minutes
		Post_Metrics_Retry_Delay_Seconds:   5 * 60, // Wait 5 minutes before retrying metric value posts
		Post_Metrics_Retry_Max:             60,     // Retry up to 60 times (5min * 60 = 5hrs)
		Post_Metrics_Buffer_Size:           6 * 60, // Keep metric values of 6 hours span in the queue
	},
}
