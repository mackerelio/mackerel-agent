// +build linux darwin freebsd netbsd

package util

import (
	"fmt"
	"os/exec"
	"time"

	"github.com/Songmu/timeout"
	"github.com/mackerelio/mackerel-agent/logging"
)

var utilLogger = logging.GetLogger("util")

// timeoutDuration is option of `Runcommand()` set timeout limit of command execution.
var timeoutDuration = 30 * time.Second

// timeoutKillAfter is option of `RunCommand()` set waiting limit to `kill -kill` after terminating the command.
var timeoutKillAfter = 10 * time.Second

// RunCommand runs command (in one string) and returns stdout, stderr strings and its exit code.
func RunCommand(command string) (string, string, int, error) {
	tio := &timeout.Timeout{
		Cmd:       exec.Command("/bin/sh", "-c", command),
		Duration:  timeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	exitStatus, stdout, stderr, err := tio.Run()

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("command timed out")
	}
	if err != nil {
		utilLogger.Errorf("RunCommand error command: %s, error: %s", command, err)
	}
	return stdout, stderr, exitStatus.GetChildExitCode(), err
}
