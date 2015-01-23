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
	filesystems, err := windows.CollectDfValues()
	if err != nil {
		return nil, err
	}

	logger.Debugf("%q", filesystems)

	ret := make(map[string]float64)
	for name, values := range filesystems {
		if matches := regexp.MustCompile(`^(.*):`).FindStringSubmatch(name); matches != nil {
			device := regexp.MustCompile(`[^A-Za-z0-9_-]`).ReplaceAllString(matches[1], "_")

			for key, value := range values {
				floatValue, valueTypeOk := value.(float64)
				if valueTypeOk {
					if(key == "kb_size") {
						// kilo bytes -> bytes
						ret["filesystem."+device+".size"] = floatValue * 1024
					} else if(key == "kb_used") {
						// kilo bytes -> bytes
						ret["filesystem."+device+".used"] = floatValue * 1024
					}
				}
			}
		}
	}

	logger.Debugf("%q", ret)

	return metrics.Values(ret), nil
}
