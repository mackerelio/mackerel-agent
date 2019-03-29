package spec

import mkr "github.com/mackerelio/mackerel-client-go"

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
