package cmdutil

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/Songmu/timeout"
	"github.com/mackerelio/golib/logging"
)

var logger = logging.GetLogger("cmdutil")

// defaultTimeoutDuration is the duration after which a command execution will be timeout.
// timeoutKillAfter is option of `RunCommand()` set waiting limit to `kill -kill` after
// terminating the command.
var (
	defaultTimeoutDuration = 30 * time.Second
	timeoutKillAfter       = 10 * time.Second
)

var cmdBase = []string{"sh", "-c"}

func init() {
	if runtime.GOOS == "windows" {
		cmdBase = []string{"cmd", "/c"}
	}
}

// CommandOption carries a timeout duration.
type CommandOption struct {
	User            string
	Env             []string
	TimeoutDuration time.Duration
}

// RunCommand runs command (in two string) and returns stdout, stderr strings and its exit code.
func RunCommand(command string, opt CommandOption) (stdout, stderr string, exitCode int, err error) {
	cmdArgs := append(cmdBase, command)
	return RunCommandArgs(cmdArgs, opt)
}

// RunCommandArgs run the command
func RunCommandArgs(cmdArgs []string, opt CommandOption) (stdout, stderr string, exitCode int, err error) {
	args := append([]string{}, cmdArgs...)
	if opt.User != "" {
		if runtime.GOOS == "windows" {
			logger.Warningf("RunCommand ignore option: user = %q", opt.User)
		} else {
			args = append([]string{"sudo", "-Eu", opt.User}, args...)
		}
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), opt.Env...)
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  defaultTimeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	if opt.TimeoutDuration != 0 {
		tio.Duration = opt.TimeoutDuration
	}
	exitStatus, stdout, stderr, err := tio.Run()

	if err == nil && exitStatus.IsTimedOut() && (runtime.GOOS == "windows" || exitStatus.Signaled) {
		err = fmt.Errorf("command timed out")
	}
	if err != nil {
		logger.Errorf("RunCommand error command: %v, error: %s", cmdArgs, err.Error())
	}
	return stdout, stderr, exitStatus.GetChildExitCode(), err
}
