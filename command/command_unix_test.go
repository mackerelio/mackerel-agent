// +build linux darwin freebsd netbsd

package command

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
)

var diceCommand = "../example/metrics-plugins/dice-with-meta.rb"

func TestRunOnce(t *testing.T) {
	if testing.Short() {
		origMetricsInterval := metricsInterval
		metricsInterval = 1 * time.Second
		defer func() {
			metricsInterval = origMetricsInterval
		}()
	}

	conf := &config.Config{
		Plugin: map[string]config.PluginConfigs{
			"metrics": map[string]config.PluginConfig{
				"metric1": config.PluginConfig{
					Command: diceCommand,
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
	if os.Getenv("TRAVIS") != "" {
		t.Skip("Skip in travis")
	}

	if testing.Short() {
		origMetricsInterval := metricsInterval
		metricsInterval = 1
		defer func() {
			metricsInterval = origMetricsInterval
		}()
	}

	conf := &config.Config{
		Plugin: map[string]config.PluginConfigs{
			"metrics": map[string]config.PluginConfig{
				"metric1": config.PluginConfig{
					Command: diceCommand,
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
