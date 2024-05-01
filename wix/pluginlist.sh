#!/bin/sh

set -eu

d=$(dirname "$0")

# `main` package can't import, but it only manages versions in go.mod, so I use the `-e` option to ignore errors and list them.
go list -e -f '{{range .Imports}}{{println .}}{{end}}' "$d/plugins_windows.go"
