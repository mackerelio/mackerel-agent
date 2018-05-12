// +build !windows

package main

import (
	"os"

	"gopkg.in/pathlist.v0"
	"gopkg.in/pathlist.v0/env"
)

func init() {
	prependTo := func(list pathlist.List, dir string) pathlist.List {
		return pathlist.Must(pathlist.PrependTo(list, dir))
	}
	list := prependTo(env.Path(), "/usr/bin")
	list = prependTo(list, "/bin")
	list = prependTo(list, "/usr/sbin")
	list = prependTo(list, "/sbin")
	env.SetPath(list)
	// prevent changing outputs of some command, e.g. ifconfig.
	os.Setenv("LANG", "C")
	os.Setenv("LC_ALL", "C")
}
