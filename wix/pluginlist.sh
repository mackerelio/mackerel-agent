#!/bin/sh

d=$(dirname "$0")
go list -f '{{range .Imports}}{{println .}}{{end}}' "$d/plugins_windows.go"
