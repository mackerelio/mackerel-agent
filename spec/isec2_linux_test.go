package spec

import (
	"testing"
)

func TestIsEC2UUID(t *testing.T) {
	tests := []struct {
		uuid   string
		expect bool
	}{
		{"ec2e1916-9099-7caf-fd21-01234abcdef", true},
		{"EC2E1916-9099-7CAF-FD21-01234ABCDEF", true},
		{"45e12aec-dcd1-b213-94ed-01234abcdef", true}, // litte endian
		{"45E12AEC-DCD1-B213-94ED-01234ABCDEF", true}, // litte endian
		{"abcd1916-9099-7caf-fd21-01234abcdef", false},
		{"ABCD1916-9099-7CAF-FD21-01234ABCDEF", false},
		{"", false},
	}

	for _, tc := range tests {
		if isEC2UUID(tc.uuid) != tc.expect {
			t.Errorf("isEC2() should be %t: %q\n", tc.expect, tc.uuid)
		}
	}
}
