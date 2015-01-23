// +build windows

package windows

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// Key XX
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var filesystemLogger = logging.GetLogger("spec.filesystem")

// Generate XXX
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	return windows.CollectDfValues()
}
