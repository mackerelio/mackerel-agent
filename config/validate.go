package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/agext/levenshtein"
	"sort"
	"strings"
)

// UnexpectedKey represents result of validation
type UnexpectedKey struct {
	Name        string
	SuggestName string
}

func nameSuggestion(given string, suggestions []string) string {
	for _, suggestion := range suggestions {
		dist := levenshtein.Distance(given, suggestion, nil)
		if dist < 3 {
			return suggestion
		}
	}
	return ""
}

// ValidateConfigFile detect unexpected key in configfile
func ValidateConfigFile(file string) ([]UnexpectedKey, error) {
	config := &Config{}
	md, err := toml.DecodeFile(file, config)
	if err != nil {
		return nil, fmt.Errorf("failed to test config: %s", err)
	}

	var suggestions = []string{
		// from type Config
		"apikey",
		"pidfile",
		"root",
		"pidfile",
		"conffile",
		"roles",
		"verbose",
		"silent",
		"diagnostic",
		"display_name",
		"host_status",
		"on_start",
		"on_stop",
		"filesystems",
		"ignore",
		"use_mountpoint",
		"interfaces",
		"http_proxy",
		"https_proxy",
		"cloud_platform",
		"plugin",
		"include",
		// from type PluginConfig
		"notification_interval",
		"check_interval",
		"execution_interval",
		"max_check_attempts",
		"custom_identifier",
		"prevent_alert_auto_close",
		"include_pattern",
		"exclude_pattern",
		"action",
		"memo",
		// from type CommandConfig
		"command",
		"user",
		"env",
		"timeout_seconds",
	}

	var unexpectedKeys []UnexpectedKey
	for _, v := range md.Undecoded() {
		splitedKey := strings.Split(v.String(), ".")
		key := splitedKey[0]
		if key == "host_status" || key == "filesystems" || key == "interfaces" || key == "plugin" {
			/*
					if conffile is following, UnexpectedKey.SuggestName should be `filesystems.use_mountpoint`, not `filesystems`
					```
					[filesystems]
				  use_mntpoint = true
					```
			*/
			suggestName := nameSuggestion(splitedKey[len(splitedKey)-1], suggestions)
			if suggestName == "" {
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					v.String(),
					"",
				})
			} else {
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					v.String(),
					strings.Join(splitedKey[:len(splitedKey)-1], ".") + "." + suggestName,
				})
			}
		} else {
			// don't accept duplicate unexpectedKey
			if !containKey(unexpectedKeys, key) {
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					key,
					nameSuggestion(key, suggestions),
				})
			}
		}
	}

	for k1, v := range config.Plugin {
		/*
			detect [plugin.<unexpected>.<???>]
			default suggestion of [plugin.<unexpected>.<???>] is plugin.metrics.<???>
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
			suggestName := nameSuggestion(k1, []string{"metrics", "checks", "metadata"})
			for k2 := range v {
				if suggestName == "" {
					unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
						fmt.Sprintf("plugin.%s.%s", k1, k2),
						fmt.Sprintf("plugin.metrics.%s", k2),
					})
				} else {
					unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
						fmt.Sprintf("plugin.%s.%s", k1, k2),
						fmt.Sprintf("plugin.%s.%s", suggestName, k2),
					})
				}
			}
		}
	}

	sort.Slice(unexpectedKeys, func(i, j int) bool {
		return unexpectedKeys[i].Name < unexpectedKeys[j].Name
	})

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
