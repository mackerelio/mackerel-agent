//go:build linux
// +build linux

package linux

import (
	"os"
	"testing"
)

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
	// In Docker enabled Travis environment, there's "docker0" interface
	// which does not have defaultGateway..
	if iface.Name == "docker0" && os.Getenv("TRAVIS") != "" {
		iface = value[1]
	}
	if len(iface.IPv4Addresses) <= 0 {
		t.Error("interface should have ipv4Addresses")
	}
	if iface.MacAddress == "" {
		t.Error("interface should have macAddress")
	}
}
