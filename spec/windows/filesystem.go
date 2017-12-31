// +build windows

package windows

import (
	"github.com/mackerelio/golib/logging"
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
			"percent_used": f.PercentUsed,
			"kb_used":      f.KbUsed,
			"kb_size":      f.KbSize,
			"kb_available": f.KbAvailable,
			"mount":        f.Mount,
			"label":        f.Label,
			"volume_name":  f.VolumeName,
			"fs_type":      f.FsType,
		}
	}

	return ret, err
}
