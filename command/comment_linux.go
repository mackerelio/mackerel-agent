package command

import (
	"github.com/mackerelio/mackerel-agent/mackerel"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsLinux "github.com/mackerelio/mackerel-agent/metrics/linux"
	"github.com/mackerelio/mackerel-agent/spec"
	specLinux "github.com/mackerelio/mackerel-agent/spec/linux"
)

func metaGenerators() []spec.Generator {
	return []spec.Generator{
		&specLinux.KernelGenerator{},
		&specLinux.CPUGenerator{},
		&specLinux.MemoryGenerator{},
		&specLinux.BlockDeviceGenerator{},
		&specLinux.FilesystemGenerator{},
		&specLinux.InterfaceGenerator{},
	}
}

func metricsGenerators(config mackerel.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metricsLinux.Loadavg5Generator{},
		&metricsLinux.CpuusageGenerator{Interval: 60},
		&metricsLinux.MemoryGenerator{},
		&metricsLinux.UptimeGenerator{},
		&metricsLinux.InterfaceGenerator{Interval: 60},
		&metricsLinux.DiskGenerator{Interval: 60},
	}
	for _, pluginConfig := range config.Plugin["metrics"] {
		generators = append(generators, &metricsLinux.PluginGenerator{pluginConfig})
	}

	return generators
}
