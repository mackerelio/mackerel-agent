// +build netbsd

package netbsd

import (
	"github.com/mackerelio/mackerel-agent/spec"
	mkr "github.com/mackerelio/mackerel-client-go"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]mkr.Interface, error) {
	// TODO
	return []mkr.Interface{}, nil
}
