package util

import (
	"testing"
	"time"
)

func TestPeriodically(t *testing.T) {
	quit := make(chan struct{})
	counter := 0
	go Periodically(
		func() {
			counter++
		},
		400*time.Millisecond,
		quit,
	)
	time.Sleep(time.Second)
	quit <- struct{}{}

	if counter != 2 {
		t.Errorf("counter should be 2, but", counter)
	}
}
