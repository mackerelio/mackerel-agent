// +build darwin

package darwin

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/spec"
	mkr "github.com/mackerelio/mackerel-client-go"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

var interfaceLogger = logging.GetLogger("spec.interface")

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]mkr.Interface, error) {
	var interfaces spec.Interfaces

	interfaces, err := g.generateByIfconfigCommand()
	if err != nil {
		return nil, err
	}
	var results []mkr.Interface
	for _, iface := range interfaces {
		if spec.IsLoopback(iface) {
			continue
		}
		if len(iface.IPv4Addresses) == 0 && len(iface.IPv6Addresses) == 0 {
			continue
		}
		results = append(results, iface)
	}
	return results, nil
}

var (
	// ex.) en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
	nameReg = regexp.MustCompile(`^([0-9a-zA-Z@\-_]+):\s+`)
	// ex.) ether 10:93:00:00:00:00
	macReg = regexp.MustCompile(`^\s*ether\s+([0-9a-f:]+)`)
	// ex.) inet 10.0.3.1 netmask 0xffffff00 broadcast 10.0.3.255
	ipv4Reg = regexp.MustCompile(`^\s*inet\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+netmask\s+(0x[0-9a-f]+)`)
	// ex.) inet6 2001:268:c044:1111:1111:1111:1111:1111 prefixlen 64 autoconf
	ipv6Reg = regexp.MustCompile(`^\s*inet6\s+([0-9a-f:]+)\s+prefixlen\s+(\d+)`)
)

func (g *InterfaceGenerator) generateByIfconfigCommand() (spec.Interfaces, error) {
	interfaces := make(spec.Interfaces)

	// ifconfig -a
	out, err := exec.Command("ifconfig", "-a").Output()
	if err != nil {
		interfaceLogger.Errorf("Failed to run ifconfig command (skip this spec): %s", err)
		return nil, err
	}

	lineScanner := bufio.NewScanner(bytes.NewReader(out))
	name := ""
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if matches := nameReg.FindStringSubmatch(line); matches != nil {
			name = matches[1]
		}
		if matches := macReg.FindStringSubmatch(line); matches != nil {
			interfaces.SetMacAddress(name, matches[1])
		}
		if matches := ipv4Reg.FindStringSubmatch(line); matches != nil {
			interfaces.AppendIPv4Address(name, matches[1])
		}
		// ex.) inet6 2001:268:c044:1111:1111:1111:1111:1111 prefixlen 64 autoconf
		if matches := ipv6Reg.FindStringSubmatch(line); matches != nil {
			interfaces.AppendIPv6Address(name, matches[1])
		}
	}
	return interfaces, nil
}
