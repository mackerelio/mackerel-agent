// +build linux

package linux

import (
	"os/exec"
	"regexp"
	"strings"

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
	_, err := exec.LookPath("ip")
	// has ip command
	if err == nil {
		interfaces, err = g.generateByIPCommand()
		if err != nil {
			return nil, err
		}
	} else {
		interfaces, err = g.generateByIfconfigCommand()
		if err != nil {
			return nil, err
		}
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
	// ex.) 2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
	ipCmdNameReg = regexp.MustCompile(`^(\d+): ([0-9a-zA-Z@:\.\-_]*?)(@[0-9a-zA-Z]+|):\s`)
	// ex.) link/ether 12:34:56:78:9a:bc brd ff:ff:ff:ff:ff:ff
	ipCmdEncapReg = regexp.MustCompile(`link\/(\w+) ([\da-f\:]+) `)
	// ex.) inet 10.0.4.7/24 brd 10.0.5.255 scope global eth0
	ipCmdIPv4Reg = regexp.MustCompile(`inet (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(\/(\d{1,2}))?`)
	//inet6 fe80::44b3:b3ff:fe1c:d17c/64 scope link
	ipCmdIPv6Reg = regexp.MustCompile(`inet6 ([a-f0-9\:]+)\/(\d+) scope (\w+)`)
)

func (g *InterfaceGenerator) generateByIPCommand() (spec.Interfaces, error) {
	interfaces := make(spec.Interfaces)
	name := ""
	{
		// ip addr
		out, err := exec.Command("ip", "addr").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run ip command (skip this spec): %s", err)
			return nil, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			if matches := ipCmdNameReg.FindStringSubmatch(line); matches != nil {
				name = matches[2]
			}
			if matches := ipCmdEncapReg.FindStringSubmatch(line); matches != nil {
				interfaces.SetMacAddress(name, matches[2])
			}
			if matches := ipCmdIPv4Reg.FindStringSubmatch(line); matches != nil {
				interfaces.AppendIPv4Address(name, matches[1])
			}
			if matches := ipCmdIPv6Reg.FindStringSubmatch(line); matches != nil {
				interfaces.AppendIPv6Address(name, matches[1])
			}
		}
	}
	return interfaces, nil
}

var (
	// ex.) eth0      Link encap:Ethernet  HWaddr 12:34:56:78:9a:bc
	ifconfigNameReg = regexp.MustCompile(`^([0-9a-zA-Z@\.\:\-_]+)\s+`)
	// ex.) eth0      Link encap:Ethernet  HWaddr 12:34:56:78:9a:bc
	ifconfigEncapReg = regexp.MustCompile(`Link encap:(Local Loopback)|Link encap:(.+?)\s`)
	// ex.) eth0      Link encap:Ethernet  HWaddr 00:16:3e:4f:f3:41
	ifconfigMacAddrReg = regexp.MustCompile(`HWaddr (.+?)\s`)
	// ex.) inet addr:10.0.4.7  Bcast:10.0.5.255  Mask:255.255.255.0
	ifconfigV4AddrReg = regexp.MustCompile(`inet addr:(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
	// ex.) inet6 addr: fe80::44b3:b3ff:fe1c:d17c/64 Scope:Link
	ifconfigV6AddrReg = regexp.MustCompile(`inet6 addr: ([a-f0-9\:]+)\/(\d+) Scope:(\w+)`)
)

func (g *InterfaceGenerator) generateByIfconfigCommand() (spec.Interfaces, error) {
	interfaces := make(spec.Interfaces)
	name := ""

	// ifconfig -a
	out, err := exec.Command("ifconfig", "-a").Output()
	if err != nil {
		interfaceLogger.Errorf("Failed to run ifconfig command (skip this spec): %s", err)
		return nil, err
	}
	for _, line := range strings.Split(string(out), "\n") {
		if matches := ifconfigNameReg.FindStringSubmatch(line); matches != nil {
			name = matches[1]
		}
		if matches := ifconfigMacAddrReg.FindStringSubmatch(line); matches != nil {
			interfaces.SetMacAddress(name, matches[1])
		}
		if matches := ifconfigV4AddrReg.FindStringSubmatch(line); matches != nil {
			interfaces.AppendIPv4Address(name, matches[1])
		}
		if matches := ifconfigV6AddrReg.FindStringSubmatch(line); matches != nil {
			interfaces.AppendIPv6Address(name, matches[1])
		}
	}
	return interfaces, nil
}
