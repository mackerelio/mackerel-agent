// +build linux darwin freebsd netbsd

package command

import (
	"reflect"
	"testing"
	"time"

	"github.com/mackerelio/mackerel-agent/config"
	mkr "github.com/mackerelio/mackerel-client-go"
)

var diceCommand = "go run ../_example/metrics-plugins/dice-with-meta.go"

func TestRunOnce(t *testing.T) {
	if testing.Short() {
		origMetricsInterval := metricsInterval
		metricsInterval = 1 * time.Second
		defer func() {
			metricsInterval = origMetricsInterval
		}()
	}

	conf := &config.Config{
		MetricPlugins: map[string]*config.MetricPlugin{
			"metric1": {
				Command: config.Command{Cmd: diceCommand},
			},
		},
		CheckPlugins: map[string]*config.CheckPlugin{
			"check1": {
				Command: config.Command{Cmd: "echo 1"},
			},
		},
	}
	err := RunOnce(conf, &AgentMeta{})
	if err != nil {
		t.Errorf("RunOnce() should be nomal exit: %s", err)
	}
}

func TestRunOncePayload(t *testing.T) {
	if testing.Short() {
		origMetricsInterval := metricsInterval
		metricsInterval = 1 * time.Second
		defer func() {
			metricsInterval = origMetricsInterval
		}()
	}

	conf := &config.Config{
		MetricPlugins: map[string]*config.MetricPlugin{
			"metric1": {
				Command: config.Command{Cmd: diceCommand},
			},
		},
		CheckPlugins: map[string]*config.CheckPlugin{
			"check1": {
				Command: config.Command{Cmd: "echo 1"},
			},
		},
	}
	graphdefs, hostSpec, metrics, err := runOncePayload(conf, &AgentMeta{})
	if err != nil {
		t.Errorf("RunOnce() should be nomal exit: %s", err)
	}

	if !reflect.DeepEqual(graphdefs[0], &mkr.GraphDefsParam{
		Name:        "custom.dice",
		DisplayName: "My Dice",
		Unit:        "integer",
		Metrics: []*mkr.GraphDefsMetric{
			{
				Name:        "custom.dice.d6",
				DisplayName: "Die (d6)",
				IsStacked:   false,
			},
			{
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
	if hostSpec.Checks[0].Name != "check1" {
		t.Errorf("first check name should be check1")
	}

	if len(metrics.Values) != 1 {
		t.Errorf("there must be some metric values")
	}
	if metrics.Values[0].Values["custom.dice.d20"] == 0 {
		t.Errorf("custom.dice.d20 name should be set")
	}

}
