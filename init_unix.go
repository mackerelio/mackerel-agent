// +build !windows

package main

import "os"

func init() {
	os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+os.Getenv("PATH"))
	// prevent changing outputs of some command, e.g. ifconfig.
	os.Setenv("LANG", "C")
	os.Setenv("LC_ALL", "C")
}
