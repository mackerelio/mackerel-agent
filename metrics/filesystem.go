//go:build !windows

package metrics

import (
	"regexp"
	"strings"

	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator is common filesystem metrics generator on unix os.
type FilesystemGenerator struct {
	IgnoreRegexp  *regexp.Regexp
	UseMountpoint bool
}

// Generate the metrics of filesystems
func (g *FilesystemGenerator) Generate() (Values, error) {
	filesystems, err := util.CollectDfValues()
	if err != nil {
		return nil, err
	}
	ret := Values{}
	for _, dfs := range filesystems {
		name := dfs.Name
		if g.IgnoreRegexp != nil && g.IgnoreRegexp.MatchString(name) {
			continue
		}
		if device, ok := strings.CutPrefix(name, "/dev/"); ok {
			var metricName string
			if g.UseMountpoint {
				metricName = util.SanitizeMetricKey(dfs.Mounted)
			} else {
				metricName = util.SanitizeMetricKey(device)
			}
			// kilo bytes -> bytes
			ret["filesystem."+metricName+".size"] = float64(dfs.Used+dfs.Available) * 1024
			ret["filesystem."+metricName+".used"] = float64(dfs.Used) * 1024
		}
	}
	return ret, nil
}
