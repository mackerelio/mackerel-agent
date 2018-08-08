// +build freebsd

package freebsd

import "github.com/mackerelio/mackerel-agent/spec"

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]spec.NetInterface, error) {
	// TODO
	return []spec.NetInterface{}, nil
}
