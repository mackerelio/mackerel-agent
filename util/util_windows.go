package util

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
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
	cmdArgs := []string{"cmd", "/c", command}
	return RunCommandArgs(cmdArgs, user)
}

// RunCommandArgs run the command
func RunCommandArgs(cmdArgs []string, user string) (string, string, int, error) {
	if user != "" {
		utilLogger.Warningf("RunCommand ignore option: user = %q", user)
	}
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
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
		if exitErr, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode := waitStatus.ExitStatus()
				return stdout, stderr, exitCode, nil
			}
		}
		return stdout, stderr, -1, err
	}

	return stdout, stderr, 0, nil
}

// GetWmic XXX
func GetWmic(target string, query string) (string, error) {
	cpuGet, err := exec.Command("wmic", target, "get", query).Output()
	if err != nil {
		return "", err
	}

	percentages := string(cpuGet)

	lines := strings.Split(percentages, "\r\r\n")

	if len(lines) <= 2 {
		return "", fmt.Errorf("wmic result malformed: [%q]", lines)
	}

	return strings.Trim(lines[1], " "), nil
}

// GetWmicToFloat XXX
func GetWmicToFloat(target string, query string) (float64, error) {
	value, err := GetWmic(target, query)
	if err != nil {
		return 0, err
	}

	ret, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return ret, nil
}
