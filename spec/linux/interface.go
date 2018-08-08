// +build linux

package linux

import (
	"os/exec"
	"regexp"
	"strings"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/spec"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

var interfaceLogger = logging.GetLogger("spec.interface")

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]spec.NetInterface, error) {
	var interfaces spec.NetInterfaces
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
	var results []spec.NetInterface
	for _, iface := range interfaces {
		if iface.Encap == "" || iface.Encap == "Loopback" {
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

func (g *InterfaceGenerator) generateByIPCommand() (spec.NetInterfaces, error) {
	interfaces := make(spec.NetInterfaces)
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
				interfaces.SetEncap(name, translateEncap(matches[1]))
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

	for _, family := range []string{"inet", "inet6"} {
		// ip -f inet route show
		out, err := exec.Command("ip", "-f", family, "route", "show").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run ip command (skip this spec): %s", err)
			return interfaces, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			name, addr, v6addr, defaultGateway := parseIProuteLine(line)
			if name == "" {
				continue
			}
			if addr != "" {
				interfaces.SetAddress(name, addr)
			}
			if v6addr != "" {
				interfaces.SetV6Address(name, v6addr)
			}
			if defaultGateway != "" {
				interfaces.SetDefaultGateway(name, defaultGateway)
			}
		}
	}
	return interfaces, nil
}

var (
	// ex.) 10.0.3.0/24 dev eth0  proto kernel  scope link  src 10.0.4.7
	// ex.) fe80::/64 dev eth0  proto kernel  metric 256
	ipRouteLineReg   = regexp.MustCompile(`^([^\s]+)\s(.*)$`)
	ipRouteDeviceReg = regexp.MustCompile(`\bdev\s+([^\s]+)\b`)
	// ex.) 10.0.3.0/24
	ipRouteIPv4Reg = regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(\/(\d{1,2}))?`)
	// ex.) fe80::/64
	ipRouteIPv6Reg = regexp.MustCompile(`([a-f0-9\:]+)\/(\d+)`)
	// ex.) default via 10.0.3.1 dev eth0
	ipRouteDefaultGatewayReg = regexp.MustCompile(`\bvia\s+([^\s]+)\b`)
)

func parseIProuteLine(line string) (name, addr, v6addr, defaultGateway string) {
	matches := ipRouteLineReg.FindStringSubmatch(line)
	if matches == nil && len(matches) < 3 {
		return
	}
	if matches := ipRouteDeviceReg.FindStringSubmatch(matches[2]); matches != nil {
		name = matches[1]
	} else {
		return
	}
	if matches := ipRouteIPv4Reg.FindStringSubmatch(matches[1]); matches != nil {
		addr = matches[1]
	}
	if matches := ipRouteIPv6Reg.FindStringSubmatch(matches[1]); matches != nil {
		v6addr = matches[1]
	}
	if matches := ipRouteDefaultGatewayReg.FindStringSubmatch(matches[2]); matches != nil {
		defaultGateway = matches[1]
	}
	return
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

func (g *InterfaceGenerator) generateByIfconfigCommand() (spec.NetInterfaces, error) {
	interfaces := make(spec.NetInterfaces)
	name := ""

	{
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
			if matches := ifconfigEncapReg.FindStringSubmatch(line); matches != nil {
				encap := matches[1]
				if encap == "" {
					encap = matches[2]
				}
				interfaces.SetEncap(name, translateEncap(encap))
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
	}

	{
		// route -n
		out, err := exec.Command("route", "-n").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run route command (skip this spec): %s", err)
			return interfaces, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			name, defaultGateway := parseRouteLine(line)
			if name == "" {
				continue
			}
			interfaces.SetDefaultGateway(name, defaultGateway)
		}
	}

	{
		// arp -an
		out, err := exec.Command("arp", "-an").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run arp command (skip this spec): %s", err)
			return interfaces, err
		}
		for _, line := range strings.Split(string(out), "\n") {
			name, addr := parseArpLine(line)
			if name == "" {
				continue
			}
			interfaces.SetAddress(name, addr)
		}
	}
	return interfaces, nil
}

// Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
// 0.0.0.0         10.0.3.1        0.0.0.0         UG    0      0        0 eth0
func parseRouteLine(line string) (name, defaultGateway string) {
	if !strings.HasPrefix(line, "0.0.0.0") {
		return
	}
	routeResults := strings.Fields(line)
	if len(routeResults) < 8 {
		return
	}
	return routeResults[7], routeResults[1]
}

// ex.) ? (10.0.3.2) at 01:23:45:67:89:ab [ether] on eth0
var arpRegexp = regexp.MustCompile(`^\S+ \((\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\) at ([a-fA-F0-9\:]+) \[(\w+)\] on ([0-9a-zA-Z\.\:\-]+)`)

func parseArpLine(line string) (name, addr string) {
	if matches := arpRegexp.FindStringSubmatch(line); matches != nil {
		return matches[4], matches[1]
	}
	return
}

func translateEncap(encap string) string {
	switch encap {
	case "Local Loopback", "loopback":
		return "Loopback"
	case "Point-to-Point Protocol":
		return "PPP"
	case "Serial Line IP":
		return "SLIP"
	case "VJ Serial Line IP":
		return "VJSLIP"
	case "IPIP Tunnel":
		return "IPIP"
	case "IPv6-in-IPv4":
		return "6to4"
	case "ether":
		return "Ethernet"
	}
	return encap
}
