//go:build windows

package windows

import (
	"testing"
)

func TestProcessorQueueLengthGenerator(t *testing.T) {
	g, _ := NewProcessorQueueLengthGenerator()

	_, err := g.Generate()
	if err != nil {
		t.Errorf("Generate() failed: %s", err)
	}
}
