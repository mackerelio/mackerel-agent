// +build linux darwin freebsd netbsd

package util

import (
	"reflect"
	"testing"
)

func TestCollectDfValues(t *testing.T) {
	dfColumnSpecs := []DfColumnSpec{
		{"kb_size", true},
		{"kb_used", true},
		{"kb_available", true},
		{"percent_used", false},
		{"mount", false},
	}

	filesystems, err := CollectDfValues(dfColumnSpecs)
	if err != nil {
		t.Skipf("collectValues() failed: %s", err)
	}

	// tmpfs may be exists
	tmpfs, hasTmpfsEntry := filesystems["tmpfs"]

	if hasTmpfsEntry {
		for _, spec := range dfColumnSpecs {
			value, hasColumn := tmpfs[spec.Name]

			if hasColumn {
				t.Logf("Value '%s' collected: %#v", spec.Name, value)

				valueType := reflect.TypeOf(value).Name()
				var expectedType string
				if spec.IsInt {
					expectedType = "int64"
				} else {
					expectedType = "string"
				}

				if valueType != expectedType {
					t.Errorf("Type mismatch of value '%s': expected %s but got %s", spec.Name, expectedType, valueType)
				}
			} else {
				t.Errorf("Value '%s' should be collected", spec.Name)
			}
		}
	} else {
		t.Log("Could not detect filesystem tmpfs")
	}
}
