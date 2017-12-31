package command

import (
	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/metrics"
	metricsFreebsd "github.com/mackerelio/mackerel-agent/metrics/freebsd"
	"github.com/mackerelio/mackerel-agent/spec"
	specFreebsd "github.com/mackerelio/mackerel-agent/spec/freebsd"
)

func specGenerators() []spec.Generator {
	return []spec.Generator{
		&specFreebsd.KernelGenerator{},
		&specFreebsd.MemoryGenerator{},
		&specFreebsd.CPUGenerator{},
		&spec.FilesystemGenerator{},
	}
}

func interfaceGenerator() spec.InterfaceGenerator {
	return &specFreebsd.InterfaceGenerator{}
}

func metricsGenerators(conf *config.Config) []metrics.Generator {
	generators := []metrics.Generator{
		&metrics.Loadavg5Generator{},
		&metricsFreebsd.CPUUsageGenerator{},
		&metrics.FilesystemGenerator{IgnoreRegexp: conf.Filesystems.Ignore.Regexp, UseMountpoint: conf.Filesystems.UseMountpoint},
		&metricsFreebsd.MemoryGenerator{},
		&metrics.InterfaceGenerator{Interval: metricsInterval},
	}

	return generators
}
