// +build windows

package windows

import (
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util/windows"
	"regexp"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// NewFilesystemGenerator XXX
func NewFilesystemGenerator() (*FilesystemGenerator, error) {
	return &FilesystemGenerator{}, nil
}

var logger = logging.GetLogger("metrics.filesystem")

// Generate XXX
func (g *FilesystemGenerator) Generate() (metrics.Values, error) {
	filesystems, err := windows.CollectFilesystemValues()
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for name, values := range filesystems {
		if matches := regexp.MustCompile(`^(.*):`).FindStringSubmatch(name); matches != nil {
			device := regexp.MustCompile(`[^A-Za-z0-9_-]`).ReplaceAllString(matches[1], "_")

			ret["filesystem."+device+".size"] = values.KbSize * 1024
			ret["filesystem."+device+".used"] = values.KbUsed * 1024
		}
	}

	logger.Debugf("%q", ret)

	return metrics.Values(ret), nil
}
