package spec

import (
	"net"

	mkr "github.com/mackerelio/mackerel-client-go"
)

// IsLoopback returns true if iface contains only loopback addresses.
// Is it possible that a interface contains mixed IPs both loopback address and else?
func IsLoopback(iface mkr.Interface) bool {
	// IPAddress field is defined by client library, but not used in agent.
	n4 := len(iface.IPv4Addresses)
	n6 := len(iface.IPv6Addresses)
	if n4+n6 == 0 {
		return false
	}
	addrs := make([]string, n4+n6)
	copy(addrs[:], iface.IPv4Addresses)
	copy(addrs[n4:], iface.IPv6Addresses)

	for _, addr := range addrs {
		ip := net.ParseIP(addr)
		if ip == nil {
			continue
		}
		if !ip.IsLoopback() {
			return false
		}
	}
	return true
}

// Interfaces are map of network interfaces per name
type Interfaces map[string]mkr.Interface

func (ifs Interfaces) getOrNew(name string) mkr.Interface {
	iface, ok := ifs[name]
	if ok {
		return iface
	}
	return mkr.Interface{Name: name}
}

// SetMacAddress sets the macaddress
func (ifs Interfaces) SetMacAddress(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.MacAddress = addr
	ifs[name] = iface
}

// AppendIPv4Address appends ipv4address
func (ifs Interfaces) AppendIPv4Address(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.IPv4Addresses = append(iface.IPv4Addresses, addr)
	ifs[name] = iface
}

// AppendIPv6Address appends ipv6address
func (ifs Interfaces) AppendIPv6Address(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.IPv6Addresses = append(iface.IPv6Addresses, addr)
	ifs[name] = iface
}

// InterfaceGenerator retrieve network informations
type InterfaceGenerator interface {
	Generate() ([]mkr.Interface, error)
}
