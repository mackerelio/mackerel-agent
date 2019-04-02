// +build netbsd

package netbsd

import (
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
