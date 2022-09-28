package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"strings"
)

type UnexpectedKey struct {
	Name string
}

// ValidateConfigFile detect unexpected key in configfile
func ValidateConfigFile(file string) ([]UnexpectedKey, error) {
	config := &Config{}
	md, err := toml.DecodeFile(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to test config: %s", err)
	}

	var unexpectedKeys []UnexpectedKey
	for _, v := range md.Undecoded() {
		key := strings.Split(v.String(), ".")[0]
		if !containKey(unexpectedKeys, key) {
			unexpectedKeys = append(unexpectedKeys, UnexpectedKey{Name: key})
		}
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
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{Name: fmt.Sprintf("plugin.%s.%s", k1, k2)})
			}
		}
	}

	return unexpectedKeys, nil
}

func containKey(target []UnexpectedKey, want string) bool {
	for _, v := range target {
		if v.Name == want {
			return true
		}
	}
	return false
}
