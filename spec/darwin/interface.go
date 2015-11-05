// +build darwin

package darwin

import (
	"bufio"
	"bytes"
	"os/exec"
	"regexp"

	"github.com/mackerelio/mackerel-agent/logging"
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
func (g *InterfaceGenerator) Generate() (interface{}, error) {
	var interfaces map[string]map[string]interface{}

	interfaces, err := g.generateByIfconfigCommand()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for key, iface := range interfaces {
		if iface["loopback"] != nil {
			continue
		}
		if iface["ipAddress"] == nil && iface["ipv6Address"] == nil {
			continue
		}
		iface["name"] = key
		results = append(results, iface)
	}

	return results, nil
}

func (g *InterfaceGenerator) generateByIfconfigCommand() (map[string]map[string]interface{}, error) {
	interfaces := make(map[string]map[string]interface{})

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
			// ex.) en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
			if matches := regexp.MustCompile(`^([0-9a-zA-Z@\-_]+):\s+`).FindStringSubmatch(line); matches != nil {
				name = matches[1]
				interfaces[name] = make(map[string]interface{}, 0)
			}
			// ex.) lo0: flags=8049<UP,LOOPBACK,RUNNING,MULTICAST> mtu 16384
			if regexp.MustCompile(`LOOPBACK`).MatchString(line) {
				interfaces[name]["loopback"] = true
			}
			// ex.) ether 10:93:00:00:00:00
			if matches := regexp.MustCompile(`^\s*ether\s+([0-9a-f:]+)`).FindStringSubmatch(line); matches != nil {
				interfaces[name]["macAddress"] = matches[1]
			}
			// ex.) inet 10.0.3.1 netmask 0xffffff00 broadcast 10.0.3.255
			if matches := regexp.MustCompile(`^\s*inet\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})\s+netmask\s+(0x[0-9a-f]+)`).FindStringSubmatch(line); matches != nil {
				interfaces[name]["ipAddress"] = matches[1]
				interfaces[name]["netmask"] = matches[2]
			}
			// ex.) inet6 2001:268:c044:1111:1111:1111:1111:1111 prefixlen 64 autoconf
			if matches := regexp.MustCompile(`^\s*inet6\s+([0-9a-f:]+)\s+prefixlen\s+(\d+)`).FindStringSubmatch(line); matches != nil {
				interfaces[name]["ipv6Address"] = matches[1]
				interfaces[name]["v6netmask"] = matches[2]
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
		routeRegexp := regexp.MustCompile(`^default\s+(\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})`)
		lineScanner := bufio.NewScanner(bytes.NewReader(out))
		// ex.)
		// Routing tables

		// Internet:
		// Destination        Gateway            Flags        Refs      Use   Netif Expire
		// default            10.0.3.1      UGSc           25        0     en0
		for lineScanner.Scan() {
			line := lineScanner.Text()
			if routeRegexp.FindStringSubmatch(line) != nil {
				routeResults := regexp.MustCompile(`[ \t]+`).Split(line, 6)
				if len(routeResults) < 6 || interfaces[routeResults[5]] == nil {
					continue
				}
				interfaces[routeResults[5]]["defaultGateway"] = routeResults[1]
				break
			}
		}
	}

	return interfaces, nil
}
