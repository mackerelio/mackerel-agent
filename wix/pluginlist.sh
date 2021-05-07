#!/bin/sh

d=$(git rev-parse --show-toplevel)
go list -f '{{range .Imports}}{{println .}}{{end}}' "$d/plugins_windows.go"
