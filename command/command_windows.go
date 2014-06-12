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
	}
}

func interfaceGenerator() spec.Generator {
	return &specWindows.InterfaceGenerator{}
}

func metricsGenerators(conf config.Config) []metrics.Generator {
	var g metrics.Generator
	var err error

	generators := []metrics.Generator{}
	if g, err = metricsWindows.NewLoadavg5Generator(); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewCpuusageGenerator(60); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewMemoryGenerator(); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewUptimeGenerator(); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewInterfaceGenerator(60); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewDiskGenerator(60); err == nil {
		generators = append(generators, g)
	}
	for _, pluginConfig := range conf.Plugin["metrics"] {
		if g, err = metricsWindows.NewPluginGenerator(pluginConfig); err == nil {
			generators = append(generators, g)
		}
	}

	return generators
}
