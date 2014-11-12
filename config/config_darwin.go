package config

import "path/filepath"
import "os"

var mackerelRoot = filepath.Join(os.Getenv("HOME"), "Library", "mackerel-agent")

// The default configuration for dawrin.
var DefaultConfig = &Config{
	Apibase:  "https://mackerel.io",
	Root:     mackerelRoot,
	Pidfile:  filepath.Join(mackerelRoot, "pid"),
	Conffile: filepath.Join(mackerelRoot, "mackerel-agent.conf"),
	Roles:    []string{},
	Verbose:  false,
	Connection: ConnectionConfig{
		Post_Metrics_Dequeue_Delay_Seconds: 30,     // Check the metric values queue for each half minutes
		Post_Metrics_Retry_Delay_Seconds:   60,     // Wait a minute before retrying metric value posts
		Post_Metrics_Retry_Max:             60,     // Retry up to 60 times (5min * 60 = 5hrs)
		Post_Metrics_Buffer_Size:           6 * 60, // Keep metric values of 6 hours span in the queue
	},
}
