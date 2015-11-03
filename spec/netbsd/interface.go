// +build netbsd

package netbsd

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

// Key XXX
func (g *InterfaceGenerator) Key() string {
	return "interface"
}

// Generate XXX
func (g *InterfaceGenerator) Generate() (interface{}, error) {
	var interfaces map[string]map[string]interface{}

	// TODO

	return interfaces, nil
}
