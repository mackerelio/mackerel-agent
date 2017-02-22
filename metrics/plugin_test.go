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
	conf := &config.MetricPlugin{
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
	g := &pluginGenerator{Config: &config.MetricPlugin{
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
