// +build windows

package windows

import (
	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/util/windows"
)

// FilesystemGenerator generates filesystem spec.
type FilesystemGenerator struct {
}

var filesystemLogger = logging.GetLogger("spec.filesystem")

// Generate specs of filesystems.
func (g *FilesystemGenerator) Generate() (interface{}, error) {

	ret := make(mackerel.FileSystem)

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
