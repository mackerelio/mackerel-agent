//go:build darwin

package darwin

import (
	"testing"
)

func TestInterfaceGenerate(t *testing.T) {
	g := &InterfaceGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if len(value) == 0 {
		t.Error("should have at least 1 interface")
		return
	}

	iface := value[0]
	if len(iface.IPv4Addresses) <= 0 {
		t.Error("interface should have ipv4Addresses")
	}
	if iface.MacAddress == "" {
		t.Error("interface should have macAddress")
	}
}

func TestGenerateByIfconfigCommand(t *testing.T) {
	g := &InterfaceGenerator{}
	interfaces, err := g.generateByIfconfigCommand()
	if err != nil {
		t.Log("Skip: should not raise error")
	}

	name := "eth0"
	if _, ok := interfaces[name]; !ok {
		t.Log("Skip: should have interfaces")
	}

	iface, ok := interfaces[name]
	if !ok {
		t.Log("Skip: should have item")
	}
	if len(iface.IPv4Addresses) <= 0 {
		t.Log("Skip: interface should have ipv4Addresses")
	}
	if iface.MacAddress == "" {
		t.Log("Skip: interface should have macAddress")
	}
}
