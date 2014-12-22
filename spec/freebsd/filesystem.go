// +build freebsd

package freebsd

import (
	"github.com/mackerelio/mackerel-agent/logging"
	utilFreebsd "github.com/mackerelio/mackerel-agent/util/freebsd"
)

type FilesystemGenerator struct {
}

func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var logger = logging.GetLogger("spec.filesystem")

var dfColumnSpecs = []utilFreebsd.DfColumnSpec{
	utilFreebsd.DfColumnSpec{Name: "kb_size", IsInt: true},
	utilFreebsd.DfColumnSpec{Name: "kb_used", IsInt: true},
	utilFreebsd.DfColumnSpec{Name: "kb_available", IsInt: true},
	utilFreebsd.DfColumnSpec{Name: "percent_used", IsInt: false},
	utilFreebsd.DfColumnSpec{Name: "mount", IsInt: false},
}

func (g *FilesystemGenerator) Generate() (interface{}, error) {
	return utilFreebsd.CollectDfValues(dfColumnSpecs)
}
