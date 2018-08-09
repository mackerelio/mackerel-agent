package spec

// NetInterface represents network interface informations
type NetInterface struct {
	Name           string   `json:"name"`
	Encap          string   `json:"encap,omitempty"`
	IPv4Addresses  []string `json:"ipv4Addresses"`
	IPv6Addresses  []string `json:"ipv6Addresses"`
	Address        string   `json:"address,omitempty"`
	V6Address      string   `json:"v6address,omitempty"`
	MacAddress     string   `json:"macAddress,omitempty"`
	DefaultGateway string   `json:"defaultGateway,omitempty"`
}

// NetInterfaces are map of network interfaces per name
type NetInterfaces map[string]NetInterface

func (ifs NetInterfaces) getOrNew(name string) NetInterface {
	iface, ok := ifs[name]
	if ok {
		return iface
	}
	return NetInterface{Name: name}
}

// SetEncap sets the encap
func (ifs NetInterfaces) SetEncap(name, encap string) {
	iface := ifs.getOrNew(name)
	iface.Encap = encap
	ifs[name] = iface
}

// SetAddress sets the address
func (ifs NetInterfaces) SetAddress(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.Address = addr
	ifs[name] = iface
}

// SetV6Address sets the v6address
func (ifs NetInterfaces) SetV6Address(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.V6Address = addr
	ifs[name] = iface
}

// SetMacAddress sets the macaddress
func (ifs NetInterfaces) SetMacAddress(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.MacAddress = addr
	ifs[name] = iface
}

// SetDefaultGateway sets the defaultGateway
func (ifs NetInterfaces) SetDefaultGateway(name, gateway string) {
	iface := ifs.getOrNew(name)
	iface.DefaultGateway = gateway
	ifs[name] = iface
}

// AppendIPv4Address appends ipv4address
func (ifs NetInterfaces) AppendIPv4Address(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.IPv4Addresses = append(iface.IPv4Addresses, addr)
	ifs[name] = iface
}

// AppendIPv6Address appends ipv6address
func (ifs NetInterfaces) AppendIPv6Address(name, addr string) {
	iface := ifs.getOrNew(name)
	iface.IPv6Addresses = append(iface.IPv6Addresses, addr)
	ifs[name] = iface
}

// InterfaceGenerator retrieve network informations
type InterfaceGenerator interface {
	Generate() ([]NetInterface, error)
}
