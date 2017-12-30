// +build windows

package main

import (
	"os"
	"path/filepath"
)

func init() {
	if p, err := os.Executable(); err == nil {
		os.Setenv("PATH", filepath.Dir(p)+
			string(filepath.ListSeparator)+os.Getenv("PATH"))
	}
}
