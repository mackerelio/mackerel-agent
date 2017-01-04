// +build windows

package metrics

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestPluginCollectValuesCommand(t *testing.T) {
	g := &pluginGenerator{Config: &config.PluginConfig{
		Command: "echo just.echo.1	1	1397822016",
	},
	}

	values, err := g.collectValues()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if len(values) != 1 {
		t.Error("Only 1 value shoud be generated")
	}

	for name, value := range values {
		if name != "custom.just.echo.1" {
			t.Errorf("Wrong name: %s", name)
		}
		if value != 1.0 {
			t.Errorf("Wrong value: %+v", value)
		}
	}
}

func TestPluginCollectValuesCommandWithSpaces(t *testing.T) {
	g := &pluginGenerator{Config: &config.PluginConfig{
		Command: `echo just.echo.2   2   1397822016`,
	}}

	values, err := g.collectValues()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if len(values) != 1 {
		t.Error("Only 1 value shoud be generated")
	}

	for name, value := range values {
		if name != "custom.just.echo.2" {
			t.Errorf("Wrong name: %s", name)
		}
		if value != 2.0 {
			t.Errorf("Wrong value: %+v", value)
		}
	}
}

func TestPluginLoadPluginMeta(t *testing.T) {
	// TODO test parse

	generatorWithoutConf := &pluginGenerator{
		Config: &config.PluginConfig{
			Command: "echo just.echo.1	1	1397822016",
		},
	}

	err = generatorWithoutConf.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}

	generatorWithBadVersion := &pluginGenerator{
		Config: &config.PluginConfig{
			Command: `echo # mackerel-agent-plugin version=666`,
		},
	}

	err = generatorWithBadVersion.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}
}
