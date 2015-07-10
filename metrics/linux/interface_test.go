// +build linux

package linux

import (
	"os"
	"testing"
	"time"
)

func TestInterfaceGenerator(t *testing.T) {
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		t.Skip("The OS seems to be Fedora. Skipping interface test for now")
	}

	if os.Getenv("TRAVIS") != "" {
		t.Skip("Skip: in Travis, Skipping interface test for now")
	}

	g := &InterfaceGenerator{1 * time.Second}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	metrics := []string{
		"rxBytes", "txBytes",
	}

	for _, metric := range metrics {
		if value, ok := values["interface.eth0."+metric+".delta"]; !ok {
			t.Errorf("Value for interface.eth0.%s.delta should be collected", metric)
		} else {
			t.Logf("Interface eth0 '%s' delta collected: %+v", metric, value)
		}
	}

	for _, metric := range metrics {
		if _, ok := values["interface.lo."+metric+".delta"]; ok {
			t.Errorf("Value for interface.lo.%s should NOT be collected", metric)
		} else {
			t.Logf("Interface lo '%s' NOT collected", metric)
		}
	}
}
