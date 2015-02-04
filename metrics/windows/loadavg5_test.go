// +build windows

package windows

import (
	"testing"
)

func TestLoadavg5Generator(t *testing.T) {
	g, _ := NewLoadavg5Generator()

	_, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() failed: %s", err)
	}
}
