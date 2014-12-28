// +build windows

package windows

import (
	"net"

	"github.com/mackerelio/mackerel-agent/logging"
	. "github.com/mackerelio/mackerel-agent/util/windows"
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
	results := make([]map[string]interface{}, 0)

	ifs, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	ai, err := GetAdapterList()
	if err != nil {
		return nil, err
	}

	for _, ifi := range ifs {
		addr, err := ifi.Addrs()
		if err != nil {
			return nil, err
		}
		name := ifi.Name
		for ; ai != nil; ai = ai.Next {
			if ifi.Index == int(ai.Index) {
				name = BytePtrToString(&ai.Description[0])
			}
		}

		results = append(results, map[string]interface{}{
			"name":       name,
			"ipAddress":  addr[0].String(),
			"macAddress": ifi.HardwareAddr.String(),
		})
	}

	return results, nil
}
