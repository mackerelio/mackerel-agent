package spec

import (
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"
)

var logger = logging.GetLogger("spec")

// Generator interface for generating spec values
type Generator interface {
	Generate() (interface{}, error)
}

// Collect spec values
func Collect(specGenerators []Generator) mackerel.HostMeta {
	var specs mackerel.HostMeta
	for _, g := range specGenerators {
		value, err := g.Generate()
		if err != nil {
			logger.Warningf("Failed to collect meta in %T (skip this spec): %s", g, err.Error())
			continue
		}
		switch v := value.(type) {
		case mackerel.BlockDevice:
			specs.BlockDevice = v
		case mackerel.CPU:
			specs.CPU = v
		case mackerel.FileSystem:
			specs.Filesystem = v
		case mackerel.Kernel:
			specs.Kernel = v
		case mackerel.Memory:
			specs.Memory = v
		case *mackerel.Cloud:
			specs.Cloud = v
		default:
		}
	}
	return specs
}
