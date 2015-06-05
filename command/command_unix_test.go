// +build linux darwin freebsd

package command

import (
	"reflect"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
)

func init() {
	metricsInterval = 1
}

func TestRunOnce(t *testing.T) {
	conf := &config.Config{
		Plugin: map[string]config.PluginConfigs{
			"metrics": map[string]config.PluginConfig{
				"metric1": config.PluginConfig{
					Command: "ruby ../example/metrics-plugins/dice-with-meta.rb",
				},
			},
			"checks": map[string]config.PluginConfig{
				"check1": config.PluginConfig{
					Command: "echo 1",
				},
			},
		},
	}
	err := RunOnce(conf)
	if err != nil {
		t.Errorf("RunOnce() should be nomal exit: %s", err)
	}
}

func TestRunOncePayload(t *testing.T) {
	conf := &config.Config{
		Plugin: map[string]config.PluginConfigs{
			"metrics": map[string]config.PluginConfig{
				"metric1": config.PluginConfig{
					Command: "ruby ../example/metrics-plugins/dice-with-meta.rb",
				},
			},
			"checks": map[string]config.PluginConfig{
				"check1": config.PluginConfig{
					Command: "echo 1",
				},
			},
		},
	}
	graphdefs, hostSpec, metrics, err := runOncePayload(conf)
	if err != nil {
		t.Errorf("RunOnce() should be nomal exit: %s", err)
	}

	if !reflect.DeepEqual(graphdefs[0], mackerel.CreateGraphDefsPayload{
		Name:        "custom.dice",
		DisplayName: "My Dice",
		Unit:        "integer",
		Metrics: []mackerel.CreateGraphDefsPayloadMetric{
			mackerel.CreateGraphDefsPayloadMetric{
				Name:        "custom.dice.d6",
				DisplayName: "Die (d6)",
				IsStacked:   false,
			},
			mackerel.CreateGraphDefsPayloadMetric{
				Name:        "custom.dice.d20",
				DisplayName: "Die (d20)",
				IsStacked:   false,
			},
		},
	}) {
		t.Errorf("graphdefs are invalid")
	}

	if hostSpec.Name == "" {
		t.Errorf("hostname should be set")
	}
	if hostSpec.Checks[0] != "check1" {
		t.Errorf("first check name should be check1")
	}

	if metrics.Values["custom.dice.d20"] == 0 {
		t.Errorf("custom.dice.d20 name should be set")
	}

}
