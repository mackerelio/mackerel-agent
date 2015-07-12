package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsFreebsd "github.com/mackerelio/mackerel-agent/metrics/freebsd"
	"github.com/mackerelio/mackerel-agent/spec"
	specFreebsd "github.com/mackerelio/mackerel-agent/spec/freebsd"
)

func specGenerators() []spec.Generator {
	specs := []spec.Generator{
		&specFreebsd.KernelGenerator{},
		&specFreebsd.MemoryGenerator{},
		&specFreebsd.CPUGenerator{},
		&specFreebsd.FilesystemGenerator{},
	}
	cloudGenerator, err := spec.NewCloudGenerator("")
	if err != nil {
		logger.Errorf("Failed to create cloudGenerator: %s", err.Error())
	} else {
		specs = append(specs, cloudGenerator)
	}
	return specs
}

func interfaceGenerator() spec.Generator {
	return &specFreebsd.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metricsFreebsd.Loadavg5Generator{},
		&metricsFreebsd.CPUUsageGenerator{},
		&metricsFreebsd.FilesystemGenerator{},
		&metricsFreebsd.MemoryGenerator{},
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
