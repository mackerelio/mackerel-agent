package config

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

// duration represents a non-negative time duration.
type duration int32 // in minutes

func (m *duration) UnmarshalText(text []byte) error {
	i, err := strconv.ParseInt(string(text), 10, 32)
	if err == nil {
		if i < 0 {
			return fmt.Errorf("duration out of range: %d", i)
		}
		*m = duration(i)
		return nil
	}
	if dur, err2 := time.ParseDuration(string(text)); err2 == nil {
		minutes := dur.Minutes()
		if minutes < 0 || float64(math.MaxInt32) < minutes {
			return fmt.Errorf("duration out of range: %v", dur)
		}
		if dur != dur.Round(time.Minute) {
			return fmt.Errorf("duration not multiple of 1m: %v", dur)
		}
		*m = duration(minutes)
		return nil
	}
	return err
}

func (m *duration) Minutes() *int32 {
	if m == nil {
		return nil
	}
	i := int32(*m)
	return &i
}
