package config

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/agext/levenshtein"
	"reflect"
	"sort"
	"strings"
	"unicode"
)

// UnexpectedKey represents result of validation
type UnexpectedKey struct {
	Name        string
	SuggestName string
}

func normalizeKeyname(f reflect.StructField) string {
	name := string(unicode.ToLower([]rune(f.Name)[0])) + f.Name[1:]
	if s := f.Tag.Get("toml"); s != "" {
		name = s
	}
	return name
}

func addKeynameToCandidates(f reflect.StructField, candidates []string) []string {
	name := normalizeKeyname(f)
	return append(candidates, name)
}

func makeCandidates(t reflect.Type) []string {
	if t.Kind() == reflect.Map || t.Kind() == reflect.Ptr {
		return makeCandidates(t.Elem())
	}
	var candidates []string
	fields := reflect.VisibleFields(t)
	for _, f := range fields {
		if s := f.Tag.Get("conf"); s == "ignore" {
			continue
		}
		if s := f.Tag.Get("conf"); s == "parent" {
			candidates = addKeynameToCandidates(f, candidates)
			childCandidates := makeCandidates(f.Type)
			candidates = append(candidates, childCandidates...)
			continue
		}
		candidates = addKeynameToCandidates(f, candidates)
	}

	return candidates
}

func nameSuggestion(given string, candidates []string) string {
	for _, candidate := range candidates {
		dist := levenshtein.Distance(given, candidate, nil)
		if dist < 3 {
			return candidate
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

	var c Config
	candidates := makeCandidates(reflect.TypeOf(c))

	var parentKeys []string
	configFields := reflect.VisibleFields(reflect.TypeOf(c))
	for _, f := range configFields {
		if s := f.Tag.Get("conf"); s == "parent" {
			parentKeys = append(parentKeys, normalizeKeyname(f))
		}
	}

	var unexpectedKeys []UnexpectedKey
	var unexpectedKeyNames []string

	undecodedKeys := md.Undecoded()
	sort.Slice(undecodedKeys, func(i, j int) bool {
		return undecodedKeys[i].String() < undecodedKeys[j].String()
	})

	for _, v := range undecodedKeys {
		splitedKey := strings.Split(v.String(), ".")
		key := splitedKey[0]
		parentKey := strings.Join(splitedKey[:len(splitedKey)-1], ".")
		if containKeyName(parentKeys, key) {
			// if parent is already exists in unexpectedKeyNames, child isn't detected.
			if containKeyName(unexpectedKeyNames, parentKey) {
				continue
			}
			/*
					if conffile is following, UnexpectedKey.SuggestName should be `filesystems.use_mountpoint`, not `filesystems`
					```
					[filesystems]
				  use_mntpoint = true
					```
			*/
			suggestName := nameSuggestion(splitedKey[len(splitedKey)-1], candidates)
			if suggestName == "" {
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					v.String(),
					"",
				})
			} else {
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					v.String(),
					parentKey + "." + suggestName,
				})
			}
			unexpectedKeyNames = append(unexpectedKeyNames, v.String())
		} else {
			// don't accept duplicate unexpectedKey
			if !containKey(unexpectedKeys, key) {
				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					key,
					nameSuggestion(key, candidates),
				})
				unexpectedKeyNames = append(unexpectedKeyNames, key)
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
				unexpectedKeyNames = append(unexpectedKeyNames, fmt.Sprintf("plugin.%s.%s", k1, k2))
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

func containKeyName(target []string, want string) bool {
	for _, v := range target {
		if v == want {
			return true
		}
	}
	return false
}
