// +build darwin

package darwin

import (
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

	for _, metric := range metrics {
		if _, ok := values["interface.en0."+metric+".delta"]; !ok {
			t.Errorf("Value for interface.en0.%s.delta should be collected", metric)
		}
	}

	for _, metric := range metrics {
		if _, ok := values["interface.lo0."+metric+".delta"]; ok {
			t.Errorf("Value for interface.lo0.%s should NOT be collected", metric)
		}
	}

	t.Logf("interface metrics: %+v", values)
}
