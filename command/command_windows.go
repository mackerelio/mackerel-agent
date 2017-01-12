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

func interfaceGenerator() spec.InterfaceGenerator {
	return &specWindows.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	var g metrics.Generator
	var err error

	generators := []metrics.Generator{}
	if g, err = metricsWindows.NewProcessorQueueLengthGenerator(); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewCPUUsageGenerator(); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewMemoryGenerator(); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewFilesystemGenerator(conf.Filesystems.Ignore.Regexp); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewInterfaceGenerator(metricsInterval); err == nil {
		generators = append(generators, g)
	}
	if g, err = metricsWindows.NewDiskGenerator(metricsInterval); err == nil {
		generators = append(generators, g)
	}

	return generators
}
