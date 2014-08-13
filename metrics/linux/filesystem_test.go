// +build linux

package linux

import (
	"reflect"
	"testing"
)

func TestFilesystemGenerate(t *testing.T) {
	g := &FilesystemGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Skipf("Generate() failed: %s", err)
	}

	if _, ok := result["filesystem.sda1.size"]; !ok {
		t.Errorf("filesystem should has sda1.size")
	}

	if _, ok := result["filesystem.sda1.used"]; !ok {
		t.Errorf("filesystem should has sda1.size")
	}
}

func TestFilesystemCollectValues(t *testing.T) {
	g := &FilesystemGenerator{}

	filesystems, err := g.collectValues()
	if err != nil {
		t.Skipf("collectValues() failed: %s", err)
	}

	// tmpfs may be exists
	tmpfs, hasTmpfsEntry := filesystems["tmpfs"]

	if hasTmpfsEntry {
		for _, spec := range dfColumnSpecs {
			value, hasColumn := tmpfs[spec.name]

			if hasColumn {
				t.Logf("Value '%s' collected: %#v", spec.name, value)

				valueType := reflect.TypeOf(value).Name()
				var expectedType string
				if spec.isInt {
					expectedType = "int64"
				} else {
					expectedType = "string"
				}

				if valueType != expectedType {
					t.Errorf("Type mismatch of value '%s': expected %s but got %s", spec.name, expectedType, valueType)
				}
			} else {
				t.Errorf("Value '%s' should be collected", spec.name)
			}
		}
	} else {
		t.Log("Could not detect filesystem tmpfs")
	}
}
