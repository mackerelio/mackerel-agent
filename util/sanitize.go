package util

import "regexp"

var sanitizerReg = regexp.MustCompile(`[^A-Za-z0-9_-]`)

// SanitizeMetricKey sanitize metric keys to be Mackerel friendly
func SanitizeMetricKey(key string) string {
	return sanitizerReg.ReplaceAllString(key, "_")
}
