package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsDarwin "github.com/mackerelio/mackerel-agent/metrics/darwin"
	"github.com/mackerelio/mackerel-agent/spec"
	specDarwin "github.com/mackerelio/mackerel-agent/spec/darwin"
)

func specGenerators() []spec.Generator {
	return []spec.Generator{
		&specDarwin.KernelGenerator{},
		&specDarwin.MemoryGenerator{},
		&specDarwin.CPUGenerator{},
		&specDarwin.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.Generator {
	return &specDarwin.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metricsDarwin.Loadavg5Generator{},
		&metricsDarwin.CPUUsageGenerator{},
		&metricsDarwin.MemoryGenerator{},
		&metricsDarwin.SwapGenerator{},
		&metricsDarwin.FilesystemGenerator{},
		&metricsDarwin.InterfaceGenerator{Interval: 60},
	}

	return generators
}

func pluginGenerators(conf *config.Config) []metrics.PluginGenerator {
	generators := []metrics.PluginGenerator{}

	for _, pluginConfig := range conf.Plugin["metrics"] {
		generators = append(generators, metrics.NewPluginGenerator(pluginConfig))
	}

	return generators
}
