package config

import (
	"fmt"
	"os"
	"testing"
)

var configFile = `
apikey = "123456abcde"
podfile = "/path/to/pidfile"

[foobar]
command = "test command"
env = { FOO = "BAR" }

[filesystems]
ignore = "/path/to/ignore"
use_mntpoint = true

[plugin.foo.bar]
command = "test command"
env = { FOO = "BAR" }

[plugin.metric.1]
command = "test command"

[plugin.check.1]
command = "test command"

[plugin.check.2]
command = "test command"

[plugins.check.1]
command = "test command"

[plugin.metrics.correct]
command = "test command"

[plugin.checks.correct]
command = "test command"

[plugin.metrics.1]
command = "test command"
action = { command = "test command", use = "test user", env = { TEST_KEY = "VALUE_1" } }

[plugin.metrics.2]
command = "test command"
action = { command = "test command", xxx = "yyy" }

[plugin.checks.1]
command = "test command"
action = { command = "test command", user = "test user", en = { TEST_KEY = "VALUE_1" } }

[plugin.metrics.aaa]
command = ["exit 1"]
[plugin.metrics.aaa.action]
command = ["exit 2"]
max_check_attempts = 3
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
		{"plugin.check.1", "plugin.checks.1"},
		{"plugin.check.2", "plugin.checks.2"},
		{"plugin.checks.1.action.en", "plugin.checks.1.action.env"},
		{"plugin.checks.1.action.en.TEST_KEY", ""},
		{"plugin.foo.bar", "plugin.metrics.bar"},
		{"plugin.metric.1", "plugin.metrics.1"},
		{"plugin.metrics.1.action.use", "plugin.metrics.1.action.user"},
		{"plugin.metrics.2.action.xxx", ""},
		{"plugin.metrics.aaa.action.max_check_attempts", ""},
		{"plugins", "plugin"},
		{"podfile", "pidfile"},
	}

	if len(wantUnexpectedKey) != len(validResult) {
		t.Errorf("should be more undecoded keys: want %v, validResult: %v", len(wantUnexpectedKey), len(validResult))
	}

	for i, v := range validResult {
		fmt.Println(i, v)
		if wantUnexpectedKey[i].Name != v.Name {
			t.Errorf("expect Name: %v, actual Name: %v", wantUnexpectedKey[i].Name, v.Name)
		}

		if wantUnexpectedKey[i].SuggestName != v.SuggestName {
			t.Errorf("expect SuggestName: %v, actual SuggestName: %v", wantUnexpectedKey[i].SuggestName, v.SuggestName)
		}
	}
}
