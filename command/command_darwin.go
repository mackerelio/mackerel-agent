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
		&spec.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.InterfaceGenerator {
	return &specDarwin.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metricsDarwin.Loadavg5Generator{},
		&metricsDarwin.CPUUsageGenerator{},
		&metricsDarwin.MemoryGenerator{},
		&metricsDarwin.SwapGenerator{},
		&metrics.FilesystemGenerator{IgnoreRegexp: conf.Filesystems.Ignore.Regexp, UseMountpoint: conf.Filesystems.UseMountpoint},
		&metricsDarwin.InterfaceGenerator{Interval: metricsInterval},
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
