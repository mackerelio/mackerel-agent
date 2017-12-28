package util

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Songmu/timeout"
	"github.com/mackerelio/golib/logging"
)

var utilLogger = logging.GetLogger("util")

const (
	// defaultTimeoutDuration is the duration after which a command execution will be timeout.
	defaultTimeoutDuration = 30 * time.Second
)

// TimeoutKillAfter is option of `RunCommand()` set waiting limit to `kill -kill` after terminating the command.
var TimeoutKillAfter = 10 * time.Second

var cmdBase = []string{"sh", "-c"}

// CommandContext carries a timeout duration.
type CommandContext struct {
	TimeoutDuration time.Duration
}

func init() {
	if runtime.GOOS == "windows" {
		cmdBase = []string{"cmd", "/c"}
	}
}

// RunCommand runs command (in two string) and returns stdout, stderr strings and its exit code.
func RunCommand(ctx CommandContext, command, user string, env []string) (stdout, stderr string, exitCode int, err error) {
	cmdArgs := append(cmdBase, command)
	return RunCommandArgs(ctx, cmdArgs, user, env)
}

// RunCommandArgs run the command
func RunCommandArgs(ctx CommandContext, cmdArgs []string, user string, env []string) (stdout, stderr string, exitCode int, err error) {
	args := append([]string{}, cmdArgs...)
	if user != "" {
		if runtime.GOOS == "windows" {
			utilLogger.Warningf("RunCommand ignore option: user = %q", user)
		} else {
			args = append([]string{"sudo", "-Eu", user}, args...)
		}
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), env...)
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  defaultTimeoutDuration,
		KillAfter: TimeoutKillAfter,
	}
	if ctx.TimeoutDuration != 0 {
		tio.Duration = ctx.TimeoutDuration
	}
	exitStatus, stdout, stderr, err := tio.Run()

	if err == nil && exitStatus.IsTimedOut() {
		err = fmt.Errorf("command timed out")
	}
	if err != nil {
		utilLogger.Errorf("RunCommand error command: %v, error: %s", cmdArgs, err.Error())
	}
	return stdout, stderr, exitStatus.GetChildExitCode(), err
}
