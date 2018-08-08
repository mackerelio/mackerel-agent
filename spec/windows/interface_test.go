// +build windows

package windows

import (
	"testing"
)

func TestInterfaceGenerate(t *testing.T) {
	g := &InterfaceGenerator{}
	interfaces, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if len(interfaces) == 0 {
		t.Skipf("testing environment has no interface")
		return
	}

	iface := interfaces[0]
	if iface.Name == "" {
		t.Error("interface should have name")
	}
	if len(iface.IPv4Addresses) == 0 {
		t.Error("interface should have IPv4Addresses")
	}
	if iface.MacAddress == "" {
		t.Error("interface should have macAddress")
	}
}
