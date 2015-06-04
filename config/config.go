package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mackerelio/mackerel-agent/logging"
)

var configLogger = logging.GetLogger("config")

var apibase string

func getApibase() string {
	if apibase != "" {
		return apibase
	}
	return "https://mackerel.io"
}

// Config represents mackerel-agent's configuration file.
type Config struct {
	Apibase     string
	Apikey      string
	Root        string
	Pidfile     string
	Conffile    string
	Roles       []string
	Verbose     bool
	Connection  ConnectionConfig
	DisplayName string     `toml:"display_name"`
	HostStatus  HostStatus `toml:"host_status"`

	// Corresponds to the set of [plugin.<kind>.<name>] sections
	// the key of the map is <kind>, which should be one of "metrics" or "checks".
	Plugin map[string]PluginConfigs

	DeprecatedSensu map[string]PluginConfigs `toml:"sensu"` // DEPRECATED this is for backward compatibility
	Include         string
}

// PluginConfigs represents a set of [plugin.<kind>.<name>] sections in the configuration file
// under a specific <kind>. The key of the map is <name>, for example "mysql" of "plugin.metrics.mysql".
type PluginConfigs map[string]PluginConfig

// PluginConfig represents a section of [plugin.*].
type PluginConfig struct {
	Command string
}

const postMetricsDequeueDelaySecondsMax = 59   // max delay seconds for dequeuing from buffer queue
const postMetricsRetryDelaySecondsMax = 3 * 60 // max delay seconds for retrying a request that caused errors

// PostMetricsInterval XXX
var PostMetricsInterval = 1 * time.Minute

// ConnectionConfig XXX
type ConnectionConfig struct {
	PostMetricsDequeueDelaySeconds int `toml:"post_metrics_dequeue_delay_seconds"` // delay for dequeuing from buffer queue
	PostMetricsRetryDelaySeconds   int `toml:"post_metrics_retry_delay_seconds"`   // delay for retrying a request that caused errors
	PostMetricsRetryMax            int `toml:"post_metrics_retry_max"`             // max numbers of retries for a request that causes errors
	PostMetricsBufferSize          int `toml:"post_metrics_buffer_size"`           // max numbers of requests stored in buffer queue.
}

type HostStatus struct {
	Start string
	Stop  string
}

// CheckNames return list of plugin.checks._name_
func (conf *Config) CheckNames() []string {
	checks := []string{}
	for name := range conf.Plugin["checks"] {
		checks = append(checks, name)
	}
	return checks
}

// LoadConfig XXX
func LoadConfig(conffile string) (*Config, error) {
	config, err := loadConfigFile(conffile)

	// set default values if config does not have values
	if config.Apibase == "" {
		config.Apibase = DefaultConfig.Apibase
	}
	if config.Root == "" {
		config.Root = DefaultConfig.Root
	}
	if config.Pidfile == "" {
		config.Pidfile = DefaultConfig.Pidfile
	}
	if config.Verbose == false {
		config.Verbose = DefaultConfig.Verbose
	}
	if config.Connection.PostMetricsDequeueDelaySeconds == 0 {
		config.Connection.PostMetricsDequeueDelaySeconds = DefaultConfig.Connection.PostMetricsDequeueDelaySeconds
	}
	if config.Connection.PostMetricsDequeueDelaySeconds > postMetricsDequeueDelaySecondsMax {
		configLogger.Warningf("'post_metrics_dequese_delay_seconds' is set to %d (Maximum Value).", postMetricsDequeueDelaySecondsMax)
		config.Connection.PostMetricsDequeueDelaySeconds = postMetricsDequeueDelaySecondsMax
	}
	if config.Connection.PostMetricsRetryDelaySeconds == 0 {
		config.Connection.PostMetricsRetryDelaySeconds = DefaultConfig.Connection.PostMetricsRetryDelaySeconds
	}
	if config.Connection.PostMetricsRetryDelaySeconds > postMetricsRetryDelaySecondsMax {
		configLogger.Warningf("'post_metrics_retry_delay_seconds' is set to %d (Maximum Value).", postMetricsRetryDelaySecondsMax)
		config.Connection.PostMetricsRetryDelaySeconds = postMetricsRetryDelaySecondsMax
	}
	if config.Connection.PostMetricsRetryMax == 0 {
		config.Connection.PostMetricsRetryMax = DefaultConfig.Connection.PostMetricsRetryMax
	}
	if config.Connection.PostMetricsBufferSize == 0 {
		config.Connection.PostMetricsBufferSize = DefaultConfig.Connection.PostMetricsBufferSize
	}

	return config, err
}

func loadConfigFile(file string) (*Config, error) {
	config := &Config{}
	if _, err := toml.DecodeFile(file, config); err != nil {
		return config, err
	}

	if config.Include != "" {
		if err := includeConfigFile(config, config.Include); err != nil {
			return config, err
		}
	}

	// for backward compatibility
	// merges sensu configs to plugin configs
	if _, ok := config.DeprecatedSensu["checks"]; ok {
		configLogger.Warningf("'sensu.checks.*' config format is DEPRECATED. Please use 'plugin.metrics.*' format.")

		if config.Plugin == nil {
			config.Plugin = map[string]PluginConfigs{}
		}
		if _, ok := config.Plugin["metrics"]; !ok {
			config.Plugin["metrics"] = PluginConfigs{}
		}
		for k, v := range config.DeprecatedSensu["checks"] {
			config.Plugin["metrics"]["DEPRECATED-sensu-"+k] = v
		}
	}

	return config, nil
}

func includeConfigFile(config *Config, include string) error {
	files, err := filepath.Glob(include)
	if err != nil {
		return err
	}

	for _, file := range files {
		// Save current "roles" value and reset it
		// because toml.DecodeFile()-ing on a fulfilled struct
		// produces bizarre array values.
		rolesSaved := config.Roles
		config.Roles = nil

		// Also, save plugin values for later merging
		pluginSaved := map[string]PluginConfigs{}
		for kind, plugins := range config.Plugin {
			pluginSaved[kind] = plugins
		}

		meta, err := toml.DecodeFile(file, &config)
		if err != nil {
			return fmt.Errorf("while loading included config file %s: %s", file, err)
		}

		// If included config does not have "roles" key,
		// use the previous roles configuration value.
		if meta.IsDefined("roles") == false {
			config.Roles = rolesSaved
		}

		for kind, plugins := range config.Plugin {
			for key, conf := range plugins {
				if pluginSaved[kind] == nil {
					pluginSaved[kind] = PluginConfigs{}
				}
				pluginSaved[kind][key] = conf
			}
		}

		config.Plugin = pluginSaved
	}

	return nil
}
