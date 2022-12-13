package config

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/BurntSushi/toml"
	"github.com/agext/levenshtein"
)

// UnexpectedKey represents result of validation
type UnexpectedKey struct {
	Key        string
	SuggestKey string
}

func normalizeKey(f reflect.StructField) string {
	key := string(unicode.ToLower([]rune(f.Name)[0])) + f.Name[1:]
	if s := f.Tag.Get("toml"); s != "" {
		key = s
	}
	return key
}

func addKeyToCandidates(f reflect.StructField, candidates []string) []string {
	key := normalizeKey(f)
	return append(candidates, key)
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
			candidates = addKeyToCandidates(f, candidates)
			childCandidates := makeCandidates(f.Type)
			candidates = append(candidates, childCandidates...)
			continue
		}
		candidates = addKeyToCandidates(f, candidates)
	}

	return candidates
}

func keySuggestion(given string, candidates []string) string {
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

	var parentConfKeys []string
	configFields := reflect.VisibleFields(reflect.TypeOf(c))
	for _, f := range configFields {
		if s := f.Tag.Get("conf"); s == "parent" {
			parentConfKeys = append(parentConfKeys, normalizeKey(f))
		}
	}

	var unexpectedKeys []UnexpectedKey
	var detectedKeys []string

	undecodedKeys := md.Undecoded()
	sort.Slice(undecodedKeys, func(i, j int) bool {
		return undecodedKeys[i].String() < undecodedKeys[j].String()
	})

	/**
		```
		[plugin.checks.incorrect]
		command = "test command"
		action = { command = "test command", user = "test user", en = { TEST_KEY = "VALUE_1" }

		[plugins.check.incorrect]
		command = "test command"
		```

		undecodedKeys -> [
			plugin.checks.incorrect.action.en,
			plugin.checks.incorrect.action.en.TEST_KEY,
			plugins.check.incorrect,
			plugins.check.incorrect.command,
		]
	**/

	for _, v := range undecodedKeys {
		/*
			v: plugin.checks.incorrect.action.en
			splitedKey -> [plugin, checks, incorrect, action, en]
			topKey -> plugin
			lastKey -> en
			parentKey -> plugin.checks.incorrect.action
		*/
		splitedKey := strings.Split(v.String(), ".")
		topKey := splitedKey[0]
		lastKey := splitedKey[len(splitedKey)-1]
		parentKey := strings.Join(splitedKey[:len(splitedKey)-1], ".")
		// When parentKey (e.g., plugin.checks.incorrect.action.en) or topKey (e.g., plugins) already exists in detected keys,
		// childKey (e.g., plugin.checks.incorrect.action.en.TEST_KEY, plugins.check.incorrect.command) isn't detected.
		if containKey(detectedKeys, parentKey) || containKey(detectedKeys, topKey) {
			continue
		}

		var key string
		var suggestKey string

		if containKey(parentConfKeys, topKey) {
			key = v.String() // same as parantKey + "." + lastKey
			suggestResult := keySuggestion(lastKey, candidates)
			if suggestResult == "" {
				suggestKey = ""
			} else {
				suggestKey = parentKey + "." + suggestResult
			}
		} else {
			key = topKey
			suggestKey = keySuggestion(topKey, candidates)
		}

		unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
			key,
			suggestKey,
		})

		detectedKeys = append(detectedKeys, key)
	}

	for k1, v := range config.Plugin {
		/*
			detect [plugin.<unexpected>.<???>]
			<unexpected> should be "metrics" or "checks" or "metadata"
			default suggestion of [plugin.<unexpected>.<???>] is plugin.metrics.<???>

			don't have to detect critical syntax error about plugin here because error should have occured while loading config in config.go

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
			suggestResult := keySuggestion(k1, []string{"metrics", "checks", "metadata"})
			for k2 := range v {
				var key string = fmt.Sprintf("plugin.%s.%s", k1, k2)
				var suggestKey string

				if suggestResult == "" {
					suggestKey = fmt.Sprintf("plugin.metrics.%s", k2)
				} else {
					suggestKey = fmt.Sprintf("plugin.%s.%s", suggestResult, k2)
				}

				unexpectedKeys = append(unexpectedKeys, UnexpectedKey{
					key,
					suggestKey,
				})
			}
		}
	}

	sort.Slice(unexpectedKeys, func(i, j int) bool {
		return unexpectedKeys[i].Key < unexpectedKeys[j].Key
	})

	return unexpectedKeys, nil
}

func containKey(target []string, want string) bool {
	for _, v := range target {
		if v == want {
			return true
		}
	}
	return false
}
