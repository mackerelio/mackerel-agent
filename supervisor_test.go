// +build !windows

package main

import (
	"os"
	"os/exec"
	"testing"
)

func init() {
	err := exec.Command("go", "build", "-o", "testdata/stub-agent", "testdata/stub-agent.go").Run()
	if err != nil {
		panic(err)
	}
}

func TestSupervise(t *testing.T) {
	sv := &supervisor{
		prog: "testdata/stub-agent",
	}
	sv.start()
	sv.stop(os.Interrupt)
	err := sv.wait()

	if err == nil {
		t.Errorf("something went wrong")
	}
}
