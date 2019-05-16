package metrics

import (
	"regexp"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
	mkr "github.com/mackerelio/mackerel-client-go"
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
	conf := &config.MetricPlugin{
		Command: config.Command{Cmd: "go run ../_example/metrics-plugins/dice.go"},
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
	g := &pluginGenerator{Config: &config.MetricPlugin{
		Command: config.Command{Cmd: "go run ../_example/metrics-plugins/dice.go"},
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

func TestPluginCollectValuesWithIncludePattern(t *testing.T) {
	g := &pluginGenerator{Config: &config.MetricPlugin{
		Command:        config.Command{Cmd: "go run ../_example/metrics-plugins/dice-with-meta.go"},
		IncludePattern: regexp.MustCompile(`^dice\.d6`),
	},
	}
	values, err := g.collectValues()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if len(values) != 1 {
		t.Errorf("Collected metrics are unexpected ")
	}
	if _, ok := values["custom.dice.d6"]; !ok {
		t.Errorf("Value for dice.d6 should be present ")
	}
}

func TestPluginCollectValuesWithExcludePattern(t *testing.T) {
	g := &pluginGenerator{Config: &config.MetricPlugin{
		Command:        config.Command{Cmd: "go run ../_example/metrics-plugins/dice-with-meta.go"},
		ExcludePattern: regexp.MustCompile(`^dice\.d20`),
	},
	}
	values, err := g.collectValues()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if len(values) != 1 {
		t.Errorf("Collected metrics are unexpected ")
	}
	if _, ok := values["custom.dice.d6"]; !ok {
		t.Errorf("Value for dice.d6 should be present ")
	}
}

func TestPluginCollectValuesWithBothPattern(t *testing.T) {
	g := &pluginGenerator{Config: &config.MetricPlugin{
		Command:        config.Command{Cmd: "go run ../_example/metrics-plugins/dice-with-meta.go"},
		IncludePattern: regexp.MustCompile(`^dice\.d20`),
		ExcludePattern: regexp.MustCompile(`^dice\.d20`),
	},
	}
	values, err := g.collectValues()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}
	if len(values) != 0 {
		t.Errorf("No values should be present")
	}
}

func TestPluginMakeGraphDefsParam(t *testing.T) {
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

	payloads := g.makeGraphDefsParam()

	if len(payloads) != 2 {
		t.Errorf("Bad payload created: %+v", payloads)
	}

	var payloadOne *mkr.GraphDefsParam
	for _, payload := range payloads {
		if payload.Name == "custom.one" {
			payloadOne = payload
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

	var metricOneFoo1 *mkr.GraphDefsMetric
	for _, metric := range payloadOne.Metrics {
		if metric.Name == "custom.one.foo1" {
			metricOneFoo1 = metric
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
