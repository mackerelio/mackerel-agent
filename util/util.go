// +build linux darwin freebsd

package util

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/Songmu/timeout"
)

// timeoutDuration is option of `Runcommand()` set timeout limit of command execution.
var timeoutDuration = 30 * time.Second

// timeoutKillAfter is option of `RunCommand()` set waiting limit to `kill -kill` after terminating the command.
var timeoutKillAfter = 10 * time.Second

// RunCommand runs command (in one string) and returns stdout, stderr strings.
func RunCommand(command string) (string, string, error) {
	tio := &timeout.Timeout{
		Cmd:       exec.Command("/bin/sh", "-c", command),
		Duration:  timeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	exitStatus, stdout, stderr, err := tio.Run()

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("command timed out")
	}
	return stdout, stderr, err
}
