//go:build linux

package linux

import (
	"github.com/vishvananda/netlink"

	"github.com/mackerelio/mackerel-agent/spec"
	mkr "github.com/mackerelio/mackerel-client-go"
)

// InterfaceGenerator XXX
type InterfaceGenerator struct {
}

// Generate XXX
func (g *InterfaceGenerator) Generate() ([]mkr.Interface, error) {
	var interfaces spec.Interfaces
	interfaces, err := g.generateByNetLink()
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

func (g *InterfaceGenerator) generateByNetLink() (spec.Interfaces, error) {
	interfaces := make(spec.Interfaces)

	links, err := netlink.LinkList()
	if err != nil {
		return nil, err
	}

	for _, link := range links {
		attr := link.Attrs()
		name := attr.Name

		interfaces.SetMacAddress(name, attr.HardwareAddr.String())

		addrs, err := netlink.AddrList(link, netlink.FAMILY_V4)
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			interfaces.AppendIPv4Address(name, addr.IP.String())
		}

		addrs, err = netlink.AddrList(link, netlink.FAMILY_V6)
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			interfaces.AppendIPv6Address(name, addr.IP.String())
		}
	}
	return interfaces, nil
}
