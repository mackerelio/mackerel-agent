// +build windows

package windows

import (
	"testing"
)

func TestInterfaceKey(t *testing.T) {
	g := &InterfaceGenerator{}

	if g.Key() != "interface" {
		t.Error("key should be interface")
	}
}

func TestInterfaceGenerate(t *testing.T) {
	g := &InterfaceGenerator{}
	value, err := g.Generate()
	if err != nil {
		t.Error("should not raise error")
	}

	interfaces, typeOk := value.([]map[string]interface{})
	if !typeOk {
		t.Errorf("value should be slice of map. %+v", value)
	}
	if len(interfaces) == 0 {
		t.Error("should have at least 1 interface")
	}

	iface := interfaces[0]
	if _, ok := iface["name"]; !ok {
		t.Error("interface should have name")
	}
	if _, ok := iface["ipAddress"]; !ok {
		t.Error("interface should have ipAddress")
	}
	if _, ok := iface["macAddress"]; !ok {
		t.Error("interface should have macAddress")
	}
}
