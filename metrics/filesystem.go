// +build !windows

package metrics

import (
	"regexp"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator is common filesystem metrics generator on unix os.
type FilesystemGenerator struct {
	IgnoreRegexp *regexp.Regexp
}

var logger = logging.GetLogger("metrics")

var dfColumnSpecs = []util.DfColumnSpec{
	util.DfColumnSpec{Name: "size", IsInt: true},
	util.DfColumnSpec{Name: "used", IsInt: true},
}

var (
	devDirReg    = regexp.MustCompile(`^/dev/(.*)$`)
	sanitizerReg = regexp.MustCompile(`[^A-Za-z0-9_-]`)
)

// Generate the metrics of filesystems
func (g *FilesystemGenerator) Generate() (Values, error) {
	filesystems, err := util.CollectDfValues(dfColumnSpecs)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for name, values := range filesystems {
		// https://github.com/docker/docker/blob/v1.5.0/daemon/graphdriver/devmapper/deviceset.go#L981
		if strings.HasPrefix(name, "/dev/mapper/docker-") ||
			(g.IgnoreRegexp != nil && g.IgnoreRegexp.MatchString(name)) {
			continue
		}
		if matches := devdirReg.FindStringSubmatch(name); matches != nil {
			device := sanitizerReg.ReplaceAllString(matches[1], "_")
			for key, value := range values {
				intValue, valueTypeOk := value.(int64)
				if valueTypeOk {
					// kilo bytes -> bytes
					ret["filesystem."+device+"."+key] = float64(intValue) * 1024
				}
			}
		}
	}

	return Values(ret), nil
}
