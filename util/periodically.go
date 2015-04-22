package util

import (
	"time"
)

// Periodically invokes function proc with specified interval. The precision is 1/100 of the interval.
func Periodically(proc func(), interval time.Duration, cancel <-chan struct{}) {
	checkInterval := interval / 100

	nextTime := time.Now().Add(interval)

	ticker := time.NewTicker(checkInterval)
	defer ticker.Stop()

	for {
		select {
		case t := <-ticker.C:
			if t.After(nextTime) {
				nextTime = nextTime.Add(interval)
				go proc()
			}

		case <-cancel:
			return
		}
	}
}
