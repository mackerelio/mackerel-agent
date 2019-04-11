package spec

import (
	"testing"

	mkr "github.com/mackerelio/mackerel-client-go"
)

func TestIsLoopback(t *testing.T) {
	t.Run("should be loopback", func(t *testing.T) {
		tests := []mkr.Interface{
			// single IPv4Address
			{
				IPv4Addresses: []string{"127.0.0.1"},
			},
			// single IPv6Address
			{
				IPv6Addresses: []string{"::1"},
			},
			// two IPv4Addresses
			{
				IPv4Addresses: []string{
					"127.0.0.1",
					"127.0.0.2",
				},
			},
			// two IPv6Addresses
			{
				IPv6Addresses: []string{
					"::1",
					"::0:1",
				},
			},
			// IPv4Address and IPv6Address
			{
				IPv4Addresses: []string{"127.0.0.1"},
				IPv6Addresses: []string{"::1"},
			},
		}
		for _, tt := range tests {
			if !IsLoopback(tt) {
				t.Errorf("IsLoopback(%+v) = false; want true", tt)
			}
		}
	})

	t.Run("should not be loopback", func(t *testing.T) {
		tests := []mkr.Interface{
			// empty
			{
				IPv4Addresses: nil,
			},
			// single IPAddress
			{
				IPAddress: "127.0.0.1",
			},
			// single IPv4Address
			{
				IPv4Addresses: []string{"227.0.0.1"},
			},
			// single IPv6Address
			{
				IPv6Addresses: []string{"::2"},
			},
			// two IPv4Addresses
			{
				IPv4Addresses: []string{
					"227.0.0.1",
					"127.0.0.2",
				},
			},
			// two IPv6Addresses
			{
				IPv6Addresses: []string{
					"::1",
					"::2",
				},
			},
			// IPAddress and IPv4Address
			{
				IPv4Addresses: []string{"128.0.0.1"},
				IPv6Addresses: []string{"::1"},
			},
			// IPv4Address and IPv6Address
			{
				IPv4Addresses: []string{"127.0.0.1"},
				IPv6Addresses: []string{"::2"},
			},
		}
		for _, tt := range tests {
			if IsLoopback(tt) {
				t.Errorf("IsLoopback(%+v) = true; want false", tt)
			}
		}
	})
}
