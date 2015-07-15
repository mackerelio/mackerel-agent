// +build linux darwin freebsd

package metrics

import (
	"regexp"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
)

func containsKeyRegexp(values Values, reg string) bool {
	for name := range values {
		if matches := regexp.MustCompile(reg).FindStringSubmatch(name); matches != nil {
			return true
		}
	}
	return false
}

func TestPluginGenerate(t *testing.T) {
	conf := config.PluginConfig{
		Command: "ruby ../example/metrics-plugins/dice.rb",
	}
	g := &pluginGenerator{Config: conf}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if !containsKeyRegexp(values, "dice") {
		t.Errorf("Value for dice should be collected")
	}
}

func TestPluginCollectValues(t *testing.T) {
	g := &pluginGenerator{Config: config.PluginConfig{
		Command: "ruby ../example/metrics-plugins/dice.rb",
	},
	}
	values, err := g.collectValues()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if !containsKeyRegexp(values, "dice") {
		t.Errorf("Value for dice should be collected")
	}
}

func TestPluginCollectValuesCommand(t *testing.T) {
	g := &pluginGenerator{Config: config.PluginConfig{
		Command: "echo \"just.echo.1\t1\t1397822016\"",
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
	g := &pluginGenerator{Config: config.PluginConfig{
		Command: `echo "just.echo.2   2   1397822016"`,
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
	g := &pluginGenerator{
		Config: config.PluginConfig{
			Command: `echo '# mackerel-agent-plugin version=1
{
  "graphs": {
    "query": {
      "label": "MySQL query",
      "unit": "integer",
      "metrics": [
        {
          "name": "foo1",
          "label": "Foo-1"
        },
        {
          "name": "foo2",
          "label": "Foo-2",
          "stacked": true
        }
      ]
    },
    "memory": {
      "label": "MySQL memory",
      "unit": "float",
      "metrics": [
        {
          "name": "bar1",
          "label": "Bar-1"
        },
        {
          "name": "bar2",
          "label": "Bar-2"
        }
      ]
    }
  }
}
'`,
		},
	}

	err := g.loadPluginMeta()
	if g.Meta == nil {
		t.Errorf("should parse meta: %s", err)
	}

	if g.Meta.Graphs["query"].Label != "MySQL query" ||
		g.Meta.Graphs["query"].Metrics[0].Label != "Foo-1" ||
		g.Meta.Graphs["query"].Unit != "integer" ||
		g.Meta.Graphs["query"].Metrics[1].Label != "Foo-2" ||
		g.Meta.Graphs["query"].Metrics[1].Stacked != true ||
		g.Meta.Graphs["memory"].Metrics[0].Label != "Bar-1" ||
		g.Meta.Graphs["memory"].Unit != "float" {

		t.Errorf("loading meta failed got: %+v", g.Meta)
	}

	generatorWithoutConf := &pluginGenerator{
		Config: config.PluginConfig{
			Command: "echo \"just.echo.1\t1\t1397822016\"",
		},
	}

	err = generatorWithoutConf.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}

	generatorWithBadVersion := &pluginGenerator{
		Config: config.PluginConfig{
			Command: `echo "# mackerel-agent-plugin version=666"`,
		},
	}

	err = generatorWithBadVersion.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}
}

func TestPluginMakeCreateGraphDefsPayload(t *testing.T) {
	// this plugin emits "one.foo1", "one.foo2" and "two.bar1" metrics
	g := &pluginGenerator{
		Meta: &pluginMeta{
			Graphs: map[string]customGraphDef{
				"one": {
					Label: "My Graph One",
					Unit:  "integer",
					Metrics: []customGraphMetricDef{
						{
							Name:    "foo1",
							Label:   "Foo(1)",
							Stacked: true,
						},
						{
							Name:    "foo2",
							Label:   "Foo(2)",
							Stacked: true,
						},
					},
				},
				"two": {
					Label: "My Graph Two",
					Metrics: []customGraphMetricDef{
						{
							Name:  "bar1",
							Label: "Bar(1)",
						},
					},
				},
			},
		},
	}

	payloads := g.makeCreateGraphDefsPayload()

	if len(payloads) != 2 {
		t.Errorf("Bad payload created: %+v", payloads)
	}

	var payloadOne *mackerel.CreateGraphDefsPayload
	for _, payload := range payloads {
		if payload.Name == "custom.one" {
			payloadOne = &payload
			break
		}
	}

	if payloadOne == nil {
		t.Errorf("Payload with name custom.one not found: %+v", payloads)
	}

	if payloadOne.DisplayName != "My Graph One" ||
		len(payloadOne.Metrics) != 2 ||
		payloadOne.Unit != "integer" {
		t.Errorf("Bad payload created: %+v", payloadOne)
	}

	var metricOneFoo1 *mackerel.CreateGraphDefsPayloadMetric
	for _, metric := range payloadOne.Metrics {
		if metric.Name == "custom.one.foo1" {
			metricOneFoo1 = &metric
			break
		}
	}
	if metricOneFoo1 == nil {
		t.Errorf("Metric payload with name custom.one.foo1 not fonud: %+v", payloadOne)
	}

	if metricOneFoo1.DisplayName != "Foo(1)" ||
		metricOneFoo1.IsStacked != true {

		t.Errorf("Bat metric payload created: %+v", metricOneFoo1)
	}
}
