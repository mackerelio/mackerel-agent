// +build darwin

package darwin

import "testing"

func TestSwapGenerator(t *testing.T) {
	g := &SwapGenerator{}
	values, err := g.Generate()
	if err != nil {
		t.Errorf("should not raise error: %v", err)
	}

	for _, name := range []string{
		"swap_total",
		"swap_free",
	} {
		if v, ok := values["memory."+name]; !ok {
			t.Errorf("memory should has %s", name)
		} else {
			t.Logf("memory '%s' collected: %+v", name, v)
		}
	}
}
