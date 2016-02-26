package main

import (
	"testing"

	"github.com/motemen/go-cli"
)

func TestCommandRegisterd(t *testing.T) {
	_, ok := cli.Default.Commands[""]
	if !ok {
		t.Errorf("main command is not registerd. It may be `commands_gen.go` not generated")
	}
}
