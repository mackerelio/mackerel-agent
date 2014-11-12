package config

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/mackerelio/mackerel-agent/logging"
)

var configLogger = logging.GetLogger("config")

type Config struct {
	Apibase         string
	Apikey          string
	Root            string
	Pidfile         string
	Conffile        string
	Roles           []string
	Verbose         bool
	Connection      ConnectionConfig
	Plugin          map[string]PluginConfigs
	DeprecatedSensu map[string]PluginConfigs `toml:"sensu"` // DEPRECATED this is for backward compatibility
	Include         string
}

type PluginConfigs map[string]PluginConfig

type PluginConfig struct {
	Command string
}

const POST_METRICS_DEQUEUE_DELAY_SECONDS_MAX = 59   // max delay seconds for dequeuing from buffer queue
const POST_METRICS_RETRY_DELAY_SECONDS_MAX = 3 * 60 // max delay seconds for retrying a request that caused errors

var PostMetricsInterval = 1 * time.Minute

type ConnectionConfig struct {
	Post_Metrics_Dequeue_Delay_Seconds int // delay for dequeuing from buffer queue
	Post_Metrics_Retry_Delay_Seconds   int // delay for retrying a request that caused errors
	Post_Metrics_Retry_Max             int // max numbers of retries for a request that causes errors
	Post_Metrics_Buffer_Size           int // max numbers of requests stored in buffer queue.
}

func LoadConfig(conffile string) (*Config, error) {
	config, err := LoadConfigFile(conffile)

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
	if config.Connection.Post_Metrics_Dequeue_Delay_Seconds == 0 {
		config.Connection.Post_Metrics_Dequeue_Delay_Seconds = DefaultConfig.Connection.Post_Metrics_Dequeue_Delay_Seconds
	}
	if config.Connection.Post_Metrics_Dequeue_Delay_Seconds > POST_METRICS_DEQUEUE_DELAY_SECONDS_MAX {
		configLogger.Warningf("'post_metrics_dequese_delay_seconds' is set to %d (Maximum Value).", POST_METRICS_DEQUEUE_DELAY_SECONDS_MAX)
		config.Connection.Post_Metrics_Dequeue_Delay_Seconds = POST_METRICS_DEQUEUE_DELAY_SECONDS_MAX
	}
	if config.Connection.Post_Metrics_Retry_Delay_Seconds == 0 {
		config.Connection.Post_Metrics_Retry_Delay_Seconds = DefaultConfig.Connection.Post_Metrics_Retry_Delay_Seconds
	}
	if config.Connection.Post_Metrics_Retry_Delay_Seconds > POST_METRICS_RETRY_DELAY_SECONDS_MAX {
		configLogger.Warningf("'post_metrics_retry_delay_seconds' is set to %d (Maximum Value).", POST_METRICS_RETRY_DELAY_SECONDS_MAX)
		config.Connection.Post_Metrics_Retry_Delay_Seconds = POST_METRICS_RETRY_DELAY_SECONDS_MAX
	}
	if config.Connection.Post_Metrics_Retry_Max == 0 {
		config.Connection.Post_Metrics_Retry_Max = DefaultConfig.Connection.Post_Metrics_Retry_Max
	}
	if config.Connection.Post_Metrics_Buffer_Size == 0 {
		config.Connection.Post_Metrics_Buffer_Size = DefaultConfig.Connection.Post_Metrics_Buffer_Size
	}

	return config, err
}

func LoadConfigFile(file string) (*Config, error) {
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
