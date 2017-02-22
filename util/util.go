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
func RunCommand(command, user string) (stdout, stderr string, exitCode int, err error) {
	cmdArgs := []string{"/bin/sh", "-c", command}
	return RunCommandArgs(cmdArgs, user)
}

// RunCommandArgs run the command
func RunCommandArgs(cmdArgs []string, user string) (stdout, stderr string, exitCode int, err error) {
	args := append([]string{}, cmdArgs...)
	if user != "" {
		args = append([]string{"sudo", "-u", user}, args...)
	}
	cmd := exec.Command(args[0], args[1:]...)
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
		utilLogger.Errorf("RunCommand error command: %T, error: %s", cmdArgs, err.Error())
	}
	return stdout, stderr, exitStatus.GetChildExitCode(), err
}
