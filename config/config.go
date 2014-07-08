package config

import (
	"github.com/BurntSushi/toml"
	"github.com/mackerelio/mackerel-agent/logging"
	"time"
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
	Metrics         MetricsConfig
	Plugin          map[string]PluginConfigs
	DeprecatedSensu map[string]PluginConfigs `toml:"sensu"` // DEPRECATED this is for backward compatibility
}

type PluginConfigs map[string]PluginConfig

type PluginConfig struct {
	Command string
}

type MetricsConfig struct {
	Dequeue_Delay time.Duration // delay for dequeuing from buffer queue
	Retry_Delay   time.Duration // delay for retring a request that causes errors
	Retry_Max     int           // max numbers of retries for a request that causes errors
	Buffer_Size   int           // max numbers of requests stored in buffer queue.
}

func LoadConfig(conffile string) (Config, error) {
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
	if config.Metrics.Dequeue_Delay == 0 {
		config.Metrics.Dequeue_Delay = DefaultConfig.Metrics.Dequeue_Delay
	}
	if config.Metrics.Retry_Delay == 0 {
		config.Metrics.Retry_Delay = DefaultConfig.Metrics.Retry_Delay
	}
	if config.Metrics.Retry_Max == 0 {
		config.Metrics.Retry_Max = DefaultConfig.Metrics.Retry_Max
	}
	if config.Metrics.Buffer_Size == 0 {
		config.Metrics.Buffer_Size = DefaultConfig.Metrics.Buffer_Size
	}

	return config, err
}

func LoadConfigFile(file string) (Config, error) {
	var config Config
	if _, err := toml.DecodeFile(file, &config); err != nil {
		return config, err
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
