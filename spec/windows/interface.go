// +build windows

package windows

import (
	"net"

	"github.com/mackerelio/golib/logging"
	"github.com/mackerelio/mackerel-agent/util/windows"
	mkr "github.com/mackerelio/mackerel-client-go"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

var interfaceLogger = logging.GetLogger("spec.interface")

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]mkr.Interface, error) {
	var results []mkr.Interface

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
		if name == "" {
			continue
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
			case *net.IPNet:
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
			results = append(results, mkr.Interface{
				Name:          name,
				IPv4Addresses: ipv4Addresses,
				IPv6Addresses: ipv6Addresses,
				MacAddress:    ifi.HardwareAddr.String(),
			})
		}
	}
	return results, nil
}
