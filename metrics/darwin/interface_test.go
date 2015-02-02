// +build darwin

package darwin

import "testing"

func TestInterfaceGenerator(t *testing.T) {
	g := &InterfaceGenerator{1}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	metrics := []string{
		"rxBytes", "txBytes",
	}

	for _, metric := range metrics {
		if value, ok := values["interface.en0."+metric+".delta"]; !ok {
			t.Errorf("Value for interface.en0.%s.delta should be collected", metric)
		} else {
			t.Logf("Interface en0 '%s' delta collected: %+v", metric, value)
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
