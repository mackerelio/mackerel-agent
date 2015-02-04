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

	ret := make(map[string]map[string]interface{})

	fileSystems, err := windows.CollectFilesystemValues()
	for drive, f := range fileSystems {
		ret[drive] = map[string]interface{}{
			"percent_used": f.Percent_used,
			"kb_used":      f.Kb_used,
			"kb_size":      f.Kb_size,
			"kb_available": f.Kb_available,
			"mount":        f.Mount,
			"label":        f.Label,
			"volume_name":  f.Volume_name,
			"fs_type":      f.Fs_type,
		}
	}

	return ret, err
}
