// +build !windows

package spec

import "testing"

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

	_, resultTypeOk := result.(map[string]map[string]interface{})
	if !resultTypeOk {
		t.Errorf("Return type of Generate() shuold be map[string]map[string]interface{}")
	}
}
