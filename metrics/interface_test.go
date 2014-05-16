package metrics

import (
	"os"
	"testing"
)

func TestInterfaceGenerator(t *testing.T) {
	if _, err := os.Stat("/etc/fedora-release"); err == nil {
		t.Skip("The OS seems to be Fedora. Skipping interface test for now")
	}

	g := &InterfaceGenerator{1}
	values, err := g.Generate()
	if err != nil {
		t.Error("should not raise error")
	}

	metrics := []string{
		"rxBytes", "rxPackets", "rxErrors", "rxDrops",
		"rxFifo", "rxFrame", "rxCompressed", "rxMulticast",
		"txBytes", "txPackets", "txErrors", "txDrops",
		"txFifo", "txColls", "txCarrier", "txCompressed",
	}

	for _, metric := range metrics {
		if value, ok := values["interface.eth0."+metric]; !ok {
			t.Errorf("Value for interface.eth0.%s should be collected", metric)
		} else {
			t.Logf("Interface eth0 '%s' collected: %+v", metric, value)
		}
	}

	for _, metric := range metrics {
		if value, ok := values["interface.eth0."+metric+".delta"]; !ok {
			t.Errorf("Value for interface.eth0.%s.delta should be collected", metric)
		} else {
			t.Logf("Interface eth0 '%s' delta collected: %+v", metric, value)
		}
	}

	for _, metric := range metrics {
		if _, ok := values["interface.lo."+metric]; ok {
			t.Errorf("Value for interface.lo.%s should NOT be collected", metric)
		} else {
			t.Logf("Interface lo '%s' NOT collected", metric)
		}
	}
}
