package config

import (
	"os"
	"testing"
)

var configFile = `
apikey = "123456abcde"
podfile = "/path/to/pidfile"

[foobar]
command = ["test command"]
env = { FOO = "BAR" }

[filesystems]
ignore = "/path/to/ignore"
use_mntpoint = true

[plugin.foo.bar]
command = ["test command"]
env = { FOO = "BAR" }

[plugin.metric.incorrect1]
command = ["test command"]

[plugin.check.incorrect2]
command = ["test command"]

[plugin.check.incorrect3]
command = ["test command"]

[plugins.check.incorrect4]
command = ["test command"]

[plugin.metrics.correct]
command = ["test command"]

[plugin.checks.correct]
command = ["test command"]

[plugin.metrics.incorrect5]
command = ["test command"]
action = { command = "test command", use = "test user", env = { TEST_KEY = "VALUE_1" } }

[plugin.metrics.incorrect6]
command = ["test command"]
action = { command = "test command", xxx = "yyy" }

[plugin.metrics.incorrect7.incorrect8]
command = ["test command"]

[plugin.checks.incorrect9]
command = ["test command"]
action = { command = "test command", user = "test user", en = { TEST_KEY = "VALUE_1" } }
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

	wantUnexpectedKey := []UnexpectedKey{
		{"filesystems.use_mntpoint", "filesystems.use_mountpoint"},
		{"foobar", ""},
		{"plugin.check.incorrect2", "plugin.checks.incorrect2"},
		{"plugin.check.incorrect3", "plugin.checks.incorrect3"},
		{"plugin.checks.incorrect9.action.en", "plugin.checks.incorrect9.action.env"},
		{"plugin.foo.bar", "plugin.metrics.bar"},
		{"plugin.metric.incorrect1", "plugin.metrics.incorrect1"},
		{"plugin.metrics.incorrect5.action.use", "plugin.metrics.incorrect5.action.user"},
		{"plugin.metrics.incorrect6.action.xxx", ""},
		{"plugin.metrics.incorrect7.incorrect8", ""},
		{"plugins", "plugin"},
		{"podfile", "pidfile"},
	}

	if len(wantUnexpectedKey) != len(validResult) {
		t.Errorf("should be more undecoded keys: want %v, validResult: %v", len(wantUnexpectedKey), len(validResult))
	}

	for i, v := range validResult {
		if wantUnexpectedKey[i].Key != v.Key {
			t.Errorf("expect Key: %v, actual Key: %v", wantUnexpectedKey[i].Key, v.Key)
		}

		if wantUnexpectedKey[i].SuggestKey != v.SuggestKey {
			t.Errorf("expect SuggestKey: %v, actual SuggestKey: %v", wantUnexpectedKey[i].SuggestKey, v.SuggestKey)
		}
	}
}
