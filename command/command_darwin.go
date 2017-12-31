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
		&metrics.Loadavg5Generator{},
		&metricsDarwin.CPUUsageGenerator{},
		&metricsDarwin.MemoryGenerator{},
		&metrics.FilesystemGenerator{IgnoreRegexp: conf.Filesystems.Ignore.Regexp, UseMountpoint: conf.Filesystems.UseMountpoint},
		&metrics.InterfaceGenerator{Interval: metricsInterval},
	}

	return generators
}
