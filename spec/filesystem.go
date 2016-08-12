// +build !windows

package spec

import (
	"fmt"

	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// Key XXX
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

// Generate XXX
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	filesystems, err := util.CollectDfValues()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]map[string]interface{})
	for _, v := range filesystems {
		entry := map[string]interface{}{
			"kb_size":      v.Blocks,
			"kb_used":      v.Used,
			"kb_available": v.Available,
			"percent_used": fmt.Sprintf("%d%%", v.Capacity),
			"mount":        v.Mounted,
		}
		ret[v.Name] = entry
	}
	return ret, nil
}
