// +build !windows

package spec

import "github.com/mackerelio/mackerel-agent/util"

// FilesystemGenerator XXX
type FilesystemGenerator struct {
}

// Key XXX
func (g *FilesystemGenerator) Key() string {
	return "filesystem"
}

var dfColumnSpecs = []util.DfColumnSpec{
	{Name: "kb_size", IsInt: true},
	{Name: "kb_used", IsInt: true},
	{Name: "kb_available", IsInt: true},
	{Name: "percent_used", IsInt: false},
	{Name: "mount", IsInt: false},
}

// Generate XXX
func (g *FilesystemGenerator) Generate() (interface{}, error) {
	return util.CollectDfValues(dfColumnSpecs)
}
