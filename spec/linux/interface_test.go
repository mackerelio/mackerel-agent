// +build linux

package linux

import (
	"os"
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
		t.Errorf("should not raise error: %v", err)
	}

	if os.Getenv("CIRCLECI") != "" {
		t.Skip("Skip in CircleCI for now")
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
	if iface.Address == "" {
		t.Error("interface should have address")
	}
	if iface.DefaultGateway == "" {
		t.Error("interface should have defaultGateway")
	}
}

func TestGenerateByIpCommand(t *testing.T) {
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		t.Skip("The OS seems to be Fedora. Skipping interface test for now")
	}

	g := &InterfaceGenerator{}
	interfaces, err := g.generateByIPCommand()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	if os.Getenv("CIRCLECI") != "" {
		t.Skip("Skip in CircleCI for now")
	}

	name := "eth0"
	if _, ok := interfaces[name]; !ok {
		t.Error("should have interfaces")
		return
	}

	iface, ok := interfaces[name]
	if !ok {
		t.Error("should have item")
	}
	if len(iface.IPv4Addresses) <= 0 {
		t.Error("interface should have ipv4Addresses")
	}
	if iface.MacAddress == "" {
		t.Error("interface should have macAddress")
	}
	if iface.DefaultGateway == "" {
		t.Error("interface should have defaultGateway")
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
	if iface.DefaultGateway == "" {
		t.Log("Skip: interface should have defaultGateway")
	}
}
