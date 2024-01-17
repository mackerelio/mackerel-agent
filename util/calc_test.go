package util

import (
	"fmt"
	"testing"
)

func TestDiffResettableCounter(t *testing.T) {
	cases := []struct {
		inCurrent  uint64
		inPrevious uint64
		want       uint64
	}{
		{100, 30, 70},
		{20, 50, 20}, // counter has been reset
	}
	for _, tt := range cases {
		t.Run(fmt.Sprintf("%d - %d", tt.inCurrent, tt.inPrevious), func(t *testing.T) {
			got := DiffResettableCounter(tt.inCurrent, tt.inPrevious)
			if got != tt.want {
				t.Errorf("want=%d, got=%d", tt.want, got)
			}
		})
	}
}
