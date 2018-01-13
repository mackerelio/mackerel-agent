// +build !windows

package cmdutil

import (
	"os/exec"
	"strings"
	"testing"
	"time"
)

const trapgo = "testdata/trapgo"

func init() {
	err := exec.Command("go", "build", "-o", trapgo, "testdata/trap.go").Run()
	if err != nil {
		panic(err)
	}
}

func TestRunCommandArgs_signalHandled(t *testing.T) {
	stdout, _, exitCode, err := RunCommandArgs([]string{trapgo}, CommandOption{
		TimeoutDuration: 50 * time.Millisecond,
	})
	if err != nil {
		t.Error("err should be nil but:", err)
	}
	expectOut := "signal received"
	if strings.TrimSpace(stdout) != expectOut {
		t.Errorf("stdout shoud be %q but: %s", expectOut, stdout)
	}
	if exitCode != 23 {
		t.Errorf("exitCode should be 23, but: %d", exitCode)
	}
}
