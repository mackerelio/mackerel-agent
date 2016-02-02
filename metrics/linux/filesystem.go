// +build linux

package linux

import (
	"regexp"

	"github.com/mackerelio/mackerel-agent/config"
	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
	Ignore config.Regexpwrapper
}

var logger = logging.GetLogger("metrics.filesystem")

var dfColumnSpecs = []util.DfColumnSpec{
	util.DfColumnSpec{Name: "size", IsInt: true},
	util.DfColumnSpec{Name: "used", IsInt: true},
}

// Generate XXX
func (g *FilesystemGenerator) Generate() (metrics.Values, error) {
	filesystems, err := util.CollectDfValues(dfColumnSpecs)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for name, values := range filesystems {
		// https://github.com/docker/docker/blob/v1.5.0/daemon/graphdriver/devmapper/deviceset.go#L981
		if regexp.MustCompile(`^/dev/mapper/docker-`).FindStringSubmatch(name) != nil ||
			(g.Ignore.Regexp != nil && g.Ignore.Regexp.FindStringSubmatch(name) != nil) {
			continue
		}
		if matches := regexp.MustCompile(`^/dev/(.*)$`).FindStringSubmatch(name); matches != nil {
			device := regexp.MustCompile(`[^A-Za-z0-9_-]`).ReplaceAllString(matches[1], "_")
			for key, value := range values {
				intValue, valueTypeOk := value.(int64)
				if valueTypeOk {
					// kilo bytes -> bytes
					ret["filesystem."+device+"."+key] = float64(intValue) * 1024
				}
			}
		}
	}

	return metrics.Values(ret), nil
}
