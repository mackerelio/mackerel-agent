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
		"plugin.foo.bar",
		"plugin.metric.first",
		"plugin.check.first",
		"plugin.check.second",
		"plugins",
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
