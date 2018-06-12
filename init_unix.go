// +build !windows

package main

import "os"

func init() {
	if path := os.Getenv("PATH"); path != "" {
		os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+path)
	} else {
		os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin")
	}
	// prevent changing outputs of some command, e.g. ifconfig.
	os.Setenv("LANG", "C")
	os.Setenv("LC_ALL", "C")
}
