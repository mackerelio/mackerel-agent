// +build linux

package linux

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// Key XXX
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var logger = logging.GetLogger("spec.filesystem")

var dfColumnSpecs = []util.DfColumnSpec{
	util.DfColumnSpec{Name: "kb_size", IsInt: true},
	util.DfColumnSpec{Name: "kb_used", IsInt: true},
	util.DfColumnSpec{Name: "kb_available", IsInt: true},
	util.DfColumnSpec{Name: "percent_used", IsInt: false},
	util.DfColumnSpec{Name: "mount", IsInt: false},
}

// Generate XXX
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	return util.CollectDfValues(dfColumnSpecs)
}
