// +build !windows

package spec

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
)

func TestFilesystemGenerate(t *testing.T) {
	g := &FilesystemGenerator{}

	result, err := g.Generate()
	if err != nil {
		t.Skipf("Generate() failed: %s", err)
	}

	_, resultTypeOk := result.(mackerel.FileSystem)
	if !resultTypeOk {
		t.Errorf("Return type of Generate() shuold be mackerel.FileSystem")
	}
}
