// +build windows

package windows

import (
	"net"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/spec"
	"github.com/mackerelio/mackerel-agent/util/windows"
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
	var results []spec.NetInterface

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ai, err := windows.GetAdapterList()
	if err != nil {
		return nil, err
	}

	for _, ifi := range ifs {
		if ifi.Flags&net.FlagLoopback != 0 {
			continue
		}

		// XXX occur mojibake when containing multi-byte strings
		name := ifi.Name
		for ; ai != nil; ai = ai.Next {
			if ifi.Index == int(ai.Index) {
				name = windows.BytePtrToString(&ai.Description[0])
			}
		}

		addrs, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		ipv4Addresses := []string{}
		ipv6Addresses := []string{}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ipv4 := ip.To4(); ipv4 != nil {
				ipv4Addresses = append(ipv4Addresses, ipv4.String())
				continue
			}
			if ipv6 := ip.To16(); ipv6 != nil {
				ipv6Addresses = append(ipv6Addresses, ipv6.String())
			}
		}

		if len(ipv4Addresses) > 0 || len(ipv6Addresses) > 0 {
			results = append(results, spec.NetInterface{
				Name:          name,
				IPv4Addresses: ipv4Addresses,
				IPv6Addresses: ipv6Addresses,
				MacAddress:    ifi.HardwareAddr.String(),
			})
		}
	}
	return results, nil
}
