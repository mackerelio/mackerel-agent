// +build darwin

package darwin

import (
	"reflect"
	"testing"
)

func TestFilesystemGenerator(t *testing.T) {
	g := &FilesystemGenerator{}

	if g.Key() != "filesystem" {
		t.Error("key should be 'filesystem'")
	}
}

func TestFilesystemGenerate(t *testing.T) {
	g := &FilesystemGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Skipf("Generate() failed: %s", err)
	}

	filesystems, resultTypeOk := result.(map[string]map[string]interface{})
	if !resultTypeOk {
		t.Errorf("Return type of Generate() shuold be map[string]map[string]interface{}")
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
