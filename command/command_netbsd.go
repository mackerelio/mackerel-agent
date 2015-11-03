package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsNetbsd "github.com/mackerelio/mackerel-agent/metrics/netbsd"
	"github.com/mackerelio/mackerel-agent/spec"
	specNetbsd "github.com/mackerelio/mackerel-agent/spec/netbsd"
)

func specGenerators() []spec.Generator {
	return []spec.Generator{
		&specNetbsd.KernelGenerator{},
		&specNetbsd.MemoryGenerator{},
		&specNetbsd.CPUGenerator{},
		&specNetbsd.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.Generator {
	return &specNetbsd.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metricsNetbsd.Loadavg5Generator{},
		&metricsNetbsd.CPUUsageGenerator{},
		&metricsNetbsd.FilesystemGenerator{},
		&metricsNetbsd.MemoryGenerator{},
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
