// +build linux

package linux

import (
	"regexp"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
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
	conf := config.PluginConfig{
		Command: "ruby ../../example/metrics-plugins/dice.rb",
	}
	g := &PluginGenerator{Config: conf}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if !containsKeyRegexp(values, "dice") {
		t.Errorf("Value for dice should be collected")
	}
}

func TestPluginCollectValues(t *testing.T) {
	g := &PluginGenerator{Config: config.PluginConfig{
		Command: "ruby ../../example/metrics-plugins/dice.rb",
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
	g := &PluginGenerator{Config: config.PluginConfig{
		Command: "echo \"just.echo.1\t1\t1397822016\"",
	},
	}

	values, err := g.collectValues()
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

func TestPluginLoadPluginMeta(t *testing.T) {
	g := &PluginGenerator{
		Config: config.PluginConfig{
			Command: `echo '# mackerel-agent-plugin version=1
[graphs.query]
label = "MySQL query"
unit = "integer"
[graphs.query.metrics.foo1]
label = "Foo-1"
[graphs.query.metrics.foo2]
label = "Foo-2"
stacked = true

[graphs.memory]
label = "MySQL memory"
unit = "float"
[graphs.memory.metrics.bar1]
label = "Bar-1"
[graphs.memory.metrics.bar2]
label = "Bar-2"
'
`,
		},
	}

	err := g.loadPluginMeta()
	if g.Meta == nil {
		t.Errorf("should parse meta: %s", err)
	}

	if g.Meta.Graphs["query"].Label != "MySQL query" ||
		g.Meta.Graphs["query"].Metrics["foo1"].Label != "Foo-1" ||
		g.Meta.Graphs["query"].Unit != "integer" ||
		g.Meta.Graphs["query"].Metrics["foo2"].Label != "Foo-2" ||
		g.Meta.Graphs["query"].Metrics["foo2"].Stacked != true ||
		g.Meta.Graphs["memory"].Metrics["bar1"].Label != "Bar-1" ||
		g.Meta.Graphs["memory"].Unit != "float" {

		t.Errorf("loading meta failed got: %+v", g.Meta)
	}

	if g.Meta.Graphs["query"].Metrics["foo1"].Order != 1 ||
		g.Meta.Graphs["query"].Metrics["foo2"].Order != 2 {

		t.Errorf("orders should be filled by the appearance order: %+v", g.Meta)
	}

	generatorWithoutConf := &PluginGenerator{
		Config: config.PluginConfig{
			Command: "echo \"just.echo.1\t1\t1397822016\"",
		},
	}

	err = generatorWithoutConf.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}

	generatorWithBadVersion := &PluginGenerator{
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
	g := &PluginGenerator{
		Meta: &pluginMeta{
			Graphs: map[string]*customGraphDef{
				"one": {
					Label: "My Graph One",
					Unit:  "integer",
					Metrics: map[string]*customGraphMetricDef{
						"foo1": {
							Label:   "Foo(1)",
							Stacked: true,
							Order:   1,
						},
						"foo2": {
							Label:   "Foo(2)",
							Stacked: true,
							Order:   2,
						},
					},
				},
				"two": {
					Label: "My Graph Two",
					Metrics: map[string]*customGraphMetricDef{
						"bar1": {
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

	if payloadOne.Metrics[0].Name != "custom.one.foo1" ||
		payloadOne.Metrics[1].Name != "custom.one.foo2" {

		t.Errorf("Metrics should be sorted: %+v", payloadOne)
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
