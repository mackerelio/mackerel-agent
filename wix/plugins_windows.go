//go:build plugins && windows

package main

import (
	_ "github.com/mackerelio/go-check-plugins/check-disk"
	_ "github.com/mackerelio/go-check-plugins/check-file-age"
	_ "github.com/mackerelio/go-check-plugins/check-log"
	_ "github.com/mackerelio/go-check-plugins/check-ntservice"
	_ "github.com/mackerelio/go-check-plugins/check-ping"
	_ "github.com/mackerelio/go-check-plugins/check-procs"
	_ "github.com/mackerelio/go-check-plugins/check-tcp"
	_ "github.com/mackerelio/go-check-plugins/check-uptime"
	_ "github.com/mackerelio/go-check-plugins/check-windows-eventlog"
	_ "github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-mssql"
	_ "github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-snmp"
	_ "github.com/mackerelio/mackerel-agent-plugins/mackerel-plugin-windows-server-sessions"
	_ "github.com/mackerelio/mkr"
)
