package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsWindows "github.com/mackerelio/mackerel-agent/metrics/windows"
	"github.com/mackerelio/mackerel-agent/spec"
	specWindows "github.com/mackerelio/mackerel-agent/spec/windows"
)

func specGenerators() []spec.Generator {
	return []spec.Generator{
		&specWindows.KernelGenerator{},
		&specWindows.CPUGenerator{},
		&specWindows.MemoryGenerator{},
		&specWindows.BlockDeviceGenerator{},
		&specWindows.FilesystemGenerator{},
		&specWindows.InterfaceGenerator{},
	}
}

func metricsGenerators(conf config.Config) []metrics.Generator {
	impls := []metrics.Generator{
		metricsWindows.NewLoadavg5Generator(),
		metricsWindows.NewCpuusageGenerator(60),
		metricsWindows.NewMemoryGenerator(),
		metricsWindows.NewUptimeGenerator(),
		metricsWindows.NewInterfaceGenerator(60),
		metricsWindows.NewDiskGenerator(60),
	}
	generators := []metrics.Generator{}
	for _, generator := range impls {
		if generator != nil {
			generators = append(generators, generator)
		}
	}
	for _, pluginConfig := range conf.Plugin["metrics"] {
		generators = append(generators, metricsWindows.NewPluginGenerator(pluginConfig))
	}

	return generators
}
