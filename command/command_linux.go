package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsLinux "github.com/mackerelio/mackerel-agent/metrics/linux"
	"github.com/mackerelio/mackerel-agent/spec"
	specLinux "github.com/mackerelio/mackerel-agent/spec/linux"
)

func specGenerators() []spec.Generator {
	return []spec.Generator{
		&specLinux.KernelGenerator{},
		&specLinux.CPUGenerator{},
		&specLinux.MemoryGenerator{},
		&specLinux.BlockDeviceGenerator{},
		&specLinux.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.Generator {
	return &specLinux.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metricsLinux.Loadavg5Generator{},
		&metricsLinux.CPUUsageGenerator{Interval: metricsInterval},
		&metricsLinux.MemoryGenerator{},
		&metricsLinux.UptimeGenerator{},
		&metricsLinux.InterfaceGenerator{Interval: metricsInterval},
		&metricsLinux.DiskGenerator{Interval: metricsInterval},
		&metricsLinux.FilesystemGenerator{Ignore: conf.Filesystems.Ignore},
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
