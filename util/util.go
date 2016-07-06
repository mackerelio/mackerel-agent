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

// TimeoutDuration is option of `Runcommand()` set timeout limit of command execution.
var TimeoutDuration = 30 * time.Second

// TimeoutKillAfter is option of `RunCommand()` set waiting limit to `kill -kill` after terminating the command.
var TimeoutKillAfter = 10 * time.Second

// RunCommand runs command (in two string) and returns stdout, stderr strings and its exit code.
func RunCommand(command, user string) (string, string, int, error) {
	cmd := exec.Command("/bin/sh", "-c", command)
	if user != "" {
		cmd = exec.Command("sudo", "-u", user, "/bin/sh", "-c", command)
	}
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  TimeoutDuration,
		KillAfter: TimeoutKillAfter,
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
