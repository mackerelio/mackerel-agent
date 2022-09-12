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

# TODO: detect warning
[plugin.check.hoge]
command = "test command"

[plugins.check.hoge]
command = "test command"

# this is correct
[plugin.checks.hoge]
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

	want := []string{
		"podfile",
		"foobar",
		"foobar.command",
		"plugins.check.hoge",
		"plugins.check.hoge.command",
	}

	if len(want) != len(validResult) {
		t.Errorf("should be more undecoded keys: want %v, validResult: %v", len(want), len(validResult))
	}

	for _, v := range validResult {
		if !contains(want, v.String()) {
			t.Errorf("should be Undecoded: %v", v.String())
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
