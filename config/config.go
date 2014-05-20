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
	Roles           []string
	Verbose         bool
	Plugin          map[string]PluginConfigs
	DeprecatedSensu map[string]PluginConfigs `toml:"sensu"` // DEPRECATED this is for backward compatibility
}

type PluginConfigs map[string]PluginConfig

type PluginConfig struct {
	Command string
}

var DefaultConfig = &Config{
	Apibase: "https://mackerel.io",
	Root:    "/var/lib/mackerel-agent",
	Pidfile: "/var/run/mackerel-agent.pid",
	Roles:   []string{},
	Verbose: false,
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
