// +build linux

package linux

import (
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mackerelio/mackerel-agent/logging"
)

type InterfaceGenerator struct {
}

func (g *InterfaceGenerator) Key() string {
	return "interface"
}

var interfaceLogger = logging.GetLogger("spec.interface")

func (g *InterfaceGenerator) Generate() (interface{}, error) {
	var interfaces map[string]map[string]interface{}
	_, err := exec.LookPath("ip")
	// has ip command
	if err == nil {
		interfaces, err = g.GenerateByIpCommand()
		if err != nil {
			return nil, err
		}
	} else {
		interfaces, err = g.GenerateByIfconfigCommand()
		if err != nil {
			return nil, err
		}
	}

	results := make([]map[string]interface{}, 0)
	for key, iface := range interfaces {
		if iface["encap"] == nil || iface["encap"] == "Loopback" {
			continue
		}
		if iface["ipAddress"] == nil || iface["macAddress"] == nil {
			continue
		}
		iface["name"] = key
		results = append(results, iface)
	}

	return results, nil
}

func (g *InterfaceGenerator) GenerateByIpCommand() (map[string]map[string]interface{}, error) {
	interfaces := make(map[string]map[string]interface{})
	name := ""

	{
		// ip addr
		out, err := exec.Command("ip", "addr").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run ip command (skip this spec): %s", err)
			return nil, err
		}

		for _, line := range strings.Split(string(out), "\n") {
			// ex.) 2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc pfifo_fast state UP qlen 1000
			if matches := regexp.MustCompile(`^(\d+): ([0-9a-zA-Z@:\.\-_]*?)(@[0-9a-zA-Z]+|):\s`).FindStringSubmatch(line); matches != nil {
				name = matches[2]
				interfaces[name] = make(map[string]interface{}, 0)
			}

			// ex.) link/ether 12:34:56:78:9a:bc brd ff:ff:ff:ff:ff:ff
			if matches := regexp.MustCompile(`link\/(\w+) ([\da-f\:]+) `).FindStringSubmatch(line); matches != nil {
				interfaces[name]["encap"] = g.TranslateEncap(matches[1])
				interfaces[name]["macAddress"] = matches[2]
			}

			// ex.) inet 10.0.4.7/24 brd 10.0.5.255 scope global eth0
			if matches := regexp.MustCompile(`inet (\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(\/(\d{1,2}))?`).FindStringSubmatch(line); matches != nil {
				interfaces[name]["ipAddress"] = matches[1]
				interfaces[name]["netmask"] = matches[3]
			}
		}
	}

	{
		// ip -f inet route show
		out, err := exec.Command("ip", "-f", "inet", "route", "show").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run ip command (skip this spec): %s", err)
			return interfaces, err
		}

		for _, line := range strings.Split(string(out), "\n") {
			// ex.) 10.0.3.0/24 dev eth0  proto kernel  scope link  src 10.0.4.7
			if matches := regexp.MustCompile(`^([^\s]+)\s(.*)$`).FindStringSubmatch(line); matches != nil {
				if matches := regexp.MustCompile(`\bdev\s+([^\s]+)\b`).FindStringSubmatch(matches[2]); matches != nil {
					name = matches[1]
				} else {
					continue
				}

				// ex.) 10.0.3.0/24
				if matches := regexp.MustCompile(`(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})(\/(\d{1,2}))?`).FindStringSubmatch(matches[1]); matches != nil {
					interfaces[name]["address"] = matches[1]
				}

				// ex.) default via 10.0.3.1 dev eth0
				if matches := regexp.MustCompile(`\bvia\s+([^\s]+)\b`).FindStringSubmatch(matches[2]); matches != nil {
					interfaces[name]["defaultGateway"] = matches[1]
				}
			}
		}
	}

	return interfaces, nil
}

func (g *InterfaceGenerator) GenerateByIfconfigCommand() (map[string]map[string]interface{}, error) {
	interfaces := make(map[string]map[string]interface{})
	name := ""

	{
		// ifconfig -a
		out, err := exec.Command("ifconfig", "-a").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run ifconfig command (skip this spec): %s", err)
			return nil, err
		}

		for _, line := range strings.Split(string(out), "\n") {
			// ex.) eth0      Link encap:Ethernet  HWaddr 12:34:56:78:9a:bc
			if matches := regexp.MustCompile(`^([0-9a-zA-Z@\.\:\-_]+)\s+`).FindStringSubmatch(line); matches != nil {
				name = matches[1]
				interfaces[name] = make(map[string]interface{}, 0)
			}
			// ex.) eth0      Link encap:Ethernet  HWaddr 12:34:56:78:9a:bc
			if matches := regexp.MustCompile(`Link encap:(Local Loopback)|Link encap:(.+?)\s`).FindStringSubmatch(line); matches != nil {
				if matches[1] != "" {
					interfaces[name]["encap"] = g.TranslateEncap(matches[1])
				} else {
					interfaces[name]["encap"] = g.TranslateEncap(matches[2])
				}
			}
			// ex.) eth0      Link encap:Ethernet  HWaddr 00:16:3e:4f:f3:41
			if matches := regexp.MustCompile(`HWaddr (.+?)\s`).FindStringSubmatch(line); matches != nil {
				interfaces[name]["macAddress"] = matches[1]
			}
			// ex.) inet addr:10.0.4.7  Bcast:10.0.5.255  Mask:255.255.255.0
			if matches := regexp.MustCompile(`inet addr:(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`).FindStringSubmatch(line); matches != nil {
				interfaces[name]["ipAddress"] = matches[1]
			}
			// ex.) inet addr:10.0.4.7  Bcast:10.0.5.255  Mask:255.255.255.0
			if matches := regexp.MustCompile(`Mask:(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`).FindStringSubmatch(line); matches != nil {
				netmask, _ := net.ParseIP(matches[1]).DefaultMask().Size()
				interfaces[name]["netmask"] = strconv.Itoa(netmask)
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

		routeRegexp := regexp.MustCompile(`^0\.0\.0\.0`)
		for _, line := range strings.Split(string(out), "\n") {
			// Destination     Gateway         Genmask         Flags Metric Ref    Use Iface
			// 0.0.0.0         10.0.3.1        0.0.0.0         UG    0      0        0 eth0
			if routeRegexp.FindStringSubmatch(line) != nil {
				routeResults := regexp.MustCompile(`[ \t]+`).Split(line, 8)
				if len(routeResults) < 8 || interfaces[routeResults[7]] == nil {
					continue
				}
				interfaces[routeResults[7]]["defaultGateway"] = routeResults[1]
			}
		}
	}

	{
		// arp -an
		out, err := exec.Command("arp", "-an").Output()
		if err != nil {
			interfaceLogger.Errorf("Failed to run arp command (skip this spec): %s", err)
			return interfaces, err
		}

		arpRegexp := regexp.MustCompile(`^\S+ \((\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\) at ([a-fA-F0-9\:]+) \[(\w+)\] on ([0-9a-zA-Z\.\:\-]+)`)
		for _, line := range strings.Split(string(out), "\n") {
			// ex.) ? (10.0.3.2) at 01:23:45:67:89:ab [ether] on eth0
			if matches := arpRegexp.FindStringSubmatch(line); matches != nil {
				if interfaces[matches[4]] == nil {
					continue
				}
				interfaces[matches[4]]["address"] = matches[1]
			}
		}
	}

	return interfaces, nil
}

func (g *InterfaceGenerator) TranslateEncap(encap string) string {
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
