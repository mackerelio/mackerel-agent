// +build linux

package linux

import (
	"regexp"
	"testing"

	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
)

func containsKeyRegexp(values metrics.Values, reg string) bool {
	for name, _ := range values {
		if matches := regexp.MustCompile(reg).FindStringSubmatch(name); matches != nil {
			return true
		}
	}
	return false
}

func TestPluginGenerate(t *testing.T) {
	config := mackerel.PluginConfig{
		Command: "ruby ../../example/metrics-plugins/dice.rb",
	}
	g := &PluginGenerator{config}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if !containsKeyRegexp(values, "dice") {
		t.Errorf("Value for dice should be collected")
	}
}

func TestPluginCollectValues(t *testing.T) {
	g := &PluginGenerator{}
	command := "ruby ../../example/metrics-plugins/dice.rb"
	values, err := g.collectValues(command)
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if !containsKeyRegexp(values, "dice") {
		t.Errorf("Value for dice should be collected")
	}
}

func TestPluginCollectValuesCommand(t *testing.T) {
	g := &PluginGenerator{}
	command := "echo \"just.echo.1\t1\t1397822016\""

	values, err := g.collectValues(command)
	if err != nil {
		t.Error("should not raise error")
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
