// +build !windows

package main

import (
	"os"

	"gopkg.in/pathlist.v0"
	"gopkg.in/pathlist.v0/env"
)

func init() {
	env.SetPath(pathlist.Must(pathlist.PrependTo(env.Path(),
		"/sbin", "/usr/sbin", "/bin", "/usr/bin")))
	// prevent changing outputs of some command, e.g. ifconfig.
	os.Setenv("LANG", "C")
	os.Setenv("LC_ALL", "C")
}
