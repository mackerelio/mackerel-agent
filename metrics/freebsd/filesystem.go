// +build freebsd

package freebsd

import (
	"regexp"

	"github.com/mackerelio/mackerel-agent/logging"
	"github.com/mackerelio/mackerel-agent/metrics"
	utilFreebsd "github.com/mackerelio/mackerel-agent/util/freebsd"
)

type FilesystemGenerator struct {
}

var logger = logging.GetLogger("metrics.filesystem")

var dfColumnSpecs = []utilFreebsd.DfColumnSpec{
	utilFreebsd.DfColumnSpec{Name: "size", IsInt: true},
	utilFreebsd.DfColumnSpec{Name: "used", IsInt: true},
}

func (g *FilesystemGenerator) Generate() (metrics.Values, error) {
	filesystems, err := utilFreebsd.CollectDfValues(dfColumnSpecs)
	if err != nil {
		return nil, err
	}

	ret := make(map[string]float64)
	for name, values := range filesystems {
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
