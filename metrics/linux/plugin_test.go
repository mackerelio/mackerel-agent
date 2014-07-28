// +build linux

package linux

import (
	"regexp"
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
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
	g := &PluginGenerator{conf}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if !containsKeyRegexp(values, "dice") {
		t.Errorf("Value for dice should be collected")
	}
}

func TestPluginCollectValues(t *testing.T) {
	g := &PluginGenerator{
		config.PluginConfig{
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
	g := &PluginGenerator{
		config.PluginConfig{
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

func TestPluginObtainConfiguration(t *testing.T) {
	g := &PluginGenerator{
		config.PluginConfig{
			Command: `echo '# mackerel-agent-plugin version=1
[[schema.graphs]]
name = "my.mysql.query"
label = "MySQL query"
[[schema.graphs.metrics]]
name = "query.foo1"
label = "Foo-1"
[[schema.graphs.metrics]]
name = "query.foo2"
label = "Foo-2"

[[schema.graphs]]
name = "my.mysql.memory"
label = "MySQL memory"
[[schema.graphs.metrics]]
name = "memory.bar1"
label = "Foo-1"
[[schema.graphs.metrics]]
name = "memory.bar2"
label = "Foo-2"
'
`,
		},
	}

	meta, err := g.loadPluginMeta()
	if meta == nil {
		t.Errorf("should parse meta: %s", err)
	}

	if len(meta.Schema.Graphs) != 2 ||
		meta.Schema.Graphs[0].Name != "my.mysql.query" ||
		len(meta.Schema.Graphs[0].Metrics) != 2 ||
		meta.Schema.Graphs[0].Metrics[0].Name != "query.foo1" {

		t.Errorf("loading meta failed got: %+v", meta)
	}

	generatorWithoutConf := &PluginGenerator{
		config.PluginConfig{
			Command: "echo \"just.echo.1\t1\t1397822016\"",
		},
	}

	_, err = generatorWithoutConf.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}

	generatorWithBadVersion := &PluginGenerator{
		config.PluginConfig{
			Command: `echo "# mackerel-agent-plugin version=666"`,
		},
	}

	_, err = generatorWithBadVersion.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}
}
