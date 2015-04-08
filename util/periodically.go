package util

import (
	"time"
)

func Periodically(proc func(), interval time.Duration, cancel chan struct{}) {
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
