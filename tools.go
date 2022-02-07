//go:build tools
// +build tools

package main

import (
	_ "github.com/Songmu/gocredits/cmd/gocredits"
	_ "github.com/Songmu/goxz/cmd/goxz"
	_ "github.com/mattn/goveralls"
	_ "golang.org/x/lint/golint"
)
