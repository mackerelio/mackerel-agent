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
		&spec.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.InterfaceGenerator {
	return &specNetbsd.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metrics.LoadavgGenerator{},
		&metricsNetbsd.CPUUsageGenerator{},
		&metrics.FilesystemGenerator{IgnoreRegexp: conf.Filesystems.Ignore.Regexp, UseMountpoint: conf.Filesystems.UseMountpoint},
		&metricsNetbsd.MemoryGenerator{},
		&metrics.InterfaceGenerator{Interval: metricsInterval},
	}

	return generators
}
