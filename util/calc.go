package util

// Calculate incremental difference of uint64 counters that can be reset.
// Avoiding overflow, return current value if it is smaller than previous one.
func DiffResettableCounter(current, previous uint64) uint64 {
	if current < previous {
		// counter has been reset
		return current
	}
	return current - previous
}
