package config

import (
	"os"
	"path/filepath"
)

var mackerelRoot = filepath.Join(os.Getenv("HOME"), "Library", getAgentName())

// DefaultConfig The default configuration for netbsd
var DefaultConfig = &Config{
	Apibase:    getApibase(),
	Root:       mackerelRoot,
	Conffile:   filepath.Join(mackerelRoot, getAgentName()+".conf"),
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
