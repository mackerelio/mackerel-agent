//go:build !windows
// +build !windows

package main

import "os"

const MkrPluginInstallPath = "/opt/mackerel-agent/plugins/bin"

func init() {
	if path := os.Getenv("PATH"); path != "" {
		os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+path+":"+MkrPluginInstallPath)
	} else {
		os.Setenv("PATH", "/sbin:/usr/sbin:/bin:/usr/bin:"+MkrPluginInstallPath)
	}
	// prevent changing outputs of some command, e.g. ifconfig.
	os.Setenv("LANG", "C")
	os.Setenv("LC_ALL", "C")
}
