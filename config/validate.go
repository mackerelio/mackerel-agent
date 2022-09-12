package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
)

func ValidateConfigFile(file string) ([]string, error) {
	config := &Config{}
	md, err := toml.DecodeFile(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to test config: %s", err)
	}

	var unexpectedKey []string = make([]string, 0)
	for _, v := range md.Undecoded() {
		unexpectedKey = append(unexpectedKey, v.String())
	}

	for k1, v := range config.Plugin {
		/*
			detect [plugin.<unexpected>.<unexpected>]
			don't have to detect critical syntax error about plugin here because error should have occured while loading config
			```
			[plugin.metrics.correct]
			```
			-> A configuration value of `command` should be string or string slice, but <nil>
			```
			[plugin.metrics]
			command = "test command"
			```
			-> type mismatch for config.PluginConfig: expected table but found string
		*/
		if k1 != "metrics" && k1 != "checks" && k1 != "metadata" {
			for k2 := range v {
				unexpectedKey = append(unexpectedKey, fmt.Sprintf("plugin.%s.%s", k1, k2))
			}
		}
	}

	return unexpectedKey, nil
}
