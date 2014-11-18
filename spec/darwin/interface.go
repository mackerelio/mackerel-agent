// +build darwin

package darwin

type InterfaceGenerator struct {
}

func (g *InterfaceGenerator) Key() string {
	return "interface"
}

func (g *InterfaceGenerator) Generate() (interface{}, error) {
	var interfaces map[string]map[string]interface{}

	// TODO

	return interfaces, nil
}
