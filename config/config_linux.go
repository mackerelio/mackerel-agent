package config

import "fmt"

// DefaultConfig The default configuration for linux
var DefaultConfig = &Config{
	Apibase:    getApibase(),
	Root:       fmt.Sprintf("/var/lib/%s", getAgentName()),
	Pidfile:    fmt.Sprintf("/var/run/%s.pid", getAgentName()),
	Conffile:   fmt.Sprintf("/etc/%s/%s.conf", getAgentName(), getAgentName()),
	Roles:      []string{},
	Verbose:    false,
	Diagnostic: false,
	Connection: ConnectionConfig{
		PostMetricsDequeueDelaySeconds: 30,     // Check the metric values queue for every half minute
		PostMetricsRetryDelaySeconds:   60,     // Wait a minute before retrying metric value posts
		PostMetricsRetryMax:            60,     // Retry up to 60 times (30s * 60 = 30min)
		PostMetricsBufferSize:          6 * 60, // Keep metric values of 6 hours span in the queue
	},
}
