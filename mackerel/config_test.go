package mackerel

import (
	"io/ioutil"
	"os"
	"testing"
)

var sampleConfig = `
apikey = "abcde"

[plugin.metrics.mysql]
command = "ruby /path/to/your/plugin/mysql.rb"

[sensu.checks.memory] # for backward compatibility
command = "ruby ../sensu/plugins/system/memory-metrics.rb"
type = "metric"
`

func TestLoadConfig(t *testing.T) {
	tmpFile, error := ioutil.TempFile("/tmp", "")
	if error != nil {
		t.Error("should not raise error")
	}
	if err := ioutil.WriteFile(tmpFile.Name(), []byte(sampleConfig), 0644); err != nil {
		t.Error("should not raise error")
	}

	config, err := LoadConfig(tmpFile.Name())
	if err != nil {
		t.Error("should not raise error")
	}

	if config.Apibase != "https://mackerel.io" {
		t.Error("should be https://mackerel.io (arg value should be used)")
	}

	if config.Apikey != "abcde" {
		t.Error("should be abcde (config value should be used)")
	}
}

func TestLoadConfigFile(t *testing.T) {
	tmpFile, error := ioutil.TempFile("", "mackerel-config-test")
	if error != nil {
		t.Error("should not raise error")
	}
	if _, err := tmpFile.WriteString(sampleConfig); err != nil {
		t.Fatal("should not raise error")
	}
	tmpFile.Sync()
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	config, err := LoadConfigFile(tmpFile.Name())
	if err != nil {
		t.Error("should not raise error")
	}

	if config.Apikey != "abcde" {
		t.Error("Apikey should be abcde")
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
}
