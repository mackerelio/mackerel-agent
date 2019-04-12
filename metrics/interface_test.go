// +build !windows

package metrics

import (
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestInterfaceGenerator(t *testing.T) {
	g := &InterfaceGenerator{1 * time.Second}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("error should be nil but got: %s", err)
	}

	metrics := []string{"rxBytes", "txBytes"}

	name := lookupDefaultName(values, "eth0")
	if runtime.GOOS != "linux" {
		name = "en0"
	}
	for _, metric := range metrics {
		metricName := "interface." + name + "." + metric + ".delta"
		if _, ok := values[metricName]; !ok {
			t.Errorf("Value for %s should be collected", metricName)
		}
	}

	name = "lo"
	if runtime.GOOS != "linux" {
		name = "lo0"
	}
	for _, metric := range metrics {
		metricName := "interface." + name + "." + metric + ".delta"
		if _, ok := values[metricName]; ok {
			t.Errorf("Value for %s should NOT be collected", metricName)
		}
	}

	if runtime.GOOS == "linux" {
		for k := range values {
			if strings.HasPrefix(k, "interface.veth") {
				t.Errorf("Value for %s should NOT be collected", k)
			}
		}
	}

	t.Logf("interface metrics: %+v", values)
}

// lookupDefaultName returns network interface name that seems to be default NIC.
//
// There is type-differed version at spec/linux/interface_test.go.
func lookupDefaultName(values Values, fallback string) string {
	if runtime.GOOS != "linux" {
		return fallback
	}
	trim := func(s string) string {
		return strings.Split(s, ".")[1]
	}
	for key := range values {
		switch {
		case strings.HasPrefix(key, "interface.eth"):
			return trim(key)
		case strings.HasPrefix(key, "interface.en"):
			return trim(key)
		case strings.HasPrefix(key, "interface.wl"):
			return trim(key)
		}
	}
	return fallback
}
