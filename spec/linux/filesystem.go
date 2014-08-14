// +build linux

package linux

import (
	"github.com/mackerelio/mackerel-agent/logging"
	utilLinux "github.com/mackerelio/mackerel-agent/util/linux"
)

type FilesystemGenerator struct {
}

func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var logger = logging.GetLogger("spec.filesystem")

var dfColumnSpecs = []utilLinux.DfColumnSpec{
	utilLinux.DfColumnSpec{"kb_size", true},
	utilLinux.DfColumnSpec{"kb_used", true},
	utilLinux.DfColumnSpec{"kb_available", true},
	utilLinux.DfColumnSpec{"percent_used", false},
	utilLinux.DfColumnSpec{"mount", false},
}

func (g *FilesystemGenerator) Generate() (interface{}, error) {
	filesystems, err := utilLinux.CollectDfValues(dfColumnSpecs)
	if err != nil {
		return nil, err
	}
	return filesystems, nil
}
