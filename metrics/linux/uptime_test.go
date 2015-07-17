// +build linux

package linux

import "testing"

func TestUptime(t *testing.T) {
	_, err := (&UptimeGenerator{}).Generate()

	if err != nil {
		t.Errorf("something went wrong")
	}
}
