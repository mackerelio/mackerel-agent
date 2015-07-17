// +build linux

package linux

import "testing"

func TestLoadAvg5Generate(t *testing.T) {
	_, err := (&Loadavg5Generator{}).Generate()

	if err != nil {
		t.Errorf("something went wrong")
	}
}
