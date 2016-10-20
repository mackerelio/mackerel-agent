package util

import "testing"

func TestSanitizeMetricKey(t *testing.T) {
	input := "abc*defあ.ggg"
	expect := "abc_def__ggg"
	output := SanitizeMetricKey(input)
	if output != expect {
		t.Errorf("invalid output of `SanitizeMetricKey`. expected: %s, output: %s", expect, output)
	}
}
