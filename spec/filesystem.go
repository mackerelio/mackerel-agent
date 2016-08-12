// +build !windows

package spec

import (
	"fmt"

	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator generator for filesystems implements spec.Generator interface
type FilesystemGenerator struct {
}

// Key key name of the generator for satisfying spec.Generator interface
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

// Generate specs of filesystems
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	filesystems, err := util.CollectDfValues()
	if err != nil {
		return nil, err
	}
	ret := make(map[string]map[string]interface{})
	for _, v := range filesystems {
		ret[v.Name] = map[string]interface{}{
			"kb_size":      v.Blocks,
			"kb_used":      v.Used,
			"kb_available": v.Available,
			"percent_used": fmt.Sprintf("%d%%", v.Capacity),
			"mount":        v.Mounted,
		}
	}
	return ret, nil
}
