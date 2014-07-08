package config

import (
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
}

type PluginConfigs map[string]PluginConfig

type PluginConfig struct {
	Command string
}

type ConnectionConfig struct {
	Post_Metrics_Dequeue_Delay_Seconds int // delay for dequeuing from buffer queue
	Post_Metrics_Retry_Delay_Seconds   int // delay for retring a request that causes errors
	Post_Metrics_Retry_Max             int // max numbers of retries for a request that causes errors
	Post_Metrics_Buffer_Size           int // max numbers of requests stored in buffer queue.
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
	if config.Connection.Post_Metrics_Dequeue_Delay_Seconds == 0 {
		config.Connection.Post_Metrics_Dequeue_Delay_Seconds = DefaultConfig.Connection.Post_Metrics_Dequeue_Delay_Seconds
	}
	if config.Connection.Post_Metrics_Retry_Delay_Seconds == 0 {
		config.Connection.Post_Metrics_Retry_Delay_Seconds = DefaultConfig.Connection.Post_Metrics_Retry_Delay_Seconds
	}
	if config.Connection.Post_Metrics_Retry_Max == 0 {
		config.Connection.Post_Metrics_Retry_Max = DefaultConfig.Connection.Post_Metrics_Retry_Max
	}
	if config.Connection.Post_Metrics_Buffer_Size == 0 {
		config.Connection.Post_Metrics_Buffer_Size = DefaultConfig.Connection.Post_Metrics_Buffer_Size
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
