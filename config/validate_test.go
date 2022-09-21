package config

import (
	"os"
	"testing"
)

var configFile = `
apikey = "hoge"
podfile = "/path/to/pidfile"

[foobar]
command = "test command"
env = { FOO = "BAR" }

[plugin.foo.bar]
command = "test command"
env = { FOO = "BAR" }

[plugin.metric.first]
command = "test command"

[plugin.check.first]
command = "test command"

[plugin.check.second]
command = "test command"

[plugins.check.first]
command = "test command"

[plugin.metrics.correct]
command = "test command"

[plugin.checks.correct]
command = "test command"
`

func TestValidateConfig(t *testing.T) {
	tmpFile, err := newTempFileWithContent(configFile)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

	validResult, err := ValidateConfigFile(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	wantUnexpectedKey := []string{
		"podfile",
		"foobar",
		"foobar.command",
		"foobar.env",
		"foobar.env.FOO",
		"plugin.foo.bar",
		// don't detect child of plugin.<unexpected>.<unexpected>
		// "plugin.foo.bar.command",
		// "plugin.foo.bar.env",
		// "plugin.foo.bar.env.FOO",
		"plugin.metric.first",
		"plugin.check.first",
		"plugin.check.second",
		"plugins.check.first",
		"plugins.check.first.command",
	}

	if len(wantUnexpectedKey) != len(validResult) {
		t.Errorf("should be more undecoded keys: want %v, validResult: %v", len(wantUnexpectedKey), len(validResult))
	}

	for _, v := range validResult {
		if !contains(wantUnexpectedKey, v) {
			t.Errorf("should be Undecoded: %v", v)
		}
	}
}

func contains(target []string, want string) bool {
	for _, v := range target {
		if v == want {
			return true
		}
	}
	return false
}
