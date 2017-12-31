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

	name := "eth0"
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
