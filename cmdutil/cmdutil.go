package cmdutil

import (
	"bytes"
	"context"
	"errors"
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
	return RunCommandContext(context.Background(), command, opt)
}

// RunCommandContext runs command with context
func RunCommandContext(ctx context.Context, command string, opt CommandOption) (stdout, stderr string, exitCode int, err error) {
	cmdArgs := append(cmdBase, command)
	return RunCommandArgsContext(ctx, cmdArgs, opt)
}

var errTimedOut = errors.New("command timed out")

// RunCommandArgs run the command
func RunCommandArgs(cmdArgs []string, opt CommandOption) (stdout, stderr string, exitCode int, err error) {
	return RunCommandArgsContext(context.Background(), cmdArgs, opt)
}

// RunCommandArgsContext runs command by args with context
func RunCommandArgsContext(ctx context.Context, cmdArgs []string, opt CommandOption) (stdout, stderr string, exitCode int, err error) {
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
	outbuf := &bytes.Buffer{}
	errbuf := &bytes.Buffer{}
	cmd.Stdout = outbuf
	cmd.Stderr = errbuf
	tio := &timeout.Timeout{
		Cmd:       cmd,
		Duration:  defaultTimeoutDuration,
		KillAfter: timeoutKillAfter,
	}
	if opt.TimeoutDuration != 0 {
		tio.Duration = opt.TimeoutDuration
	}
	exitStatus, err := tio.RunContext(ctx)
	stdout = outbuf.String()
	stderr = errbuf.String()
	exitCode = -1
	if err == nil && exitStatus.IsTimedOut() && (runtime.GOOS == "windows" || exitStatus.Signaled) {
		err = errTimedOut
		exitCode = exitStatus.GetChildExitCode()
	}
	if err != nil {
		logger.Errorf("RunCommand error. command: %v, error: %s", cmdArgs, err.Error())
		if terr, ok := err.(*timeout.Error); ok {
			exitCode = terr.ExitCode
		}
		return stdout, stderr, exitCode, err
	}
	return stdout, stderr, exitStatus.GetChildExitCode(), nil
}
