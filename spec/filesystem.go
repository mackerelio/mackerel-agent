// +build !windows

package spec

import (
	"fmt"

	"github.com/mackerelio/mackerel-client-go"

	"github.com/mackerelio/mackerel-agent/util"
)

// FilesystemGenerator generates filesystem spec.
type FilesystemGenerator struct {
}

// Generate specs of filesystems.
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	filesystems, err := util.CollectDfValues()
	if err != nil {
		return nil, err
	}
	ret := make(mackerel.FileSystem)
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
