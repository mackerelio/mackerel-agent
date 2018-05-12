// +build windows

package main

import (
	"os"
	"path/filepath"

	"gopkg.in/pathlist.v0"
	"gopkg.in/pathlist.v0/env"
)

func init() {
	if p, err := os.Executable(); err == nil {
		env.SetPath(pathlist.Must(pathlist.PrependTo(env.Path(), filepath.Dir(p))))
	}
}
