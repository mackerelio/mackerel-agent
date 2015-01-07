// +build darwin

package darwin

import (
	"github.com/mackerelio/mackerel-agent/logging"
	utilDarwin "github.com/mackerelio/mackerel-agent/util/darwin"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// Key XXX
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var logger = logging.GetLogger("spec.filesystem")

var dfColumnSpecs = []utilDarwin.DfColumnSpec{
	utilDarwin.DfColumnSpec{Name: "kb_size", IsInt: true},
	utilDarwin.DfColumnSpec{Name: "kb_used", IsInt: true},
	utilDarwin.DfColumnSpec{Name: "kb_available", IsInt: true},
	utilDarwin.DfColumnSpec{Name: "percent_used", IsInt: false},
	utilDarwin.DfColumnSpec{Name: "mount", IsInt: false},
}

// Generate XXX
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	return utilDarwin.CollectDfValues(dfColumnSpecs)
}
