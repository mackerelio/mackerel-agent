// +build linux

package linux

import (
	"os"
	"strings"
	"testing"

	"github.com/mackerelio/mackerel-agent/spec"
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

	name := lookupDefaultName(interfaces, "eth0")
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
}

func TestGenerateByIfconfigCommand(t *testing.T) {
	g := &InterfaceGenerator{}
	interfaces, err := g.generateByIfconfigCommand()
	if err != nil {
		t.Log("Skip: should not raise error")
	}

	name := lookupDefaultName(interfaces, "eth0")
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

// lookupDefaultName returns network interface name that seems to be default NIC.
// There are some naming rules on recent linux environment.
// 1. traditional names (eth0, eth1)
// 2. predictable names for ethernet (ens0, enp1s0)
// 3. predictable names for wireless (wls0, wls1)
//
// There is type-differed version at metric/interface_test.go.
func lookupDefaultName(ifaces spec.Interfaces, fallback string) string {
	for key := range ifaces {
		switch {
		case strings.HasPrefix(key, "eth"):
			return key
		case strings.HasPrefix(key, "en"):
			return key
		case strings.HasPrefix(key, "wl"):
			return key
		}
	}
	return fallback
}
