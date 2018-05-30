// +build windows

package main

import (
	"os"
	"path/filepath"
	"strings"
)

func init() {
	if p, err := os.Executable(); err == nil {
		dir := filepath.Dir(p)
		if strings.Index(dir, ";") > -1 {
			dir = `"` + dir + `"`
		}
		os.Setenv("PATH", dir+string(filepath.ListSeparator)+os.Getenv("PATH"))
	}
}
