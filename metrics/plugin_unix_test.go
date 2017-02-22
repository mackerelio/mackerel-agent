// +build linux darwin freebsd netbsd

package metrics

import (
	"testing"

	"github.com/mackerelio/mackerel-agent/config"
)

func TestPluginCollectValuesCommand(t *testing.T) {
	g := &pluginGenerator{Config: &config.MetricPlugin{
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
	g := &pluginGenerator{Config: &config.MetricPlugin{
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
		Config: &config.MetricPlugin{
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
		Config: &config.MetricPlugin{
			Command: "echo \"just.echo.1\t1\t1397822016\"",
		},
	}

	err = generatorWithoutConf.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}

	generatorWithBadVersion := &pluginGenerator{
		Config: &config.MetricPlugin{
			Command: `echo "# mackerel-agent-plugin version=666"`,
		},
	}

	err = generatorWithBadVersion.loadPluginMeta()
	if err == nil {
		t.Error("should raise error")
	}
}
