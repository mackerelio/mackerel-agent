// +build linux

package linux

import (
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
