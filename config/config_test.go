package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

var sampleConfig = `
apikey = "abcde"

[connection]
post_metrics_retry_delay_seconds = 600
post_metrics_retry_max = 5

[plugin.metrics.mysql]
command = "ruby /path/to/your/plugin/mysql.rb"

[sensu.checks.memory] # for backward compatibility
command = "ruby ../sensu/plugins/system/memory-metrics.rb"
type = "metric"

[plugin.checks.heartbeat]
command = "heartbeat.sh"
`

func TestLoadConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if err = ioutil.WriteFile(tmpFile.Name(), []byte(sampleConfig), 0644); err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if config.Apibase != "https://mackerel.io" {
		t.Error("should be https://mackerel.io (arg value should be used)")
	}

	if config.Apikey != "abcde" {
		t.Error("should be abcde (config value should be used)")
	}

	if config.Connection.PostMetricsDequeueDelaySeconds != 30 {
		t.Error("should be 30 (default value should be used)")
	}

	if config.Connection.PostMetricsRetryDelaySeconds != 180 {
		t.Error("should be 180 (max retry delay seconds is 180)")
	}

	if config.Connection.PostMetricsRetryMax != 5 {
		t.Error("should be 5 (config value should be used)")
	}
}

func TestLoadConfigFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "mackerel-config-test")
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if _, err := tmpFile.WriteString(sampleConfig); err != nil {
		t.Fatal("should not raise error")
	}
	tmpFile.Sync()
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	config, err := loadConfigFile(tmpFile.Name())
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if config.Apikey != "abcde" {
		t.Error("Apikey should be abcde")
	}

	if config.Connection.PostMetricsRetryMax != 5 {
		t.Error("PostMetricsRetryMax should be 5")
	}

	if config.Plugin["metrics"] == nil {
		t.Error("plugin should have metrics")
	}
	pluginConf := config.Plugin["metrics"]["mysql"]
	if pluginConf.Command != "ruby /path/to/your/plugin/mysql.rb" {
		t.Errorf("plugin conf command should be 'ruby /path/to/your/plugin/mysql.rb' but %v", pluginConf.Command)
	}

	// for backward compatibility
	sensu := config.Plugin["metrics"]["DEPRECATED-sensu-memory"]
	if sensu.Command != "ruby ../sensu/plugins/system/memory-metrics.rb" {
		t.Error("sensu command should be 'ruby ../sensu/plugins/system/memory-metrics.rb'")
	}

	if config.Plugin["checks"] == nil {
		t.Error("plugin should have checks")
	}
	checks := config.Plugin["checks"]["heartbeat"]
	if checks.Command != "heartbeat.sh" {
		t.Error("sensu command should be 'heartbeat.sh'")
	}
}

func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func assert(t *testing.T, ok bool, msg string) {
	if !ok {
		t.Error(msg)
	}
}

var tomlQuotedReplacer = strings.NewReplacer(
	"\t", "\\t",
	"\n", "\\n",
	"\r", "\\r",
	"\"", "\\\"",
	"\\", "\\\\",
)

func TestLoadConfigFileInclude(t *testing.T) {
	configDir, err := ioutil.TempDir("", "mackerel-config-test")
	assertNoError(t, err)

	configFile, err := ioutil.TempFile("", "mackerel-config-test")
	assertNoError(t, err)

	includedFile, err := os.Create(filepath.Join(configDir, "sub1.conf"))

	configContent := fmt.Sprintf(`
apikey = "not overwritten"
roles = [ "roles", "to be overwritten" ]

include = "%s/*.conf"

[plugin.metrics.foo1]
command = "foo1"

[plugin.metrics.bar]
command = "this wille be overwritten"
`, tomlQuotedReplacer.Replace(configDir))

	includedContent := `
roles = [ "Service:role" ]

[plugin.metrics.foo2]
command = "foo2"

[plugin.metrics.bar]
command = "bar"
`

	_, err = configFile.WriteString(configContent)
	assertNoError(t, err)

	_, err = includedFile.WriteString(includedContent)
	assertNoError(t, err)

	configFile.Close()
	includedFile.Close()
	defer os.Remove(configFile.Name())
	defer os.Remove(includedFile.Name())

	config, err := loadConfigFile(configFile.Name())
	assertNoError(t, err)

	assert(t, config.Apikey == "not overwritten", "apikey should not be overwritten")
	assert(t, len(config.Roles) == 1, "roles should be overwritten")
	assert(t, config.Roles[0] == "Service:role", "roles should be overwritten")
	assert(t, config.Plugin["metrics"]["foo1"].Command == "foo1", "plugin.metrics.foo1 should exist")
	assert(t, config.Plugin["metrics"]["foo2"].Command == "foo2", "plugin.metrics.foo2 should exist")
	assert(t, config.Plugin["metrics"]["bar"].Command == "bar", "plugin.metrics.bar should be overwritten")
}
