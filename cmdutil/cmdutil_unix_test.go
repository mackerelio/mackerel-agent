// +build !windows

package cmdutil

import (
	"strings"
	"testing"
	"time"
)

func TestRunCommandArgs_signalHandled(t *testing.T) {
	cmd := []string{stubcmd, "-trap=SIGTERM", "-trap-exit=23", "-sleep=10"}
	stdout, _, exitCode, err := RunCommandArgs(cmd, CommandOption{
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
