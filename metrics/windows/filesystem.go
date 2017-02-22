// +build windows

package windows

import (
	"regexp"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util"
	"github.com/mackerelio/mackerel-agent/util/windows"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
	IgnoreRegexp *regexp.Regexp
}

// NewFilesystemGenerator XXX
func NewFilesystemGenerator(ignoreReg *regexp.Regexp) (*FilesystemGenerator, error) {
	return &FilesystemGenerator{IgnoreRegexp: ignoreReg}, nil
}

var logger = logging.GetLogger("metrics.filesystem")

var driveLetterReg = regexp.MustCompile(`^(.*):`)

// Generate the metrics of filesystems
func (g *FilesystemGenerator) Generate() (metrics.Values, error) {
	filesystems, err := windows.CollectFilesystemValues()
	if err != nil {
		return nil, err
	}
	ret := metrics.Values{}
	for name, values := range filesystems {
		if g.IgnoreRegexp != nil && g.IgnoreRegexp.MatchString(name) {
			continue
		}
		if matches := driveLetterReg.FindStringSubmatch(name); matches != nil {
			device := util.SanitizeMetricKey(matches[1])
			ret["filesystem."+device+".size"] = values.KbSize * 1024
			ret["filesystem."+device+".used"] = values.KbUsed * 1024
		}
	}
	logger.Debugf("%q", ret)
	return ret, nil
}
