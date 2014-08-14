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
	utilLinux.DfColumnSpec{Name: "kb_size", IsInt: true},
	utilLinux.DfColumnSpec{Name: "kb_used", IsInt: true},
	utilLinux.DfColumnSpec{Name: "kb_available", IsInt: true},
	utilLinux.DfColumnSpec{Name: "percent_used", IsInt: false},
	utilLinux.DfColumnSpec{Name: "mount", IsInt: false},
}

func (g *FilesystemGenerator) Generate() (interface{}, error) {
	return utilLinux.CollectDfValues(dfColumnSpecs)
}
