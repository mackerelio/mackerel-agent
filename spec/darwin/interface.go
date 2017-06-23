// +build darwin

package darwin

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/spec"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

// Key XXX
func (g *InterfaceGenerator) Key() string {
	return "interface"
}

var interfaceLogger = logging.GetLogger("spec.interface")

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]spec.NetInterface, error) {
	var interfaces spec.NetInterfaces

	interfaces, err := g.generateByIfconfigCommand()
	if err != nil {
		return nil, err
	}
	var results []spec.NetInterface
	for _, iface := range interfaces {
		if iface.Encap == "Loopback" {
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

func (g *InterfaceGenerator) generateByIfconfigCommand() (spec.NetInterfaces, error) {
	interfaces := make(spec.NetInterfaces)

	{
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
			// ex.) lo0: flags=8049<UP,LOOPBACK,RUNNING,MULTICAST> mtu 16384
			if strings.Contains(line, "LOOPBACK") {
				interfaces.SetEncap(name, "Loopback")
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
	}

	{
		// netstat -f inet -rn
		out, err := exec.Command("netstat", "-f", "inet", "-rn").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run netstat command (skip this spec): %s", err)
			return interfaces, err
		}
		name, gateway := retrieveDefaultGateway(out)
		interfaces.SetDefaultGateway(name, gateway)
	}
	return interfaces, nil
}

var routeRegexp = regexp.MustCompile(`^default\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)

func retrieveDefaultGateway(out []byte) (name, gateway string) {
	lineScanner := bufio.NewScanner(bytes.NewReader(out))
	// ex.)
	// Routing tables

	// Internet:
	// Destination        Gateway            Flags        Refs      Use   Netif Expire
	// default            10.0.3.1      UGSc           25        0     en0
	for lineScanner.Scan() {
		line := lineScanner.Text()
		if routeRegexp.MatchString(line) {
			routeResults := strings.Fields(line)
			if len(routeResults) < 6 {
				continue
			}
			return routeResults[5], routeResults[1]
		}
	}
	return
}
