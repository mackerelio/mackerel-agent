package metrics

import (
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/mackerel"
)

// ContainerGenerator definition
type ContainerGenerator struct {
	PluginGenerator
}

var containerLogger = logging.GetLogger("metrics.container")

// GetContainerGenerator returns ContainerGenerator
func GetContainerGenerator(conf *config.Config) *ContainerGenerator {
	switch conf.ContainerPlatform {
	case config.ContainerPlatformNone:
		return nil
	case config.ContainerPlatformECS:
		// TODO
		// return &ContainerGenerator{&ECSGenerator{}}
		return nil
	case config.ContainerPlatformECSFargate:
		return &ContainerGenerator{&ECSFargateGenerator{}}
	case config.ContainerPlatformKubernetes:
		// TODO
		// return &ContainerGenerator{&KubernetesGenerator{}}
		return nil
	default:
		return nil
	}
}

type ECSFargateGenerator struct{}

func (g *ECSFargateGenerator) Generate() (Values, error) {
	v := make(Values)
	// Do nothing
	return v, nil
}

func (g *ECSFargateGenerator) CustomIdentifier() *string {
	return nil
}

func (g *ECSFargateGenerator) PrepareGraphDefs() ([]mackerel.CreateGraphDefsPayload, error) {
	meta := &pluginMeta{
		Graphs: map[string]customGraphDef{
			"ecs.task.cpu": customGraphDef{
				Label: "CPU%",
				Unit:  "percentage",
				Metrics: []customGraphMetricDef{
					{Name: "usage", Label: "cpu%", Stacked: false},
				},
			},
			"ecs.task.memory": customGraphDef{
				Label: "Memory Usage",
				Unit:  "bytes",
				Metrics: []customGraphMetricDef{
					{Name: "usage", Label: "usage", Stacked: false},
				},
			},
			"ecs.container.#.cpu": customGraphDef{
				Label: "CPU% / Container",
				Unit:  "percentage",
				Metrics: []customGraphMetricDef{
					{Name: "usage", Label: "cpu%", Stacked: true},
				},
			},
			"ecs.container.#.memory": customGraphDef{
				Label: "Memory Usage / Container",
				Unit:  "bytes",
				Metrics: []customGraphMetricDef{
					{Name: "usage", Label: "usage", Stacked: true},
				},
			},
		},
	}
	return makeCreateGraphDefsPayload(meta), nil
}
